package main

/*
NOTE: These tests (functions) are dependent, and can only run once until the environment (e.g. docker-compose) recreates
*/

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"net/http"
	"time"

	"github.com/stretchr/testify/assert"

	"testing"

	"github.com/iamlockon/shortu/internal/cache"
	"github.com/iamlockon/shortu/internal/config"
	"github.com/iamlockon/shortu/internal/db"
	api "github.com/iamlockon/shortu/web/api"
)

const (
	url     = "https://example.com"
	timeout = 2 * time.Second
)

var cfg = config.New(cache.NewConfig(), db.NewConfig())
var uploadURL = fmt.Sprintf("http://%s:%s/api/v1/urls", cfg.SrvHost, cfg.SrvPort)

func TestUserCanUploadURLAndRedirect(t *testing.T) {
	if os.Getenv("env") != "e2e" {
		t.Skip()
	}
	expAt := time.Now().UTC().Add(100 * time.Second)
	body := bytes.NewBufferString(
		fmt.Sprintf(`{"url": "%s", "expireAt": "%s"}`, url, expAt.Format(api.ExpiredDateLayout)))
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadURL, body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Context-Type", "application/json")
	fmt.Println(">>>>Post uploadURL")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("<<<<Post uploadURL")
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var ret struct {
		ID       string `json:"id"`
		ShortURL string `json:"shortUrl"`
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(respBody, &ret); err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, ret.ID)
	assert.NotEmpty(t, ret.ShortURL)
	var redirectURL = fmt.Sprintf("http://%s:%s/%s", cfg.SrvHost, cfg.SrvPort, ret.ID)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	ctx, cancel = context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req, err = http.NewRequestWithContext(ctx, "GET", redirectURL, nil)
	resp2, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()

	assert.Equal(t, http.StatusFound, resp2.StatusCode)

	var deleteURL = fmt.Sprintf("http://%s:%s/api/v1/urls/%s", cfg.SrvHost, cfg.SrvPort, ret.ID)
	ctx, cancel = context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req, err = http.NewRequestWithContext(ctx, "DELETE", deleteURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch Request
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
