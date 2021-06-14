package web

import (
	"time"

	"github.com/iamlockon/shortu/internal/cache"
	"github.com/iamlockon/shortu/internal/config"
	"github.com/iamlockon/shortu/internal/db"
	filter "github.com/seiflotfy/cuckoofilter"
)

const (
	expiredDateLayout = "2006-01-02T15:04:05.000Z"
)

type uploadURLRes struct {
	ID       string `json:"id"`
	ShortURL string `json:"shortUrl"`
}

type uploadURLReq struct {
	URL       string    `json:"url" binding:"url,required"`
	ExpiredAt time.Time `json:"expired_at" binding:"required"`
}

type deleteURLUri struct {
	ID string `uri:"id" binding:"max=10,required"`
}

type redirectURLUri struct {
	ID string `uri:"id" binding:"max=10,required"`
}

type APIController struct {
	cache  cache.CacheClient
	db     db.DBClient
	cfg    *config.Config
	filter *filter.Filter
}
