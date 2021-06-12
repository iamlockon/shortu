package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheClient interface {
	GetText(ctx context.Context, key string) string
}

type RedisConfig struct {
	user     string
	password string
	host     string
	port     string
	db       string
	timeout  int
}

var _ CacheClient = (*RedisClient)(nil)

type RedisClient struct {
	client  *redis.Client
	timeout time.Duration
}
