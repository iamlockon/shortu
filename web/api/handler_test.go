package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/iamlockon/shortu/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewAPIController(t *testing.T) {
	ctrl := gomock.NewController(t)
	cc, dc := mock.NewMockCacheClient(ctrl), mock.NewMockDBClient(ctrl)
	c := NewAPIController(cc, dc)
	assert.NotNil(t, c.cache)
	assert.NotNil(t, c.db)
}

func TestSetURLHandler(t *testing.T) {
	srv := NewTestAPIServer(t)
	w := httptest.NewRecorder()
	body := bytes.NewBufferString(`{"url": "https://abc.com"}`)
	req, _ := http.NewRequestWithContext(context.Background(), "POST", "/api/v1/urls", body)
	srv.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var ret SetURLRes
	if err := json.Unmarshal(w.Body.Bytes(), &ret); err != nil {
		t.Fail()
	}
	assert.Equal(t, "Xrf2", ret.Res)
}

func TestGetURLHandler(t *testing.T) {
	srv := NewTestAPIServer(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("/%s", "Xfl2"), nil)
	srv.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "getURL", w.Body.String())
}

func TestDeleteURLHandler(t *testing.T) {
	srv := NewTestAPIServer(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(context.Background(), "DELETE", fmt.Sprintf("/api/v1/urls/%s", "Xfl2"), nil)
	srv.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
