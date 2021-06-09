package cache

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/iamlockon/shortu/internal/error"
	"github.com/iamlockon/shortu/internal/config"
)

// New creates a redis client or nil if error presents
func New(config config.StorageConfig) (*RedisClient, *error.Error) {
	opt, err := redis.ParseURL(config.GetConnStr())
	if err != nil {
		return nil, error.New(error.InvalidConfigError, fmt.Sprintf("redis.ParseURL failed: %v", err))
	}

	return &RedisClient{
		client: redis.NewClient(opt),
	}, nil
}

func (c *RedisClient) GetText(key string) string {
	return ""
}
