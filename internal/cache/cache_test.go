package cache

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	"github.com/golang/mock/gomock"
	"github.com/iamlockon/shortu/internal/errors"
	"github.com/iamlockon/shortu/mock"
	"github.com/stretchr/testify/assert"
)

const (
	key         = "key1"
	val         = "val1"
	nonexistKey = "nonexist"
)

func TestNewCacheClient_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mock.NewMockStorageConfig(ctrl)
	sc.EXPECT().GetConnStr().Times(1).Return("redis://abc:6379")
	sc.EXPECT().GetTimeout().Times(1).Return(time.Duration(10))
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
	assert.Equal(t, err.Code, errors.InvalidConfigError)
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

// func TestGetText_Timeout(t *testing.T) {
// 	c, _ := New(NewConfig())
// 	rc, m := redismock.NewClientMock()
// 	c.client = rc
// 	c.timeout = 1 * time.Nanosecond
// 	m.ExpectGet(key).SetVal(val)
// 	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
// 	defer cancel()
// 	assert.Empty(t, c.GetText(ctx, key))
// }
