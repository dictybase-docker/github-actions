package fake

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	gh "github.com/google/go-github/v62/github"
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

func fetchRoute() []*route {
	return []*route{
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

func postRoute() []*route {
	return []*route{
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
	}
}

func delRoute() []*route {
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
	}
}

func routeTable() []*route {
	var route []*route
	route = append(route, fetchRoute()...)
	route = append(route, postRoute()...)

	return append(route, delRoute()...)
}

func handleNoContent(file string, w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
	if _, err := fmt.Fprint(w, file); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
	}
}

func handleAccepted(file string, wrt http.ResponseWriter, _ *http.Request) {
	bfl, err := payloadFile(file)
	if err != nil {
		http.Error(
			wrt,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}
	wrt.WriteHeader(http.StatusAccepted)
	fmt.Fprint(wrt, string(bfl))
}

func handleSuccess(file string, wrt http.ResponseWriter, _ *http.Request) {
	bfl, err := payloadFile(file)
	if err != nil {
		http.Error(
			wrt,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}
	if _, err := wrt.Write(bfl); err != nil {
		http.Error(
			wrt,
			err.Error(),
			http.StatusInternalServerError,
		)
	}
}

func router(wrt http.ResponseWriter, req *http.Request) {
	for _, rtbl := range routeTable() {
		if req.Method != rtbl.method {
			continue
		}
		if !rtbl.regexp.MatchString(req.URL.Path) {
			continue
		}
		rtbl.fn(rtbl.file, wrt, req)

		return
	}
	http.NotFound(wrt, req)
}

func payloadFile(file string) ([]byte, error) {
	path, err := filepath.Abs(file)
	if err != nil {
		return []byte(""), fmt.Errorf("unable to get current file %s", err)
	}
	b, err := os.ReadFile(path)
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
