package gcloud

import (
	"github.com/dictyBase-docker/github-actions/internal/runner"
	"github.com/urfave/cli"
)

func K8sClusterCredentials(c *cli.Context) error {
	gcloud, err := runner.NewGcloud()
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	err = gcloud.GetClusterCredentials(
		c.String("project"),
		c.String("zone"),
		c.String("cluster"),
	)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	return nil
}
