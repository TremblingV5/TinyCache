package base

import (
	"github.com/caarlos0/env/v9"
)

type Config struct {
	Port     string `env:"TINY_CACHE_PORT",default:"8001"`
	ApiPort  string `env:"TINY_CACHE_API_PORT",default:"9999"`
	StartApi bool   `env:"TINY_CACHE_START_API",default:"false"`
	Master   string `env:"TINY_CACHE_MASTER",default:"localhost:8001"`

	MaxBytes int64 `env:"TINY_CACHE_MAX_BYTES",default:"20480"`
}

func LoadConfig() *Config {
	config := Config{}
	if err := env.Parse(&config); err != nil {
		panic(err)
	}

	return &config
}
