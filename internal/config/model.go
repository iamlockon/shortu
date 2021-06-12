package config

import (
	"time"
)

type StorageConfig interface {
	GetConnStr() string
	GetTimeout() time.Duration
}
