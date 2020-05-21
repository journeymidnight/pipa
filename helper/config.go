package helper

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

var Config PipaConfig

const (
	DEFAULT_PIPA_CONF_PATH  = "/etc/pipa/pipa.toml"
	DEFAULT_PIPA_FRONT_PATH = "/usr/share/fonts/Chinese_fonts/"
)

type PipaConfig struct {
	S3Domain      []string `toml:"s3domain"`
	LogLevel      string   `toml:"log_level"`
	LogPath       string   `toml:"log_path"`
	WorkersNumber int      `toml:"workers_number"`

	RedisStore           string   `toml:"redis_store"`    // Choose redis connection method
	RedisAddress         string   `toml:"redis_address"`  // redis connection string, e.g localhost:1234
	RedisGroup           []string `toml:"redis_group"`    // Redis cluster connection address
	RedisPassword        string   `toml:"redis_password"` // redis auth password
	RedisMaxRetries 	 int      `toml:"redis_max_retries"`
	RedisConnectTimeout  int      `toml:"redis_connect_timeout"`
	RedisReadTimeout     int      `toml:"redis_read_timeout"`
	RedisWriteTimeout    int      `toml:"redis_write_timeout"`
	RedisPoolMaxIdle     int      `toml:"redis_pool_max_idle"`
	RedisPoolIdleTimeout int      `toml:"redis_pool_idle_timeout"`
	RedisSetDataMaxTime  int      `toml:"redis_set_data_max_time"`

	// This property sets the amount of seconds, after tripping the circuit,
	// to reject requests before allowing attempts again to determine if the circuit should again be closed.
	CacheCircuitCloseSleepWindow int `toml:"cache_circuit_close_sleep_window"`
	// This value is how may consecutive passing requests are required before the circuit is closed
	CacheCircuitCloseRequiredCount int `toml:"cache_circuit_close_required_count"`
	// This property sets the minimum number of requests in a rolling window that will trip the circuit.
	CacheCircuitOpenThreshold     int   `toml:"cache_circuit_open_threshold"`
	CacheCircuitExecTimeout       uint  `toml:"cache_circuit_exec_timeout"`
	CacheCircuitExecMaxConcurrent int64 `toml:"cache_circuit_exec_max_concurrent"`
}

func SetupGlobalConfig() {
	data, err := ioutil.ReadFile(DEFAULT_PIPA_CONF_PATH)
	if err != nil {
		if err != nil {
			panic("Cannot open pipa.toml")
		}
	}
	var c PipaConfig
	_, err = toml.Decode(string(data), &c)
	if err != nil {
		panic("load pipa.toml error: " + err.Error())
	}

	Config.S3Domain = c.S3Domain
	Config.LogLevel = c.LogLevel
	Config.LogPath = c.LogPath
	Config.WorkersNumber = c.WorkersNumber

	Config.RedisStore = c.RedisStore
	Config.RedisAddress = c.RedisAddress
	Config.RedisGroup = c.RedisGroup
	Config.RedisPassword = c.RedisPassword
	Config.RedisMaxRetries = c.RedisMaxRetries
	Config.RedisConnectTimeout = c.RedisConnectTimeout
	Config.RedisReadTimeout = c.RedisReadTimeout
	Config.RedisWriteTimeout = c.RedisWriteTimeout
	Config.RedisPoolMaxIdle = c.RedisPoolMaxIdle
	Config.RedisPoolIdleTimeout = c.RedisPoolIdleTimeout
	Config.RedisSetDataMaxTime = c.RedisSetDataMaxTime

	Config.CacheCircuitCloseSleepWindow = c.CacheCircuitCloseSleepWindow
	Config.CacheCircuitCloseRequiredCount = c.CacheCircuitCloseRequiredCount
	Config.CacheCircuitOpenThreshold = c.CacheCircuitOpenThreshold
	Config.CacheCircuitExecTimeout = c.CacheCircuitExecTimeout
	Config.CacheCircuitExecMaxConcurrent = c.CacheCircuitExecMaxConcurrent
}
