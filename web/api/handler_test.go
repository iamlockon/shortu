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

const (
	url, shorten = "https://abc.com", "shorteeee"
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

	w := httptest.NewRecorder()
	expAt := time.Now().UTC().Add(100 * time.Second).Round(time.Millisecond)
	gomock.InOrder(
		srv.d.EXPECT().UploadURL(context.Background(), url, expAt.UTC(), gomock.Any()).Times(1).Return(shorten, nil),
		srv.c.EXPECT().SetText(gomock.Any(), shorten, url, cacheValidFor24H).Times(1).Return(nil),
	)

	body := bytes.NewBufferString(
		fmt.Sprintf(`{"url": "%s", "expireAt": "%s"}`, url, expAt.Format(ExpiredDateLayout)))
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
	w := httptest.NewRecorder()
	expAt := time.Now().UTC().Add(-100 * time.Second)
	body := bytes.NewBufferString(
		fmt.Sprintf(`{"url": "%s", "expireAt": "%s"}`, url, expAt.Format(ExpiredDateLayout)))
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

func TestRedirectURLHandler_IfIDNotExistInCuckoo_ShouldNotRedirect(t *testing.T) {
	srv := NewTestAPIServer(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("/%s", "nonexist"), nil)
	srv.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRedirectURLHandler_InvalidID(t *testing.T) {
	srv := NewTestAPIServer(t)
	tests := []struct {
		in  string
		out int
	}{
		{"...;;;", http.StatusBadRequest},
		{"X+fwef", http.StatusBadRequest},
		{"fwef033ff333f3ff", http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("/%s", tt.in), nil)
			srv.router.ServeHTTP(w, req)
			assert.Equal(t, tt.out, w.Code)
		})
	}
}

func TestRedirectURLHandler_ExistsInCache(t *testing.T) {
	srv := NewTestAPIServer(t)
	w := httptest.NewRecorder()
	srv.f.InsertUnique([]byte(shorten))
	srv.c.EXPECT().GetText(gomock.Any(), shorten).Times(1).Return(url)
	req, _ := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("/%s", shorten), nil)
	srv.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, url, w.Header().Get("Location"))
}

func TestRedirectURLHandler_ExistInDB(t *testing.T) {
	t.Skip()
}

func TestDeleteURLHandler_NotExistInCuckoo_ShouldReturnNotFound(t *testing.T) {
	srv := NewTestAPIServer(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(context.Background(), "DELETE", fmt.Sprintf("/api/v1/urls/%s", "Xfl2"), nil)
	srv.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteURLHandler_InvalidID_ShouldReturnBadRequest(t *testing.T) {
	srv := NewTestAPIServer(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(context.Background(), "DELETE", fmt.Sprintf("/api/v1/urls/%s", "+fse_"), nil)
	srv.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteURLHandler_IfPassChecksAndNoDBError_ShouldReturnOK(t *testing.T) {
	srv := NewTestAPIServer(t)
	w := httptest.NewRecorder()
	shorten := "abcdegf"
	srv.f.InsertUnique([]byte(shorten))
	req, _ := http.NewRequestWithContext(context.Background(), "DELETE", fmt.Sprintf("/api/v1/urls/%s", shorten), nil)
	srv.d.EXPECT().DeleteURL(gomock.Any(), shorten).Times(1).Return(nil)
	srv.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
