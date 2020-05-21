package circuitbreak

import (
	"errors"
	"time"

	"github.com/cep21/circuit"
	"github.com/cep21/circuit/closers/hystrix"
	"github.com/journeymidnight/pipa/helper"
)

var (
	CacheCircuitIsOpenErr = errors.New("redis circuit is open now!")
)

func NewCacheCircuit() *circuit.Circuit {
	return circuit.NewCircuitFromConfig("Pipa Redis", circuit.Config{
		General: circuit.GeneralConfig{
			OpenToClosedFactory: hystrix.CloserFactory(hystrix.ConfigureCloser{
				SleepWindow:                  time.Duration(helper.Config.CacheCircuitCloseSleepWindow) * time.Second,
				RequiredConcurrentSuccessful: int64(helper.Config.CacheCircuitCloseRequiredCount),
			}),
			ClosedToOpenFactory: hystrix.OpenerFactory(hystrix.ConfigureOpener{
				RequestVolumeThreshold: int64(helper.Config.CacheCircuitOpenThreshold),
			}),
		},
		Execution: circuit.ExecutionConfig{
			Timeout:               time.Duration(helper.Config.CacheCircuitExecTimeout) * time.Second,
			MaxConcurrentRequests: helper.Config.CacheCircuitExecMaxConcurrent,
		},
	})
}
