package gcloud

import (
	"github.com/dictyBase-docker/github-actions/internal/runner"
	"github.com/urfave/cli"
)

func K8sClusterCredentials(clt *cli.Context) error {
	gcloud, err := runner.NewGcloud()
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	err = gcloud.GetClusterCredentials(
		clt.String("project"),
		clt.String("zone"),
		clt.String("cluster"),
	)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}

	return nil
}
