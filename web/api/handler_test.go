package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/iamlockon/shortu/internal/cache"
	"github.com/iamlockon/shortu/internal/config"
	"github.com/iamlockon/shortu/internal/db"
	"github.com/iamlockon/shortu/mock"
	filter "github.com/seiflotfy/cuckoofilter"
	"github.com/stretchr/testify/assert"
)

func TestNewAPIController(t *testing.T) {
	ctrl := gomock.NewController(t)
	cc, dc := mock.NewMockCacheClient(ctrl), mock.NewMockDBClient(ctrl)
	cfg := config.New(cache.NewConfig(), db.NewConfig())
	f := filter.NewFilter(cfg.FilterCap)
	c := NewAPIController(cc, dc, cfg, f)
	assert.NotNil(t, c.cache)
	assert.NotNil(t, c.db)
}

func TestUploadURLHandler(t *testing.T) {
	srv := NewTestAPIServer(t)
	url, shorten := "https://abc.com", "shorteeee"
	w := httptest.NewRecorder()
	expAt := time.Now().UTC().Add(100 * time.Second)
	gomock.InOrder(
		srv.d.EXPECT().UploadURL(context.Background(), url, expAt.UTC().Unix(), gomock.Any()).Times(1).Return(shorten, nil),
		srv.c.EXPECT().SetText(gomock.Any(), shorten, url, cacheValidFor24H).Times(1).Return(nil),
	)

	body := bytes.NewBufferString(
		fmt.Sprintf(`{"url": "%s", "expired_at": "%s"}`, url, expAt.Format(expiredDateLayout)))
	req, _ := http.NewRequestWithContext(
		context.Background(),
		"POST",
		"/api/v1/urls",
		body,
	)
	req.Header.Add("Content-Type", "application/json")
	srv.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var ret uploadURLRes
	if err := json.Unmarshal(w.Body.Bytes(), &ret); err != nil {
		t.Fail()
	}
	assert.Equal(t, shorten, ret.ID)
	assert.Equal(t, fmt.Sprintf("https://%s/%s", "localhost", shorten), ret.ShortURL)
}

func TestUploadURLHandler_BadExpiredAt_ShouldFail(t *testing.T) {
	srv := NewTestAPIServer(t)
	url := "https://abc.com"
	w := httptest.NewRecorder()
	expAt := time.Now().UTC().Add(-100 * time.Second)
	body := bytes.NewBufferString(
		fmt.Sprintf(`{"url": "%s", "expired_at": "%s"}`, url, expAt.Format(expiredDateLayout)))
	req, _ := http.NewRequestWithContext(
		context.Background(),
		"POST",
		"/api/v1/urls",
		body,
	)
	req.Header.Add("Content-Type", "application/json")
	srv.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRedirectURLHandler_IfIDNotExist_ShouldNotRedirect(t *testing.T) {
	srv := NewTestAPIServer(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("/%s", "nonexist"), nil)
	srv.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRedirectURLHandler_Exists(t *testing.T) {
	srv := NewTestAPIServer(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("/%s", "Xfl2"), nil)
	srv.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "https://abc.com", w.Header().Get("Location"))
}

func TestDeleteURLHandler(t *testing.T) {
	srv := NewTestAPIServer(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(context.Background(), "DELETE", fmt.Sprintf("/api/v1/urls/%s", "Xfl2"), nil)
	srv.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
