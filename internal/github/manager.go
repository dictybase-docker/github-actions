package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	gh "github.com/google/go-github/v62/github"
)

type CommittedFilesParams struct {
	Client      *gh.Client
	Input       io.Reader
	Event       string
	FileSuffix  string
	SkipDeleted bool
}

type Manager struct {
	client *gh.Client
}

func NewGithubManager(c *gh.Client) *Manager {
	return &Manager{client: c}
}

func (g *Manager) CommittedFilesInPull(
	r io.Reader,
) (*ChangedFilesBuilder, error) {
	var bcf *ChangedFilesBuilder
	pev := &gh.PullRequestEvent{}
	if err := json.NewDecoder(r).Decode(pev); err != nil {
		return bcf, fmt.Errorf("error in decoding json %s", err)
	}
	var after, before string
	switch pev.GetAction() {
	case "synchronize":
		before = pev.GetBefore()
		after = pev.GetAfter()
	case "opened":
		before = pev.GetPullRequest().GetBase().GetSHA()
		after = pev.GetPullRequest().GetHead().GetSHA()
	}
	comc, _, err := g.client.Repositories.CompareCommits(
		context.Background(),
		pev.GetRepo().GetOwner().GetLogin(),
		pev.GetRepo().GetName(),
		before,
		after,
	)
	if err != nil {
		return bcf, fmt.Errorf("error in comparing commits %s", err)
	}

	return CommittedFiles(comc), nil
}

func (g *Manager) CommittedFilesInPush(
	r io.Reader,
) (*ChangedFilesBuilder, error) {
	var bfl *ChangedFilesBuilder
	pev := &gh.PushEvent{}
	if err := json.NewDecoder(r).Decode(pev); err != nil {
		return bfl, fmt.Errorf("error in decoding json %s", err)
	}
	comc, _, err := g.client.Repositories.CompareCommits(
		context.Background(),
		pev.GetRepo().GetOwner().GetLogin(),
		pev.GetRepo().GetName(),
		pev.GetBefore(),
		pev.GetAfter(),
	)
	if err != nil {
		return bfl, fmt.Errorf("error in comparing commits %s", err)
	}

	return CommittedFiles(comc), nil
}

func CommittedFiles(event *gh.CommitsComparison) *ChangedFilesBuilder {
	afc := make([]*ChangedFiles, 0)
	for _, f := range event.Files {
		afc = append(
			afc,
			&ChangedFiles{Name: f.GetFilename(), Change: f.GetStatus()},
		)
	}

	return &ChangedFilesBuilder{files: afc}
}

func FilterCommittedFiles(args *CommittedFilesParams) ([]string, error) {
	var fbl *ChangedFilesBuilder
	var err error
	switch args.Event {
	case "push":
		fbl, err = NewGithubManager(
			args.Client,
		).CommittedFilesInPush(args.Input)
	case "pull-request":
		fbl, err = NewGithubManager(
			args.Client,
		).CommittedFilesInPull(args.Input)
	default:
		err = fmt.Errorf("event type %s not supported", args.Event)
	}
	if err != nil {
		return []string{}, err
	}

	return fbl.FilterUniqueByName().
		FilterDeleted(args.SkipDeleted).
		FilterSuffix(args.FileSuffix).
		List(), nil
}
