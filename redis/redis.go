package redis

import (
	"github.com/journeymidnight/pipa/helper"
	goredis "github.com/journeymidnight/pipa/redis/go-redis"
)

type Redis interface {
	Close()
	BRPop(key string, timeout uint) ([]string, error)
	LPushSucceed(url, uuid, returnMessage string, blob []byte)
	LPushFailed(uuid, returnMessage string)
}

var RedisConn Redis

func Initialize() error {
	switch helper.Config.RedisStore {
	case "single":
		r, err := goredis.InitializeSingle()
		RedisConn = r.(Redis)
		return err
	case "cluster":
		r, err := goredis.InitializeCluster()
		RedisConn = r.(Redis)
		return err
	default:
		r, err := goredis.InitializeSingle()
		RedisConn = r.(Redis)
		return err
	}
}
