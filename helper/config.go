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

	RedisAddress         string `toml:"redis_address"`  // redis connection string, e.g localhost:1234
	RedisPassword        string `toml:"redis_password"` // redis auth password
	RedisConnectTimeout  int    `toml:"redis_connect_timeout"`
	RedisReadTimeout     int    `toml:"redis_read_timeout"`
	RedisWriteTimeout    int    `toml:"redis_write_timeout"`
	RedisPoolMaxIdle     int    `toml:"redis_pool_max_idle"`
	RedisPoolIdleTimeout int    `toml:"redis_pool_idle_timeout"`
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

	Config.RedisAddress = c.RedisAddress
	Config.RedisPassword = c.RedisPassword
	Config.RedisConnectTimeout = c.RedisConnectTimeout
	Config.RedisReadTimeout = c.RedisReadTimeout
	Config.RedisWriteTimeout = c.RedisWriteTimeout
	Config.RedisPoolMaxIdle = c.RedisPoolMaxIdle
	Config.RedisPoolIdleTimeout = c.RedisPoolIdleTimeout
}
