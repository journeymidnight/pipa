package redis

import (
	"github.com/journeymidnight/pipa/helper"
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
		r, err := InitializeSingle()
		RedisConn = r.(Redis)
		return err
	case "cluster":
		r, err := InitializeCluster()
		RedisConn = r.(Redis)
		return err
	default:
		r, err := InitializeSingle()
		RedisConn = r.(Redis)
		return err
	}
}
