package chatops

import (
	"fmt"

	"github.com/dictyBase-docker/github-actions/internal/logger"
	"github.com/sethvargo/go-githubactions"
	"github.com/urfave/cli"
)

func PullRequestDeploy(c *cli.Context) error {
	cluster := c.String("cluster")
	commit := c.String("commit")
	sha := c.String("sha")
	prid := c.String("pr-id")
	var ref, imageTag string

	a := githubactions.New()
	log := logger.GetLogger(c)

	if commit != "" {
		ref = commit
		imageTag = fmt.Sprintf("pr-%s-%s", prid, commit[0:6])
	} else {
		ref = sha
		imageTag = fmt.Sprintf("pr-%s-%s", prid, sha[0:6])
	}

	a.SetOutput("cluster", cluster)
	a.SetOutput("ref", ref)
	a.SetOutput("image_tag", imageTag)
	log.Info("added all keys to the output")

	return nil
}
