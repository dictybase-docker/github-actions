package client

import (
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

func GetGithubClient(token string) (*github.Client, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	return github.NewClient(tc), nil
}
