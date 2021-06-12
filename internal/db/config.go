package db

import (
	"fmt"
	"time"

	"github.com/iamlockon/shortu/internal/env"
)

// NewConfig create PgConfig from environment variables
func NewConfig() *PgConfig {
	return &PgConfig{
		user:     env.MustGetString("PgUser", "pg"),
		password: env.MustGetString("PgPassword", "123456"),
		host:     env.MustGetString("PgHost", "db"),
		db:       env.MustGetString("PgDb", "shortu"),
		timeout:  env.MustGetInt("PgTimeout", 10),
	}
}

// GetConnStr returns connection string
// example "postgres://username:password@localhost:5432/database_name"
func (c *PgConfig) GetConnStr() string {
	return fmt.Sprintf("postgres://%s:%s@%s:5432/%s", c.user, c.password, c.host, c.db)
}

func (c *PgConfig) GetTimeout() time.Duration {
	return time.Duration(c.timeout)
}
