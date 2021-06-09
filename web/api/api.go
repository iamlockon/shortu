package web

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/iamlockon/shortu/internal/cache"
	"github.com/iamlockon/shortu/internal/config"
	"github.com/iamlockon/shortu/internal/db"
)

func Run() {
	config := config.New(cache.NewConfig(), db.NewConfig())
	ca, err := cache.New(config.Cache)
	if err != nil {
		fmt.Println("failed to new cache: ", err)
		panic(err)
	}
	d, err := db.New(config.DB)
	if err != nil {
		fmt.Println("failed to new db: ", err)
		panic(err)
	}
	ctrl := NewApiController(ca, d)
	router := gin.Default()
	router.GET("/:url_id", ctrl.getURLHandler)
	v1 := router.Group("/api/v1")
	{
		v1.POST("/urls", ctrl.setURLHandler)
		v1.DELETE("/urls/:url_id", ctrl.deleteURLHandler)
	}

	router.Run(fmt.Sprintf("%s:%s", config.SrvHost, config.SrvPort))
}
