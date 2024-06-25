package client

import (
	"context"

	lgh "github.com/google/go-github/v32/github"
	"github.com/google/go-github/v62/github"
	"golang.org/x/oauth2"
)

func GetLegacyGithubClient(token string) (*lgh.Client, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	return lgh.NewClient(tc), nil
}

func GetGithubClient(token string) (*github.Client, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	return github.NewClient(tc), nil
}
