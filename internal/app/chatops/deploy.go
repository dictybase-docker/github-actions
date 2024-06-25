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
	Get(
		ctx context.Context,
		owner string,
		repo string,
		number int,
	) (*github.PullRequest, *github.Response, error)
}

type pullRequestClient struct {
	ctx               context.Context
	pullRequestClient pullRequestGetter
}

type branchGetter interface {
	GetBranch(
		ctx context.Context,
		owner string,
		repo string,
		branch string,
	) (*github.Branch, *github.Response, error)
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
	inp := &Inputs{}
	w := &WorkflowDispatchEvent{}
	if err := json.NewDecoder(r).Decode(w); err != nil {
		return inp, fmt.Errorf("error in decoding json %s", err)
	}
	if err := json.Unmarshal(w.Inputs, &inp); err != nil {
		return inp, fmt.Errorf("error in decoding json data to struct %s", err)
	}

	return inp, nil
}

func ParseDeployCommand(clt *cli.Context) error {
	r, err := os.Open(clt.String("payload-file"))
	if err != nil {
		return fmt.Errorf("error in reading content from file %s", err)
	}
	defer r.Close()
	pjson, err := getWorkflowInputsFromJSON(r)
	if err != nil {
		return err
	}
	act := githubactions.New()
	log := logger.GetLogger(clt)
	oinput, err := parseWorkflowInputs(pjson)
	if err != nil {
		return fmt.Errorf("error in parsing workflow inputs %s", err)
	}
	imageTag := oinput.ImageTag
	// add image tag prefixes for developers
	if clt.Bool("frontend") && pjson.Cluster == "erickube" {
		imageTag = fmt.Sprintf("ericdev-%s", oinput.ImageTag)
	}
	if clt.Bool("frontend") && pjson.Cluster == "siddkube" {
		imageTag = fmt.Sprintf("devsidd-%s", oinput.ImageTag)
	}
	act.SetOutput("image_tag", imageTag)
	act.SetOutput("ref", oinput.Ref)
	log.Info("added all keys to the output")

	return nil
}

func parseWorkflowInputs(param *Inputs) (*Output, error) {
	out := &Output{}
	client := github.NewClient(nil)
	prc := &pullRequestClient{
		ctx:               context.Background(),
		pullRequestClient: client.PullRequests,
	}
	bclient := &branchClient{
		ctx:          context.Background(),
		branchClient: client.Repositories,
	}
	if strings.Contains(param.URL, "pull") {
		o, err := parsePR(prc, param)
		if err != nil {
			return out, err
		}

		return o, nil
	}
	o, err := parseIssue(bclient, param)
	if err != nil {
		return out, err
	}

	return o, nil
}

func parsePR(prc *pullRequestClient, param *Inputs) (*Output, error) {
	out := &Output{}
	if param.Commit == "" {
		ref, err := prc.getHeadCommitFromPR(
			param.RepositoryName,
			param.RepositoryOwner,
			param.IssueNumber,
		)
		if err != nil {
			return out, err
		}
		out.ImageTag = fmt.Sprintf("pr-%s-%s", param.IssueNumber, ref[0:7])
		out.Ref = ref

		return out, nil
	}
	out.ImageTag = fmt.Sprintf("pr-%s-%s", param.IssueNumber, param.Commit[0:7])
	out.Ref = param.Commit

	return out, nil
}

func (prc *pullRequestClient) getHeadCommitFromPR(
	name, owner, id string,
) (string, error) {
	num, err := strconv.Atoi(id)
	if err != nil {
		return "", fmt.Errorf("error converting string to int %s", err)
	}
	pgr, _, err := prc.pullRequestClient.Get(
		context.Background(),
		owner,
		name,
		num,
	)
	if err != nil {
		return "", fmt.Errorf("error getting pull request info %s", err)
	}

	return pgr.GetHead().GetSHA(), nil
}

func parseIssue(bc *branchClient, param *Inputs) (*Output, error) {
	out := &Output{}
	if param.Branch != "" {
		ref, err := bc.getHeadCommitFromBranch(
			param.RepositoryName,
			param.RepositoryOwner,
			param.Branch,
		)
		if err != nil {
			return out, err
		}
		cb := strings.ReplaceAll(param.Branch, "/", "-")
		out.ImageTag = fmt.Sprintf("%s-%s", cb, ref[0:7])
		out.Ref = ref

		return out, nil
	}
	if param.Commit != "" {
		out.ImageTag = param.Commit[0:7]
		out.Ref = param.Commit

		return out, nil
	}

	return out, nil
}

func (bc *branchClient) getHeadCommitFromBranch(
	name, owner, branch string,
) (string, error) {
	brch, _, err := bc.branchClient.GetBranch(
		context.Background(),
		owner,
		name,
		branch,
	)
	if err != nil {
		return "", fmt.Errorf("error getting pull request info %s", err)
	}

	return brch.GetCommit().GetSHA(), nil
}
