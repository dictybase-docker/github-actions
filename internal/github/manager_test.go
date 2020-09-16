package github

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"

	gh "github.com/google/go-github/v32/github"
)

const (
	baseURLPath = "/api-v3"
)

func handleCommitComparison(w http.ResponseWriter, r *http.Request) {
	dir, err := os.Getwd()
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("unable to get current dir %s", err),
			http.StatusInternalServerError,
		)
		return
	}
	path := filepath.Join(
		filepath.Dir(dir), "../testdata", "event.json",
	)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		http.Error(
			w,
			"unable to read test file",
			http.StatusInternalServerError,
		)
		return
	}
	if _, err := w.Write(b); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
	}
}

func fakeGhServerClient() (*httptest.Server, *gh.Client) {
	apiHandler := http.NewServeMux()
	apiHandler.HandleFunc("/repos/o/r/compare/b...h", handleCommitComparison)
	server := httptest.NewServer(apiHandler)
	client := gh.NewClient(nil)
	url, _ := url.Parse(server.URL + baseURLPath + "/")
	client.BaseURL = url
	client.UploadURL = url
	return server, client
}
