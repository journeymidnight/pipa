package go_redis

import (
	"github.com/go-redis/redis/v7"
	"github.com/journeymidnight/pipa/helper"
	"time"
)

type SingleRedis struct {
	client *redis.Client
}

var client *redis.Client

func InitializeSingle() (interface{}, error) {
	options := &redis.Options{
		Addr:         helper.Config.RedisAddress,
		Password:     helper.Config.RedisPassword,
		DialTimeout:  time.Duration(helper.Config.RedisConnectTimeout) * time.Second,
		ReadTimeout:  time.Duration(helper.Config.RedisReadTimeout) * time.Second,
		WriteTimeout: time.Duration(helper.Config.RedisWriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(helper.Config.RedisPoolIdleTimeout) * time.Second,
	}

	client = redis.NewClient(options)
	_, err := client.Ping().Result()
	if err != nil {
		helper.Log.Error("redis PING err:", err)
		return nil, err
	}
	r := &SingleRedis{client: client}
	return interface{}(r), err
}

func (s *SingleRedis) Close() {
	if err := s.client.Close(); err != nil {
		helper.Log.Error("can not close redis client. err:", err)
	}
}

func (s *SingleRedis) BRPop(key string, timeout uint) ([]string, error) {
	conn := s.client.Conn()
	defer conn.Close()
	do := conn.BRPop(time.Duration(timeout)*time.Second, key)
	strings, err := do.Result()
	if err != nil {
		helper.Log.Error("BRPop err:", err)
	}
	return strings, err
}

func (s *SingleRedis) LPushSucceed(url, uuid, returnMessage string, blob []byte) {
	conn := s.client.Conn()
	defer conn.Close()
	tx := conn.TxPipeline()
	_, err := tx.Set(url, blob, time.Duration(1000*helper.Config.RedisSetDataMaxTime)).Result()
	if err != nil {
		tx.Discard()
		helper.Log.Error("SET do err:", err)
	}
	_, err = tx.LPush(uuid, returnMessage).Result()
	if err != nil {
		tx.Discard()
		helper.Log.Error("LPUSH do err:", err)
	}
	_, err = tx.Exec()
	if err != nil {
		helper.Log.Error("EXEC do err:", err)
	}
}

func (s *SingleRedis) LPushFailed(uuid, returnMessage string) {
	conn := s.client.Conn()
	defer conn.Close()
	_, err := conn.LPush(uuid, returnMessage).Result()
	if err != nil {
		helper.Log.Error("EXEC do err:", err)
	}
}

type ClusterRedis struct {
	cluster *redis.ClusterClient
}

var cluster *redis.ClusterClient

func InitializeCluster() (interface{}, error) {
	clusterRedis := &redis.ClusterOptions{
		Addrs:        helper.Config.RedisGroup,
		Password:     helper.Config.RedisPassword,
		DialTimeout:  time.Duration(helper.Config.RedisConnectTimeout) * time.Second,
		ReadTimeout:  time.Duration(helper.Config.RedisReadTimeout) * time.Second,
		WriteTimeout: time.Duration(helper.Config.RedisWriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(helper.Config.RedisPoolIdleTimeout) * time.Second,
	}
	cluster = redis.NewClusterClient(clusterRedis)
	_, err := cluster.Ping().Result()
	if err != nil {
		helper.Log.Error("Cluster Mode redis PING err:", err)
		return nil, err
	}
	r := &ClusterRedis{cluster: cluster}
	return interface{}(r), err
}

func (c *ClusterRedis) Close() {
	if err := c.cluster.Close(); err != nil {
		helper.Log.Error("can not close redis cluster. err:", err)
	}
}

func (c *ClusterRedis) BRPop(key string, timeout uint) ([]string, error) {
	do := c.cluster.BRPop(time.Duration(timeout)*time.Second, key)
	if _, err := do.Result(); err != nil {
		helper.Log.Error("Cluster Mode: BRPop err:", err)
	}
	return do.Result()
}

func (c *ClusterRedis) LPushSucceed(url, uuid, returnMessage string, blob []byte) {
	_, err := c.cluster.Do("MULTI").Result()
	if err != nil {
		helper.Log.Error("Cluster Mode: MULTI do err:", err)
	}
	_, err = c.cluster.Do("SET", url, blob).Result()
	if err != nil {
		c.cluster.Do("DISCARD")
		helper.Log.Error("Cluster Mode: SET do err:", err)
	}
	_, err = c.cluster.LPush(uuid, returnMessage).Result()
	if err != nil {
		c.cluster.Do("DISCARD")
		helper.Log.Error("Cluster Mode: LPUSH do err:", err)
	}
	_, err = c.cluster.Do("EXEC").Result()
	if err != nil {
		helper.Log.Error("Cluster Mode: EXEC do err:", err)
	}
}

func (c *ClusterRedis) LPushFailed(uuid, returnMessage string) {
	_, err := c.cluster.LPush(uuid, returnMessage).Result()
	if err != nil {
		helper.Log.Error("Cluster Mode: EXEC do err:", err)
	}
}
