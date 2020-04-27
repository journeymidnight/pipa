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
		r := single.InitializeSingle()
		RedisConn = r.(Redis)
		break
	case "cluster":
		r := cluster.InitializeCluster()
		RedisConn = r.(Redis)
		break
	default:
		r := single.InitializeSingle()
		RedisConn = r.(Redis)
	}
}
