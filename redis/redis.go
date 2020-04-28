package redis

import (
	"github.com/journeymidnight/pipa/helper"
	"github.com/journeymidnight/pipa/redis/cluster"
	go_redis "github.com/journeymidnight/pipa/redis/go-redis"
	"github.com/journeymidnight/pipa/redis/single"
)

type Redis interface {
	Close()
	BRPop(key string, timeout uint) ([]string, error)
	LPushSucceed(url, uuid, returnMessage string, blob []byte)
	LPushFailed(uuid, returnMessage string)
}

var RedisConn Redis

func Initialize() {
	switch helper.Config.RedisStore {
	case "single":
		r := go_redis.InitializeSingle()
		RedisConn = r.(Redis)
		break
	case "cluster":
		r := go_redis.InitializeCluster()
		RedisConn = r.(Redis)
		break
	default:
		r := go_redis.InitializeSingle()
		RedisConn = r.(Redis)
	}
}
