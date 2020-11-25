package fake

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	gh "github.com/google/go-github/v32/github"
)

const (
	baseURLPath = "/api-v3"
)

type route struct {
	regexp *regexp.Regexp
	fn     http.HandlerFunc
}

func newRoute(pattern string, fn http.HandlerFunc) *route {
	return &route{
		regexp: regexp.MustCompile(fmt.Sprintf("^%s$", pattern)),
		fn:     fn,
	}
}

func routeTable() []*route {
	return []*route{
		newRoute(
			fmt.Sprintf("%s%s",
				baseURLPath,
				`/repos/([^/]+)/([^/]+)/compare/\w+\.\.\.\w+`,
			),
			handleCommitComparison,
		)}
}

func handleCommitComparison(w http.ResponseWriter, r *http.Request) {
	b, err := payloadFile("commit-diff.json")
	if err != nil {
		http.Error(
			w,
			err.Error(),
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

func router(w http.ResponseWriter, r *http.Request) {
	for _, rt := range routeTable() {
		if rt.regexp.MatchString(r.URL.Path) {
			rt.fn(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

func payloadFile(file string) ([]byte, error) {
	dir, err := os.Getwd()
	if err != nil {
		return []byte(""), fmt.Errorf("unable to get current dir %s", err)
	}
	path := filepath.Join(
		filepath.Dir(dir), "../testdata", file,
	)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte(""), fmt.Errorf("unable to read test file %s", path)
	}
	return b, nil
}

func GhServerClient() (*httptest.Server, *gh.Client) {
	apiHandler := http.NewServeMux()
	apiHandler.HandleFunc("/", router)
	server := httptest.NewServer(apiHandler)
	client := gh.NewClient(nil)
	url, _ := url.Parse(server.URL + baseURLPath + "/")
	client.BaseURL = url
	client.UploadURL = url
	return server, client
}
