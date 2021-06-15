package main

/*
NOTE: These tests (functions) are dependent, and can only run once until the environment (e.g. docker-compose) recreates
*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	url = "https://example.com"
)

var cfg = config.New(cache.NewConfig(), db.NewConfig())
var uploadURL = fmt.Sprintf("http://%s:%s/api/v1/urls", cfg.SrvHost, cfg.SrvPort)

func UserCanUploadURLAndRedirect() string {
	t := &testing.T{}
	expAt := time.Now().UTC().Add(100 * time.Second)
	body := bytes.NewBufferString(
		fmt.Sprintf(`{"url": "%s", "expireAt": "%s"}`, url, expAt.Format(api.ExpiredDateLayout)))

	resp, err := http.Post(uploadURL, "application/json", body)
	if err != nil {
		t.Fatal(err)
	}
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
	resp2, err := http.Get(redirectURL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	assert.Equal(t, http.StatusFound, resp2.StatusCode)
	return ret.ID
}

func UserCanDeleteURL(id string) {
	t := &testing.T{}

	var deleteURL = fmt.Sprintf("http://%s:%s/api/v1/urls/%s", cfg.SrvHost, cfg.SrvPort, id)
	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Fetch Request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func main() {
	// User stories
	id := UserCanUploadURLAndRedirect()
	UserCanDeleteURL(id)
}
