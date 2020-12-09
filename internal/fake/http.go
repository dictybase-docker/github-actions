package fake

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	method string
	regexp *regexp.Regexp
	fn     httpFunc
}

func routeTable() []*route {
	return []*route{
		{
			method: "DELETE",
			file:   "",
			fn:     handleNoContent,
			regexp: regexp.MustCompile(
				fmt.Sprintf(
					"^%s%s$",
					baseURLPath,
					`/repos/([^/]+)/([^/]+)`,
				)),
		},
		{
			method: "PATCH",
			file:   "../../../testdata/edit-repo.json",
			fn:     handleSuccess,
			regexp: regexp.MustCompile(
				fmt.Sprintf(
					"^%s%s$",
					baseURLPath,
					`/repos/([^/]+)/([^/]+)`,
				)),
		},
		{
			file:   "../../../testdata/create-fork.json",
			fn:     handleAccepted,
			method: "POST",
			regexp: regexp.MustCompile(
				fmt.Sprintf(
					"^%s%s$",
					baseURLPath,
					`/repos/([^/]+)/([^/]+)/forks`,
				)),
		},
		{
			method: "GET",
			file:   "../../../testdata/get-repo.json",
			fn:     handleSuccess,
			regexp: regexp.MustCompile(
				fmt.Sprintf(
					"^%s%s$",
					baseURLPath,
					`/repos/([^/]+)/([^/]+)`,
				)),
		},
		{
			file:   "../../testdata/commit-diff.json",
			fn:     handleSuccess,
			method: "GET",
			regexp: regexp.MustCompile(
				fmt.Sprintf(
					"^%s%s$",
					baseURLPath,
					`/repos/([^/]+)/([^/]+)/compare/\w+\.\.\.\w+`,
				)),
		},
	}
}

func handleNoContent(file string, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
	if _, err := fmt.Fprint(w, file); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
	}
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
	fmt.Fprint(w, string(b))
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
		if r.Method != rt.method {
			continue
		}
		if !rt.regexp.MatchString(r.URL.Path) {
			continue
		}
		rt.fn(rt.file, w, r)
		return
	}
	http.NotFound(w, r)
}

func payloadFile(file string) ([]byte, error) {
	path, err := filepath.Abs(file)
	if err != nil {
		return []byte(""), fmt.Errorf("unable to get current file %s", err)
	}
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
