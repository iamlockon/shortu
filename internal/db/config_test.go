package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	dc := NewConfig()
	assert.Equal(t, dc.user, "pg")
	assert.Equal(t, dc.password, "123456")
	assert.Equal(t, dc.host, "db")
	assert.Equal(t, dc.db, "shortu")
	assert.Equal(t, dc.timeout, 10)
}

func TestGetConnStr(t *testing.T) {
	rc := NewConfig()
	assert.Equal(t, "postgres://pg:123456@db:5432/shortu", rc.GetConnStr())
}

func TestGetTimeout(t *testing.T) {
	rc := NewConfig()
	timeout := 10
	rc.timeout = timeout
	assert.Equal(t, time.Duration(timeout), rc.GetTimeout())
}
