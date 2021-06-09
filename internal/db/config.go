package db

import (
	"fmt"

	"github.com/iamlockon/shortu/internal/env"
)

// NewConfig create MongoConfig from environment variables
func NewConfig() *MongoConfig {
	return &MongoConfig{
		User:     env.MustGetString("MongoUser", ""),
		Password: env.MustGetString("MongoPassword", ""),
		Host:     env.MustGetString("MongoHost", "localhost"),
		Port:     env.MustGetString("MongoPort", "27017"),
		Db:       env.MustGetString("MongoDb", "shortu"),
		Timeout:  env.MustGetInt("MongoTimeout", 10),
	}
}

// GetConnStr returns connection string
func (c *MongoConfig) GetConnStr() string {
	return fmt.Sprintf("mongodb://%s:%s", c.Host, c.Port)
}

func (c *MongoConfig) GetTimeout() int {
	return c.Timeout
}
