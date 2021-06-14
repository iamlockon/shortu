package web

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/iamlockon/shortu/internal/cache"
	"github.com/iamlockon/shortu/internal/config"
	"github.com/iamlockon/shortu/internal/db"
	"github.com/iamlockon/shortu/mock"
)

type TestAPIServer struct {
	router *gin.Engine
	c      *mock.MockCacheClient
	d      *mock.MockDBClient
}

func NewTestAPIServer(t *testing.T) *TestAPIServer {
	ctrl := gomock.NewController(t)
	ca, d := mock.NewMockCacheClient(ctrl), mock.NewMockDBClient(ctrl)
	cfg := config.New(cache.NewConfig(), db.NewConfig())
	// ctx, cancel := context.WithTimeout(gomock.Any(), time.Duration(cfg.FilterWarmupTimeout)*time.Second)
	// defer cancel()
	d.EXPECT().LoadURL(gomock.Any(), gomock.Any()).Times(1).Return(nil)
	return &TestAPIServer{
		router: setupRouter(ca, d, cfg),
		c:      ca,
		d:      d,
	}
}
