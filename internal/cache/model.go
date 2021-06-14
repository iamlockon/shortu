package cache

import (
	"context"
	"time"
	"github.com/iamlockon/shortu/internal/errors"
	"github.com/go-redis/redis/v8"
)

type CacheClient interface {
	GetText(ctx context.Context, key string) string
	SetText(ctx context.Context, key, val string, expiry time.Duration) *errors.Error
	Close() error
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
