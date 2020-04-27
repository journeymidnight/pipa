package single

import (
	"github.com/garyburd/redigo/redis"
	"github.com/journeymidnight/pipa/helper"
	"time"
)

type SingleRedis struct {
	Pool *redis.Pool
}

var pool *redis.Pool

func InitializeSingle() interface{} {
	options := []redis.DialOption{
		redis.DialConnectTimeout(time.Duration(helper.Config.RedisConnectTimeout) * time.Second),
		redis.DialReadTimeout(time.Duration(helper.Config.RedisReadTimeout) * time.Second),
		redis.DialWriteTimeout(time.Duration(helper.Config.RedisWriteTimeout) * time.Second),
	}

	if helper.Config.RedisPassword != "" {
		options = append(options, redis.DialPassword(helper.Config.RedisPassword))
	}

	pool = &redis.Pool{
		MaxIdle:     helper.Config.RedisPoolMaxIdle,
		IdleTimeout: time.Duration(helper.Config.RedisPoolIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", helper.Config.RedisAddress, options...)
			if err != nil {
				helper.Log.Error("connect redis failed:", err)
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			pong, err := c.Do("PING")
			helper.Log.Error("redis PING:", pong)
			return err
		},
	}
	return SingleRedis{Pool: pool}
}

func (s *SingleRedis) Close() {
	err := s.Pool.Close()
	if err != nil {
		helper.Log.Error("can not close redis pool. err:", err)
	}
}

func (s *SingleRedis) BRPop(key string, timeout uint) ([]string, error) {
	c := s.Pool.Get()
	defer c.Close()
	return redis.Strings(c.Do("BRPOP", key, timeout))
}

func (s *SingleRedis) LPushSucceed(url, uuid, returnMessage string, blob []byte) {
	c := s.Pool.Get()
	defer c.Close()
	_, err := c.Do("MULTI")
	if err != nil {
		helper.Log.Error("MULTI do err:", err)
	}
	_, err = c.Do("SET", url, blob)
	if err != nil {
		c.Do("DISCARD")
		helper.Log.Error("SET do err:", err)
	}
	_, err = c.Do("LPUSH", uuid, returnMessage)
	if err != nil {
		c.Do("DISCARD")
		helper.Log.Error("LPUSH do err:", err)
	}
	_, err = c.Do("EXEC")
	if err != nil {
		helper.Log.Error("EXEC do err:", err)
	}
}

func (s *SingleRedis) LPushFailed(uuid, returnMessage string) {
	c := s.Pool.Get()
	defer c.Close()
	_, err := c.Do("LPUSH", uuid, returnMessage)
	if err != nil {
		helper.Log.Error("EXEC do err:", err)
	}
}
