package config

import (
	"github.com/iamlockon/shortu/internal/env"
)

type Config struct {
	Cache               StorageConfig
	DB                  StorageConfig
	SrvHost            string
	SrvPort            string
	FilterCap           uint
	FilterWarmupTimeout int
}

// New create configs
func New(cacheConfig, dbConfig StorageConfig) *Config {
	return &Config{
		Cache:               cacheConfig,
		DB:                  dbConfig,
		SrvHost:            env.MustGetString("SRV_HOST", "localhost"),
		SrvPort:            env.MustGetString("SRV_PORT", "8080"),
		FilterCap:           uint(env.MustGetInt("FilterCap", 1000000)),
		FilterWarmupTimeout: env.MustGetInt("FilterWarmupTimeout", 60),
	}
}
