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

func (g *GithubManager) CommitedFilesInPull(r io.Reader) (*ChangedFilesBuilder, error) {
	var b *ChangedFilesBuilder
	pe := &gh.PullRequestEvent{}
	if err := json.NewDecoder(r).Decode(pe); err != nil {
		return b, fmt.Errorf("error in decoding json %s", err)
	}
	var after, before string
	switch pe.GetAction() {
	case "synchronize":
		before = pe.GetBefore()
		after = pe.GetAfter()
	case "opened":
		before = pe.GetPullRequest().GetBase().GetSHA()
		after = pe.GetPullRequest().GetHead().GetSHA()
	}
	comc, _, err := g.client.Repositories.CompareCommits(
		context.Background(),
		pe.GetRepo().GetOwner().GetLogin(),
		pe.GetRepo().GetName(),
		before,
		after,
	)
	if err != nil {
		return b, fmt.Errorf("error in comparing commits %s", err)
	}
	return CommittedFiles(comc), nil
}

func (g *GithubManager) CommittedFilesInPush(r io.Reader) (*ChangedFilesBuilder, error) {
	var b *ChangedFilesBuilder
	pe := &gh.PushEvent{}
	if err := json.NewDecoder(r).Decode(pe); err != nil {
		return b, fmt.Errorf("error in decoding json %s", err)
	}
	comc, _, err := g.client.Repositories.CompareCommits(
		context.Background(),
		pe.GetRepo().GetOwner().GetLogin(),
		pe.GetRepo().GetName(),
		pe.GetBefore(),
		pe.GetAfter(),
	)
	if err != nil {
		return b, fmt.Errorf("error in comparing commits %s", err)
	}
	return CommittedFiles(comc), nil
}

func CommittedFiles(event *gh.CommitsComparison) *ChangedFilesBuilder {
	var a []*ChangedFiles
	for _, f := range event.Files {
		a = append(
			a,
			&ChangedFiles{Name: f.GetFilename(), Change: f.GetStatus()},
		)
	}
	return &ChangedFilesBuilder{files: a}
}
