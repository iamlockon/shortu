package web

import (
	"fmt"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iamlockon/shortu/internal/cache"
	"github.com/iamlockon/shortu/internal/config"
	"github.com/iamlockon/shortu/internal/db"
	"github.com/iamlockon/shortu/internal/errors"
	filter "github.com/seiflotfy/cuckoofilter"
)

const (
	cacheValidFor24H = 24 * time.Second
)

func NewAPIController(c cache.CacheClient, d db.DBClient, cfg *config.Config, f *filter.Filter) *APIController {
	return &APIController{
		cache:  c,
		db:     d,
		cfg:    cfg,
		filter: f,
	}
}

// redirectURLHandler redirects to original URL
func (ctrl *APIController) redirectURLHandler(c *gin.Context) {
	var uri redirectURLUri
	if err := c.ShouldBindUri(&uri); err != nil {
		fmt.Println("failed to bind uri :", err.Error())
		c.Status(http.StatusBadRequest)
	}
	// test invalid ID
	if !checkValidID(uri.ID) {
		fmt.Println("requested ID is invalid")
		c.Status(http.StatusBadRequest)
		return
	}
	// test against bloom filter
	if !checkBloom(uri.ID, ctrl.filter) {
		fmt.Println("requested ID does not exist")
		c.Status(http.StatusBadRequest)
		return
	}
	// try cache
	original := ctrl.cache.GetText(c.Request.Context(), uri.ID)
	// query db
	if original == "" {
		var err *errors.Error
		if original, err = ctrl.db.GetURL(c.Request.Context(), uri.ID); err != nil {
			if err.Code == errors.URLNotFoundError { // expired, deleted, not exists all go here
				c.Status(http.StatusNotFound)
				return
			}
			c.String(http.StatusInternalServerError, err.Msg)
			return
		}
		// add cache to prevent further db access
		if err = ctrl.cache.SetText(c.Request.Context(), uri.ID, original, cacheValidFor24H); err != nil {
			c.String(http.StatusInternalServerError, err.Msg)
			return
		}
	}
	c.Redirect(http.StatusFound, original)
}

// uploadURLHandler accepts json encoding request body with two string keys:
// 1. url: the original url
// 2. expireAt: expired date
// and returns in json encoding with following two string keys
// 1. id: url id
// 2. shortURL: shortened url
//
// # Example:
// $curl -X POST -H "Content-Type:application/json" http://localhost/api/v1/urls \
//  -d '{"url":"<original_url>","expireAt":"2021-02-08T09:20:41Z"}'
// # Response:
// {"id":"<url_id>","shortUrl":"http://localhost/<url_id>"}
func (ctrl *APIController) uploadURLHandler(c *gin.Context) {
	var req uploadURLReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("invalid request body:", err)
		c.String(http.StatusBadRequest, "invalid request body")
		return
	}

	if !checkExpiredAt(req.ExpiredAt) {
		fmt.Println("invalid expired date: ", req.ExpiredAt)
		c.String(http.StatusBadRequest, fmt.Sprintf("invalid expired date: %s", req.ExpiredAt))
		return
	}
	shorten, err := ctrl.db.UploadURL(c.Request.Context(), req.URL, req.ExpiredAt.UTC().Unix(), ctrl.filter)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Msg)
		return
	}
	// cache it
	if err := ctrl.cache.SetText(c.Request.Context(), shorten, req.URL, cacheValidFor24H); err != nil {
		fmt.Println("set text cache failed: ", err.Msg)
	}
	res := uploadURLRes{
		ID:       shorten,
		ShortURL: fmt.Sprintf("https://%s/%s", ctrl.cfg.SrvHost, shorten),
	}
	c.JSON(http.StatusOK, res)
}

// deleteURLHandler removes one entity if it exists, otherwise do nothing
func (ctrl *APIController) deleteURLHandler(c *gin.Context) {
	c.Status(http.StatusOK)
}
