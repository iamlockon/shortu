package web

import (
	"github.com/gin-gonic/gin"
	"github.com/iamlockon/shortu/internal/cache"
	"github.com/iamlockon/shortu/internal/db"
)

type ApiController struct {
	cache cache.CacheClient
	db    db.DbClient
}

func NewApiController(c cache.CacheClient, d db.DbClient) *ApiController {
	return &ApiController{
		cache: c,
		db:    d,
	}
}

func (ctrl *ApiController) getURLHandler(c *gin.Context) {

}

func (ctrl *ApiController) setURLHandler(c *gin.Context) {

}

func (ctrl *ApiController) deleteURLHandler(c *gin.Context) {

}
