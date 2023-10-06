package base

import (
	"github.com/caarlos0/env/v9"
)

type Config struct {
	Name          string `env:"TINY_CACHE_SERVER_NAME" envDefault:"TinyCache"`
	Port          string `env:"TINY_CACHE_PORT" envDefault:"8001"`
	ApiPort       string `env:"TINY_CACHE_API_PORT" envDefault:"9999"`
	StartApi      bool   `env:"TINY_CACHE_START_API" envDefault:"false"`
	Master        string `env:"TINY_CACHE_MASTER" envDefault:"localhost:8001"`
	SecondaryList string `env:"TINY_CACHE_SECONDARY_LIST" envDefault:""`
	SecondaryNum  int    `env:"TINY_CACHE_SECONDARY_NUM" envDefault:"0"`

	SnowFlakeNodeNum int `env:"TINY_CACHE_SNOWFLAKE_NODE_NUM" envDefault:"1"`

	EliminationMethod string `env:"TINY_CACHE_ELIMINATION_METHOD" envDefault:"LRU"`

	MaxBytes int64 `env:"TINY_CACHE_MAX_BYTES" envDefault:"20480"`

	LogFileName   string `env:"TINY_CACHE_LOG_FILE_NAME" envDefault:"/tmp/tiny_cache.log"`
	LogMaxSize    int    `env:"TINY_CACHE_LOG_MAX_SIZE" envDefault:"500"`
	LogMaxBackups int    `env:"TINY_CACHE_LOG_MAX_BACKUPS" envDefault:"3"`
	LogMaxAge     int    `env:"TINY_CACHE_LOG_MAX_AGE" envDefault:"7"`
	LogLevel      string `env:"TINY_CACHE_LOG_LEVEL" envDefault:"info"`
}

func LoadConfig() *Config {
	config := Config{}
	if err := env.Parse(&config); err != nil {
		panic(err)
	}

	return &config
}
