package config

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"

	"github.com/iamlockon/shortu/mock"
)

func TestNew_Default(t *testing.T) {
	ctrl := gomock.NewController(t)
	ca, db := mock.NewMockStorageConfig(ctrl), mock.NewMockStorageConfig(ctrl)
	rc := New(ca, db)
	assert.Equal(t, "localhost", rc.SrvHost)
	assert.Equal(t, "8080", rc.SrvPort)
}

func TestNew_SetEnv(t *testing.T) {
	ctrl := gomock.NewController(t)
	ca, db := mock.NewMockStorageConfig(ctrl), mock.NewMockStorageConfig(ctrl)
	host, port := "0.0.0.0", "2345"
	if err := os.Setenv("SrvHost", host); err != nil {
		t.Fail()
	}
	if err := os.Setenv("SrvPort", port); err != nil {
		t.Fail()
	}
	rc := New(ca, db)
	assert.Equal(t, host, rc.SrvHost)
	assert.Equal(t, port, rc.SrvPort)
}
