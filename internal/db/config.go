package db

import (
	"fmt"
	"time"

	"github.com/iamlockon/shortu/internal/env"
)

// NewConfig create PgConfig from environment variables
func NewConfig() *PgConfig {
	return &PgConfig{
		user:     env.MustGetString("PG_USER", "pg"),
		password: env.MustGetString("PG_PASSWORD", "123456"),
		host:     env.MustGetString("PG_HOST", "db"),
		db:       env.MustGetString("PG_DB", "shortu"),
		timeout:  env.MustGetInt("PG_TIMEOUT", 10),
	}
}

// GetConnStr returns connection string
// example "postgres://username:password@localhost:5432/database_name"
func (c *PgConfig) GetConnStr() string {
	return fmt.Sprintf("postgres://%s:%s@%s:5432/%s", c.user, c.password, c.host, c.db)
}

func (c *PgConfig) GetTimeout() time.Duration {
	return time.Duration(c.timeout) * time.Second
}
