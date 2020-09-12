package chatops

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/dictyBase-docker/github-actions/internal/logger"
	"github.com/google/go-github/v32/github"
	"github.com/sethvargo/go-githubactions"
	"github.com/urfave/cli"
)

type pullRequestGetter interface {
	Get(ctx context.Context, owner string, repo string, number int) (*github.PullRequest, *github.Response, error)
}

type pullRequestClient struct {
	ctx               context.Context
	pullRequestClient pullRequestGetter
}

type branchGetter interface {
	GetBranch(ctx context.Context, owner string, repo string, branch string) (*github.Branch, *github.Response, error)
}

type branchClient struct {
	ctx          context.Context
	branchClient branchGetter
}

type Inputs struct {
	Cluster         string `json:"cluster"`
	URL             string `json:"html-url"`
	IssueNumber     string `json:"issue-number"`
	RepositoryName  string `json:"repository-name"`
	RepositoryOwner string `json:"repository-owner"`
	Commit          string `json:"commit"`
	Branch          string `json:"branch"`
}

// WorkflowDispatchEvent is triggered when someone triggers a workflow run on GitHub or
// sends a POST request to the create a workflow dispatch event endpoint.
//
// GitHub API docs: https://docs.github.com/en/developers/webhooks-and-events/webhook-events-and-payloads#workflow_dispatch
type WorkflowDispatchEvent struct {
	Inputs   json.RawMessage `json:"inputs,omitempty"`
	Ref      *string         `json:"ref,omitempty"`
	Workflow *string         `json:"workflow,omitempty"`

	// The following fields are only populated by Webhook events.
	Repo   *github.Repository   `json:"repository,omitempty"`
	Org    *github.Organization `json:"organization,omitempty"`
	Sender *github.User         `json:"sender,omitempty"`
}

type Output struct {
	ImageTag string
	Ref      string
}

func getWorkflowInputsFromJSON(r io.Reader) (*Inputs, error) {
	i := &Inputs{}
	w := &WorkflowDispatchEvent{}
	if err := json.NewDecoder(r).Decode(w); err != nil {
		return i, fmt.Errorf("error in decoding json %s", err)
	}
	if err := json.Unmarshal(w.Inputs, &i); err != nil {
		return i, fmt.Errorf("error in decoding json data to struct %s", err)
	}
	return i, nil
}

func ParseDeployCommand(c *cli.Context) error {
	r, err := os.Open(c.String("payload-file"))
	if err != nil {
		return fmt.Errorf("error in reading content from file %s", err)
	}
	defer r.Close()
	p, err := getWorkflowInputsFromJSON(r)
	if err != nil {
		return err
	}
	a := githubactions.New()
	log := logger.GetLogger(c)
	o, err := parseWorkflowInputs(p)
	if err != nil {
		return fmt.Errorf("error in parsing workflow inputs %s", err)
	}
	a.SetOutput("image_tag", o.ImageTag)
	a.SetOutput("ref", o.Ref)
	log.Info("added all keys to the output")
	return nil
}

func parseWorkflowInputs(p *Inputs) (*Output, error) {
	ou := &Output{}
	client := github.NewClient(nil)
	prc := &pullRequestClient{
		ctx:               context.Background(),
		pullRequestClient: client.PullRequests,
	}
	bc := &branchClient{
		ctx:          context.Background(),
		branchClient: client.Repositories,
	}
	if strings.Contains(p.URL, "pull") {
		o, err := parsePR(prc, p)
		if err != nil {
			return ou, err
		}
		return o, nil
	} else {
		o, err := parseIssue(bc, p)
		if err != nil {
			return ou, err
		}
		return o, nil
	}
}

func parsePR(prc *pullRequestClient, p *Inputs) (*Output, error) {
	o := &Output{}
	if p.Commit == "" {
		ref, err := prc.getHeadCommitFromPR(p.RepositoryName, p.RepositoryOwner, p.IssueNumber)
		if err != nil {
			return o, err
		}
		o.ImageTag = fmt.Sprintf("pr-%s-%s", p.IssueNumber, ref[0:7])
		o.Ref = ref
		return o, nil
	}
	o.ImageTag = fmt.Sprintf("pr-%s-%s", p.IssueNumber, p.Commit[0:7])
	o.Ref = p.Commit
	return o, nil
}

func (prc *pullRequestClient) getHeadCommitFromPR(name, owner, id string) (string, error) {
	num, err := strconv.Atoi(id)
	if err != nil {
		return "", fmt.Errorf("error converting string to int %s", err)
	}
	pr, _, err := prc.pullRequestClient.Get(context.Background(), owner, name, num)
	if err != nil {
		return "", fmt.Errorf("error getting pull request info %s", err)
	}
	return pr.GetHead().GetSHA(), nil
}

func parseIssue(bc *branchClient, p *Inputs) (*Output, error) {
	o := &Output{}
	if p.Branch != "" {
		ref, err := bc.getHeadCommitFromBranch(p.RepositoryName, p.RepositoryOwner, p.Branch)
		if err != nil {
			return o, err
		}
		cb := strings.ReplaceAll(p.Branch, "/", "-")
		o.ImageTag = fmt.Sprintf("%s-%s", cb, ref[0:7])
		o.Ref = ref
		return o, nil
	}
	if p.Commit != "" {
		o.ImageTag = p.Commit[0:7]
		o.Ref = p.Commit
		return o, nil
	}
	return o, nil
}

func (bc *branchClient) getHeadCommitFromBranch(name, owner, branch string) (string, error) {
	b, _, err := bc.branchClient.GetBranch(context.Background(), owner, name, branch)
	if err != nil {
		return "", fmt.Errorf("error getting pull request info %s", err)
	}
	return b.GetCommit().GetSHA(), nil
}
