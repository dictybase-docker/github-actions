package chatops

import (
	"fmt"

	"github.com/dictyBase-docker/github-actions/internal/logger"
	"github.com/sethvargo/go-githubactions"
	"github.com/urfave/cli"
)

type Input struct {
	Commit        string
	SHA           string
	PullRequestID string
}

type Output struct {
	ImageTag string
	Ref      string
}

func PullRequestDeploy(c *cli.Context) error {
	i := &Input{
		Commit:        c.String("commit"),
		SHA:           c.String("sha"),
		PullRequestID: c.String("pr-id"),
	}
	a := githubactions.New()
	log := logger.GetLogger(c)

	o := convertToOutput(i)

	a.SetOutput("cluster", c.String("cluster"))
	a.SetOutput("ref", o.Ref)
	a.SetOutput("image_tag", o.ImageTag)
	log.Info("added all keys to the output")

	return nil
}

func convertToOutput(i *Input) *Output {
	commit := i.Commit
	sha := i.SHA
	prid := i.PullRequestID
	var ref, imageTag string
	if commit != "" {
		ref = commit
		imageTag = fmt.Sprintf("pr-%s-%s", prid, commit[0:7])
	} else {
		ref = sha
		imageTag = fmt.Sprintf("pr-%s-%s", prid, sha[0:7])
	}
	return &Output{
		ImageTag: imageTag,
		Ref:      ref,
	}
}
