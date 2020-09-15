package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	gh "github.com/google/go-github/v32/github"
)

type GithubManager struct {
	client *gh.Client
}

func NewGithubManager(c *gh.Client) *GithubManager {
	return &GithubManager{client: c}
}

func (g *GithubManager) CommitedFilesInPush(r io.Reader) ([]string, error) {
	var files []string
	pe := &gh.PushEvent{}
	if err := json.NewDecoder(r).Decode(pe); err != nil {
		return files, fmt.Errorf("error in decoding json %s", err)
	}
	comc, _, err := g.client.Repositories.CompareCommits(
		context.Background(),
		pe.GetRepo().GetOwner().GetLogin(),
		pe.GetRepo().GetName(),
		pe.GetBefore(),
		pe.GetAfter(),
	)
	if err != nil {
		return files, fmt.Errorf("error in comparing commits %s", err)
	}
	return CommittedFiles(comc), nil
}
