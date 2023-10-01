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

	EliminationMethod string `env:"TINY_CACHE_ELIMINATION_METHOD" envDefault:"LRU"`

	MaxBytes int64 `env:"TINY_CACHE_MAX_BYTES" envDefault:"20480"`
}

func LoadConfig() *Config {
	config := Config{}
	if err := env.Parse(&config); err != nil {
		panic(err)
	}

	return &config
}
