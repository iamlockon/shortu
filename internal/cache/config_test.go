package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig_Default(t *testing.T) {
	rc := NewConfig()
	assert.Equal(t, rc.user, "")
	assert.Equal(t, rc.password, "")
	assert.Equal(t, rc.host, "localhost")
	assert.Equal(t, rc.port, "6379")
	assert.Equal(t, rc.db, "0")
	assert.Equal(t, rc.timeout, 10)
}

func TestGetConnStr(t *testing.T) {
	rc := NewConfig()
	assert.Equal(t, "redis://:@localhost:6379/0", rc.GetConnStr())
}

func TestGetTimeout(t *testing.T) {
	rc := NewConfig()
	timeout := 10
	rc.timeout = timeout
	assert.Equal(t, time.Second*time.Duration(timeout), rc.GetTimeout())
}
