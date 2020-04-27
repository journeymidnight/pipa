package cluster

import (
	"github.com/chasex/redis-go-cluster"
	"github.com/journeymidnight/pipa/helper"
	"time"
)

type ClusterRedis struct {
	Cluster redis.Cluster
}

var cluster redis.Cluster

func InitializeCluster() interface{} {
	var err error
	cluster, err = redis.NewCluster(
		&redis.Options{
			StartNodes:   helper.Config.RedisGroup,
			ConnTimeout:  time.Duration(helper.Config.RedisConnectTimeout) * time.Second,
			ReadTimeout:  time.Duration(helper.Config.RedisReadTimeout) * time.Second,
			WriteTimeout: time.Duration(helper.Config.RedisWriteTimeout) * time.Second,
			KeepAlive:    helper.Config.RedisKeepalived,
			AliveTime:    time.Duration(helper.Config.RedisAlivedTime) * time.Second,
		})

	if err != nil {
		helper.Log.Error("redis.New error: %s", err.Error())
	}
	return ClusterRedis{Cluster: cluster}
}

func (c *ClusterRedis) Close() {
}

func (c *ClusterRedis) BRPop(key string, timeout uint) ([]string, error) {
	return redis.Strings(c.Cluster.Do("BRPOP", key, timeout))
}

func (c *ClusterRedis) LPushSucceed(url, uuid, returnMessage string, blob []byte) {
	_, err := c.Cluster.Do("MULTI")
	if err != nil {
		helper.Log.Error("MULTI do err:", err)
	}
	_, err = c.Cluster.Do("SET", url, blob)
	if err != nil {
		c.Cluster.Do("DISCARD")
		helper.Log.Error("SET do err:", err)
	}
	_, err = c.Cluster.Do("LPUSH", uuid, returnMessage)
	if err != nil {
		c.Cluster.Do("DISCARD")
		helper.Log.Error("LPUSH do err:", err)
	}
	_, err = c.Cluster.Do("EXEC")
	if err != nil {
		helper.Log.Error("EXEC do err:", err)
	}
}

func (c *ClusterRedis) LPushFailed(uuid, returnMessage string) {
	_, err := c.Cluster.Do("LPUSH", uuid, returnMessage)
	if err != nil {
		helper.Log.Error("EXEC do err:", err)
	}
}
