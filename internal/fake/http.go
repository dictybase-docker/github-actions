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

type httpFunc func(string, http.ResponseWriter, *http.Request)

type route struct {
	file   string
	regexp *regexp.Regexp
	fn     httpFunc
}

func newRoute(pattern, file string, fn httpFunc) *route {
	return &route{
		regexp: regexp.MustCompile(fmt.Sprintf("^%s$", pattern)),
		fn:     fn,
		file:   file,
	}
}

func routeTable() []*route {
	return []*route{
		newRoute(
			fmt.Sprintf("%s%s",
				baseURLPath,
				`/repos/([^/]+)/([^/]+)`,
			),
			"edit-repo.json",
			handleSuccess,
		),
		newRoute(
			fmt.Sprintf("%s%s",
				baseURLPath,
				`/repos/([^/]+)/([^/]+)/forks`,
			),
			"create-fork.json",
			handleAccepted,
		),
		newRoute(
			fmt.Sprintf("%s%s",
				baseURLPath,
				`/repos/([^/]+)/([^/]+)/compare/\w+\.\.\.\w+`,
			),
			"commit-diff.json",
			handleSuccess,
		)}
}

func handleAccepted(file string, w http.ResponseWriter, r *http.Request) {
	b, err := payloadFile(file)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	if _, err := fmt.Fprint(w, string(b)); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
	}
}

func handleSuccess(file string, w http.ResponseWriter, r *http.Request) {
	b, err := payloadFile(file)
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
			rt.fn(rt.file, w, r)
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
