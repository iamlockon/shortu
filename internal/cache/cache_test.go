package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	"github.com/golang/mock/gomock"
	E "github.com/iamlockon/shortu/internal/errors"
	"github.com/iamlockon/shortu/mock"
	"github.com/stretchr/testify/assert"
)

const (
	key         = "key1"
	val         = "val1"
	exp         = time.Duration(time.Hour * 24)
	nonexistKey = "nonexist"
)

func TestNewCacheClient_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mock.NewMockStorageConfig(ctrl)
	sc.EXPECT().GetConnStr().Times(1).Return("redis://abc:6379")
	sc.EXPECT().GetTimeout().Times(1).Return(10 * time.Second)
	cc, err := New(sc)
	assert.Nil(t, err)
	assert.NotNil(t, cc)
	assert.NotNil(t, cc.client)
}

func TestNewCacheClient_WithWrongConnStr_ShouldReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mock.NewMockStorageConfig(ctrl)
	sc.EXPECT().GetConnStr().Times(1).Return("randomsttrrrrr")
	cc, err := New(sc)
	assert.Nil(t, cc)
	assert.NotNil(t, err)
	assert.Equal(t, err.Code, E.InvalidConfigError)
}

func TestGetText_OK(t *testing.T) {
	c, err := New(NewConfig())
	assert.Nil(t, err)
	assert.NotNil(t, c.client)
	rc, m := redismock.NewClientMock()
	c.client = rc
	m.ExpectGet(key).SetVal(val)
	ctx := context.Background()
	assert.Equal(t, val, c.GetText(ctx, key))
	assert.Empty(t, c.GetText(ctx, nonexistKey))
}

func TestSetText_OK(t *testing.T) {
	c, _ := New(NewConfig())
	rc, m := redismock.NewClientMock()
	c.client = rc
	m.ExpectSet(key, val, exp)
	ctx := context.Background()
	assert.Nil(t, c.SetText(ctx, key, val, exp))
}

func TestSetText_Failed(t *testing.T) {
	c, _ := New(NewConfig())
	rc, m := redismock.NewClientMock()
	c.client = rc
	m.ExpectSet(key, val, exp).SetErr(errors.New("anything"))
	ctx := context.Background()
	assert.NotNil(t, c.SetText(ctx, key, val, exp))
}
