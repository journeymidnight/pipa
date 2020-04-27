package redis

import (
	"github.com/journeymidnight/pipa/helper"
	"github.com/journeymidnight/pipa/redis/cluster"
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
		RedisConn = single.InitializeSingle().(Redis)
		break
	case "cluster":
		RedisConn = cluster.InitializeCluster().(Redis)
		break
	default:
		RedisConn = single.InitializeSingle().(Redis)
	}
}
