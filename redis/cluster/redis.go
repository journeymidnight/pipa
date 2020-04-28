package cluster

import (
	"github.com/gomodule/redigo/redis"
	"github.com/journeymidnight/pipa/helper"
	"github.com/mna/redisc"
	"time"
)

type ClusterRedis struct {
	Cluster *redisc.Cluster
}

var cluster *redisc.Cluster

func InitializeCluster() interface{} {
	options := []redis.DialOption{
		redis.DialConnectTimeout(time.Duration(helper.Config.RedisConnectTimeout) * time.Second),
		redis.DialReadTimeout(time.Duration(helper.Config.RedisReadTimeout) * time.Second),
		redis.DialWriteTimeout(time.Duration(helper.Config.RedisWriteTimeout) * time.Second),
	}

	if helper.Config.RedisPassword != "" {
		options = append(options, redis.DialPassword(helper.Config.RedisPassword))
	}

	cluster = &redisc.Cluster{
		StartupNodes: helper.Config.RedisGroup,
		DialOptions:  options,
		CreatePool:   createPool,
	}
	// initialize its mapping
	if err := cluster.Refresh(); err != nil {
		helper.Log.Error("Refresh failed: %v", err)
	}
	r := &ClusterRedis{Cluster: cluster}
	return interface{}(r)
}

func createPool(addr string, opts ...redis.DialOption) (*redis.Pool, error) {
	return &redis.Pool{
		MaxIdle:     helper.Config.RedisPoolMaxIdle,
		IdleTimeout: time.Duration(helper.Config.RedisPoolIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr, opts...)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				helper.Log.Error("redis PING:", err)
			}
			return err
		},
	}, nil
}

func (c *ClusterRedis) Close() {
	err := c.Cluster.Close()
	if err != nil {
		helper.Log.Error("can not close redis pool. err:", err)
	}
}

func (c *ClusterRedis) BRPop(key string, timeout uint) ([]string, error) {
	coon := c.Cluster.Get()
	defer coon.Close()
	return redis.Strings(coon.Do("BRPOP", key, timeout))
}

func (c *ClusterRedis) LPushSucceed(url, uuid, returnMessage string, blob []byte) {
	coon := c.Cluster.Get()
	defer coon.Close()
	_, err := coon.Do("MULTI")
	if err != nil {
		helper.Log.Error("MULTI do err:", err)
	}
	_, err = coon.Do("SET", url, blob)
	if err != nil {
		coon.Do("DISCARD")
		helper.Log.Error("SET do err:", err)
	}
	_, err = coon.Do("LPUSH", uuid, returnMessage)
	if err != nil {
		coon.Do("DISCARD")
		helper.Log.Error("LPUSH do err:", err)
	}
	_, err = coon.Do("EXEC")
	if err != nil {
		helper.Log.Error("EXEC do err:", err)
	}
}

func (c *ClusterRedis) LPushFailed(uuid, returnMessage string) {
	coon := c.Cluster.Get()
	defer coon.Close()
	_, err := coon.Do("LPUSH", uuid, returnMessage)
	if err != nil {
		helper.Log.Error("EXEC do err:", err)
	}
}
