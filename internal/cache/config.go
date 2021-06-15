package cache

import (
	"fmt"
	"time"

	"github.com/iamlockon/shortu/internal/env"
)

// NewConfig create RedisConfig from environment variables
func NewConfig() *RedisConfig {

	return &RedisConfig{
		user:     env.MustGetString("REDIS_USER", ""),
		password: env.MustGetString("REDIS_PASSWORD", ""),
		host:     env.MustGetString("REDIS_HOST", "localhost"),
		port:     env.MustGetString("REDIS_PORT", "6379"),
		db:       env.MustGetString("REDIS_DB", "0"),
		timeout:  env.MustGetInt("REDIS_TIMEOUT", 10),
	}
}

// GetConnStr returns connection string
func (c *RedisConfig) GetConnStr() string {
	return fmt.Sprintf("redis://%s:%s@%s:%s/%s", c.user, c.password, c.host, c.port, c.db)
}

func (c *RedisConfig) GetTimeout() time.Duration {
	return time.Duration(c.timeout) * time.Second
}
