package chatops

import (
	"github.com/dictyBase-docker/github-actions/internal/logger"
	"github.com/sethvargo/go-githubactions"
	"github.com/urfave/cli"
)

func PullRequestDeploy(c *cli.Context) error {
	a := githubactions.New()
	log := logger.GetLogger(c)
	a.SetOutput("cluster", c.String("cluster"))
	a.SetOutput("ref", "")
	a.SetOutput("image_tag", "")
	log.Info("added all keys to the output")
	return nil
}
