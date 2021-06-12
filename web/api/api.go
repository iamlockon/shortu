package web

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/iamlockon/shortu/internal/cache"
	"github.com/iamlockon/shortu/internal/config"
	"github.com/iamlockon/shortu/internal/db"
)

func setupRouter(ca cache.CacheClient, d db.DBClient) *gin.Engine {
	router := gin.Default()
	ctrl := NewAPIController(ca, d)
	router.GET("/:url_id", ctrl.getURLHandler)
	v1 := router.Group("/api/v1")
	{
		v1.POST("/urls", ctrl.setURLHandler)
		v1.DELETE("/urls/:url_id", ctrl.deleteURLHandler)
	}
	return router
}

func Run() {
	cfg := config.New(cache.NewConfig(), db.NewConfig())
	ca, err := cache.New(cfg.Cache)
	if err != nil {
		fmt.Println("failed to new cache: ", err)
		panic(err)
	}
	d, err := db.New(cfg.DB)
	if err != nil {
		fmt.Println("failed to new db: ", err)
		panic(err)
	}

	router := setupRouter(ca, d)
	if err := router.Run(fmt.Sprintf("%s:%s", cfg.SrvHost, cfg.SrvPort)); err != nil {
		panic(err)
	}
}
