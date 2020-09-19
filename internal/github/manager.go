package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/dictyBase-docker/github-actions/internal/client"

	gh "github.com/google/go-github/v32/github"
	"github.com/urfave/cli"
)

type GithubManager struct {
	client *gh.Client
}

func NewGithubManager(c *gh.Client) *GithubManager {
	return &GithubManager{client: c}
}

func (g *GithubManager) CommittedFilesInPull(r io.Reader) (*ChangedFilesBuilder, error) {
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

func FilterCommittedFiles(c *cli.Context, in io.Reader, event string) ([]string, error) {
	gclient, err := client.GetGithubClient(c.GlobalString("token"))
	if err != nil {
		return []string{}, fmt.Errorf("error in getting github client %s", err)
	}
	var fb *ChangedFilesBuilder
	switch event {
	case "push":
		fb, err = NewGithubManager(gclient).CommittedFilesInPush(in)
	case "pull":
		fb, err = NewGithubManager(gclient).CommittedFilesInPull(in)
	default:
		return []string{}, fmt.Errorf("event type %s not supported", event)
	}
	if err != nil {
		return []string{}, err
	}
	return fb.FilterUniqueByName().
		FilterDeleted(c.BoolT("skip-deleted")).
		FilterSuffix(c.String("include-file-suffix")).
		List(), nil
}
