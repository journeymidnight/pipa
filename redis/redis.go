package redis

import (
	"github.com/garyburd/redigo/redis"
	"github.com/journeymidnight/pipa/helper"
	"time"
)

var Pool *redis.Pool

func Initialize() {
	options := []redis.DialOption{
		redis.DialConnectTimeout(time.Duration(helper.Config.RedisConnectTimeout) * time.Second),
		redis.DialReadTimeout(time.Duration(helper.Config.RedisReadTimeout) * time.Second),
		redis.DialWriteTimeout(time.Duration(helper.Config.RedisWriteTimeout) * time.Second),
	}

	if helper.Config.RedisPassword != "" {
		options = append(options, redis.DialPassword(helper.Config.RedisPassword))
	}

	Pool = &redis.Pool{
		MaxIdle:     helper.Config.RedisPoolMaxIdle,
		IdleTimeout: time.Duration(helper.Config.RedisPoolIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", helper.Config.RedisAddress, options...)
			if err != nil {
				helper.Log.Error("connect redis failed:",err)
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				helper.Log.Error("redis PING:",err)
			}
			return err
		},
	}
}

func Close() {
	err := Pool.Close()
	if err != nil {
		helper.Log.Error("can not close redis pool. err:", err)
	}
}

func BRPop(key string, timeout uint) ([]string, error) {
	c := Pool.Get()
	defer c.Close()
	return redis.Strings(c.Do("BRPOP", key, timeout))
}
