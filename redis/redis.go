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

func Initialize() {
	switch helper.Config.RedisStore {
	case "single":
		r := goredis.InitializeSingle()
		RedisConn = r.(Redis)
		break
	case "cluster":
		r := goredis.InitializeCluster()
		RedisConn = r.(Redis)
		break
	default:
		r := goredis.InitializeSingle()
		RedisConn = r.(Redis)
	}
}
