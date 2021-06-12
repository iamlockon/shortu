package cache

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/iamlockon/shortu/internal/config"
	"github.com/iamlockon/shortu/internal/errors"
)

// New creates a redis client or nil if error presents
func New(cfg config.StorageConfig) (*RedisClient, *errors.Error) {
	opt, err := redis.ParseURL(cfg.GetConnStr())
	if err != nil {
		return nil, errors.New(errors.InvalidConfigError, fmt.Sprintf("redis.ParseURL failed: %v", err))
	}

	return &RedisClient{
		client:  redis.NewClient(opt),
		timeout: cfg.GetTimeout(),
	}, nil
}

func (c *RedisClient) GetText(ctx context.Context, key string) string {
	rCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	s, err := c.client.Get(rCtx, key).Result()
	if err != nil {
		return ""
	}
	return s
}
