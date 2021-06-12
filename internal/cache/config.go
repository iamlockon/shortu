package cache

import (
	"fmt"
	"time"

	"github.com/iamlockon/shortu/internal/env"
)

// NewConfig create RedisConfig from environment variables
func NewConfig() *RedisConfig {

	return &RedisConfig{
		user:     env.MustGetString("RedisUser", ""),
		password: env.MustGetString("RedisPassword", ""),
		host:     env.MustGetString("RedisHost", "localhost"),
		port:     env.MustGetString("RedisPort", "6379"),
		db:       env.MustGetString("RedisDb", "0"),
		timeout:  env.MustGetInt("RedisTimeout", 10),
	}
}

// GetConnStr returns connection string
func (c *RedisConfig) GetConnStr() string {
	return fmt.Sprintf("redis://%s:%s@%s:%s/%s", c.user, c.password, c.host, c.port, c.db)
}

func (c *RedisConfig) GetTimeout() time.Duration {
	return time.Duration(c.timeout)
}
