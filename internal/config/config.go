package config

import (
	"github.com/iamlockon/shortu/internal/cache"
)

type Config struct {
	cache StorageConfig
	db StorageConfig
}

// New create configs
func New() *Config {
	return &Config{
		cache: cache.NewConfig()
		db: db.NewConfig()
	}
}