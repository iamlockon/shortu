package config

import (
	"github.com/iamlockon/shortu/internal/env"
)

type Config struct {
	Cache   StorageConfig
	DB      StorageConfig
	SrvHost string
	SrvPort string
}

// New create configs
func New(cacheConfig, dbConfig StorageConfig) *Config {
	return &Config{
		Cache:   cacheConfig,
		DB:      dbConfig,
		SrvHost: env.MustGetString("SrvHost", "localhost"),
		SrvPort: env.MustGetString("SrvPort", "8080"),
	}
}
