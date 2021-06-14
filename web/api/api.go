package web

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iamlockon/shortu/internal/cache"
	"github.com/iamlockon/shortu/internal/config"
	"github.com/iamlockon/shortu/internal/db"

	filter "github.com/seiflotfy/cuckoofilter"
)

func setupRouter(ca cache.CacheClient, d db.DBClient, cfg *config.Config) *gin.Engine {
	router := gin.Default()
	f := filter.NewFilter(cfg.FilterCap)
	// warm up filter, this is a blocking op
	fmt.Println(">>>> begin to warm up filter")
	startLoadTime := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.FilterWarmupTimeout)*time.Second)
	defer cancel()
	if err := d.LoadURL(ctx, f); err != nil {
		fmt.Println("failed to warm up filter: ", err.Msg)
		panic(err.Msg)
	}
	fmt.Println("<<<< finish warming up filter, took ", time.Since(startLoadTime))

	ctrl := NewAPIController(ca, d, cfg, f)
	router.GET("/:id", ctrl.redirectURLHandler)
	v1 := router.Group("/api/v1")
	{
		v1.POST("/urls", ctrl.uploadURLHandler)
		v1.DELETE("/urls/:id", ctrl.deleteURLHandler)
	}
	return router
}

func Run() {
	cfg := config.New(cache.NewConfig(), db.NewConfig())
	ca, err := cache.New(cfg.Cache)
	defer ca.Close()
	if err != nil {
		fmt.Println("failed to new cache: ", err)
		panic(err)
	}
	d, err := db.New(cfg.DB)
	defer d.Close()
	if err != nil {
		fmt.Println("failed to new db: ", err)
		panic(err)
	}

	router := setupRouter(ca, d, cfg)
	if err := router.Run(fmt.Sprintf("%s:%s", cfg.SrvHost, cfg.SrvPort)); err != nil {
		panic(err)
	}
}
