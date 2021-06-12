package web

import (
	"github.com/gin-gonic/gin"
	"github.com/iamlockon/shortu/internal/cache"
	"github.com/iamlockon/shortu/internal/db"

	"net/http"
)

type APIController struct {
	cache cache.CacheClient
	db    db.DBClient
}

func NewAPIController(c cache.CacheClient, d db.DBClient) *APIController {
	return &APIController{
		cache: c,
		db:    d,
	}
}

func (ctrl *APIController) getURLHandler(c *gin.Context) {
	c.String(http.StatusOK, "getURL")
}

func (ctrl *APIController) setURLHandler(c *gin.Context) {
	res := SetURLRes{
		Res: "Xrf2",
	}
	c.JSON(http.StatusOK, res)
}

func (ctrl *APIController) deleteURLHandler(c *gin.Context) {
	c.String(http.StatusOK, "")
}
