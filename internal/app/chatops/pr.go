package chatops

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dictyBase-docker/github-actions/internal/logger"
	"github.com/google/go-github/v32/github"
	"github.com/sethvargo/go-githubactions"
	"github.com/urfave/cli"
)

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

type Payload struct {
	Event WorkflowDispatchEvent `json:"event"`
}

func GetWorkflowInputs(data []byte) (*Inputs, error) {
	o := new(Inputs)
	if err := json.Unmarshal(data, &o); err != nil {
		return o, fmt.Errorf("error in decoding json data to string %s", err)
	}
	return o, nil
}

func ParseDeployCommand(c *cli.Context) error {
	r, err := os.Open(c.String("payload-file"))
	if err != nil {
		return fmt.Errorf("error in reading content from file %s", err)
	}
	defer r.Close()
	i := &Payload{}
	if err := json.NewDecoder(r).Decode(i); err != nil {
		return fmt.Errorf("error in decoding json %s", err)
	}
	p, err := GetWorkflowInputs(i.Event.Inputs)
	if err != nil {
		return err
	}
	a := githubactions.New()
	log := logger.GetLogger(c)
	if strings.Contains(p.URL, "pull") && p.Commit != "" {
		tag := fmt.Sprintf("pr-%s-%s", p.IssueNumber, p.Commit[0:7])
		a.SetOutput("image_tag", tag)
		a.SetOutput("ref", p.Commit)
	}
	if strings.Contains(p.URL, "pull") && p.Commit == "" {
		ref, err := getHeadCommitFromPR(p.RepositoryName, p.RepositoryOwner, p.IssueNumber)
		if err != nil {
			return err
		}
		tag := fmt.Sprintf("pr-%s-%s", p.IssueNumber, ref[0:7])
		a.SetOutput("ref", ref)
		a.SetOutput("image_tag", tag)
	}
	log.Info("added all keys to the output")
	return nil
}

func getHeadCommitFromPR(name, owner, id string) (string, error) {
	client := github.NewClient(nil)
	num, err := strconv.Atoi(id)
	if err != nil {
		return "", fmt.Errorf("error converting string to int %s", err)
	}
	pr, _, err := client.PullRequests.Get(context.Background(), owner, name, num)
	if err != nil {
		return "", err
	}
	return *pr.Head.SHA, nil
}