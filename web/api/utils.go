package web

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/iamlockon/shortu/mock"
)

func NewTestAPIServer(t *testing.T) *gin.Engine {
	ctrl := gomock.NewController(t)
	ca, d := mock.NewMockCacheClient(ctrl), mock.NewMockDBClient(ctrl)
	return setupRouter(ca, d)
}
