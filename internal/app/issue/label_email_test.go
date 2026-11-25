package issue

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/dictyBase-docker/github-actions/internal/fake"
	"github.com/google/go-github/v62/github"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

const (
	testDomain    = "sandbox.example.org"
	testAPIKey    = "key-test"
	testFromEmail = "dsc@example.org"
	maxFormMemory = 1 << 20
)

type mailgunRequest struct {
	path     string
	user     string
	password string
	form     url.Values
}

type fakeMailgun struct {
	*httptest.Server

	mu       sync.Mutex
	requests []mailgunRequest
}

func newFakeMailgun(t *testing.T, status int) *fakeMailgun {
	t.Helper()

	fmg := &fakeMailgun{}
	fmg.Server = httptest.NewServer(
		http.HandlerFunc(func(wrt http.ResponseWriter, req *http.Request) {
			req.Body = http.MaxBytesReader(wrt, req.Body, maxFormMemory)

			form, err := readMultipartForm(req)
			if err != nil {
				http.Error(wrt, err.Error(), http.StatusBadRequest)

				return
			}

			user, password, _ := req.BasicAuth()

			fmg.mu.Lock()
			fmg.requests = append(fmg.requests, mailgunRequest{
				path:     req.URL.Path,
				user:     user,
				password: password,
				form:     form,
			})
			fmg.mu.Unlock()

			wrt.Header().Set("Content-Type", "application/json")
			wrt.WriteHeader(status)
			fmt.Fprint(wrt, `{"id":"<test@sandbox.example.org>","message":"Queued. Thank you."}`)
		}),
	)
	t.Cleanup(fmg.Close)

	return fmg
}

func readMultipartForm(req *http.Request) (url.Values, error) {
	reader, err := req.MultipartReader()
	if err != nil {
		return nil, fmt.Errorf("expected multipart payload: %w", err)
	}

	form := url.Values{}

	for {
		part, err := reader.NextPart()
		if errors.Is(err, io.EOF) {
			return form, nil
		}

		if err != nil {
			return nil, fmt.Errorf("reading multipart part: %w", err)
		}

		value, err := io.ReadAll(part)
		if err != nil {
			return nil, fmt.Errorf("reading multipart value: %w", err)
		}

		form.Add(part.FormName(), string(value))
	}
}

func (fmg *fakeMailgun) received() []mailgunRequest {
	fmg.mu.Lock()
	defer fmg.mu.Unlock()

	return append([]mailgunRequest(nil), fmg.requests...)
}

func labelEmailContext(t *testing.T, apiBase, apiKey string) *cli.Context {
	t.Helper()

	set := flag.NewFlagSet("test", flag.ContinueOnError)
	set.String("token", "gh-token", "")
	set.String("owner", "dictybase-playground", "")
	set.String("repository", "learn-github-action", "")
	set.Int("issueid", 122, "")
	set.String("label", shippedLabel, "")
	set.String("domain", testDomain, "")
	set.String("fromEmail", testFromEmail, "")
	set.String("apiKey", apiKey, "")
	set.String("apiBase", apiBase, "")
	require.NoError(t, set.Parse(nil))

	return cli.NewContext(cli.NewApp(), set, nil)
}

func fakeGithubFactory(t *testing.T) (githubClientFactory, *bool) {
	t.Helper()

	server, gclient := fake.GhServerClient()
	t.Cleanup(server.Close)

	called := false

	return func(token string) (*github.Client, error) {
		called = true

		require.Equal(t, "gh-token", token)

		return gclient, nil
	}, &called
}

func TestSendIssueLabelEmailSendsOrderUpdate(t *testing.T) {
	t.Parallel()

	mailgun := newFakeMailgun(t, http.StatusOK)
	newGithubClient, _ := fakeGithubFactory(t)

	err := sendIssueLabelEmail(labelEmailContext(t, mailgun.URL, testAPIKey), newGithubClient)
	require.NoError(t, err)

	requests := mailgun.received()
	require.Len(t, requests, 1)

	sent := requests[0]
	require.Equal(t, "/v3/"+testDomain+"/messages", sent.path)
	require.Equal(t, "api", sent.user)
	require.Equal(t, testAPIKey, sent.password)
	require.Equal(t, testFromEmail, sent.form.Get("from"))
	require.Equal(t, []string{"kevin.tun@northwestern.edu"}, sent.form["to"])
	require.Equal(t, "Dicty Stock Center - Order Update #10283618", sent.form.Get("subject"))
	require.Contains(t, sent.form.Get("html"), "10283618")
	require.Contains(t, sent.form.Get("html"), shippedLabel)
}

func TestSendIssueLabelEmailMailgunFailure(t *testing.T) {
	t.Parallel()

	mailgun := newFakeMailgun(t, http.StatusUnauthorized)
	newGithubClient, _ := fakeGithubFactory(t)

	err := sendIssueLabelEmail(labelEmailContext(t, mailgun.URL, testAPIKey), newGithubClient)
	require.ErrorContains(t, err, "error sending email")

	var exitErr cli.ExitCoder

	require.ErrorAs(t, err, &exitErr)
	require.Equal(t, exitFailure, exitErr.ExitCode())
	require.Len(t, mailgun.received(), 1)
}

func TestSendIssueLabelEmailRequiresMailgunConfig(t *testing.T) {
	t.Parallel()

	mailgun := newFakeMailgun(t, http.StatusOK)
	newGithubClient, githubCalled := fakeGithubFactory(t)

	err := sendIssueLabelEmail(labelEmailContext(t, mailgun.URL, ""), newGithubClient)
	require.ErrorContains(t, err, "Mailgun domain and apiKey are required")
	require.False(t, *githubCalled, "GitHub must not be queried without Mailgun config")
	require.Empty(t, mailgun.received())
}
