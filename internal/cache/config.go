package cache

import (
	"fmt"

	"github.com/iamlockon/shortu/internal/env"
)

// NewConfig create RedisConfig from environment variables
func NewConfig() *RedisConfig {

	return &RedisConfig{
		User:     env.MustGetString("RedisUser", ""),
		Password: env.MustGetString("RedisPassword", ""),
		Host:     env.MustGetString("RedisHost", "localhost"),
		Port:     env.MustGetString("RedisPort", "6379"),
		Db:       env.MustGetString("RedisDb", "0"),
		Timeout:  env.MustGetInt("RedisTimeout", 10),
	}
}

// GetConnStr returns connection string
func (c *RedisConfig) GetConnStr() string {
	return fmt.Sprintf("redis://%s:%s@%s:%s/%s", c.User, c.Password, c.Host, c.Port, c.Db)
}

func (c *RedisConfig) GetTimeout() int {
	return c.Timeout
}
