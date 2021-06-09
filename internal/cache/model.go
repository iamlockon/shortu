package cache

import (
	"github.com/go-redis/redis/v8"
)

type CacheClient interface {
	GetText(key string) string
}

type RedisConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Db       string
	Timeout  int
}

var _ CacheClient = (*RedisClient)(nil)

type RedisClient struct {
	client *redis.Client
}
