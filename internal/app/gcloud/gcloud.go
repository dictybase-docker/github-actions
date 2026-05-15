package gcloud

import (
	"github.com/dictyBase-docker/github-actions/internal/runner"
	"github.com/urfave/cli"
)

const exitFailure = 2

func K8sClusterCredentials(clt *cli.Context) error {
	gcloud, err := runner.NewGcloud()
	if err != nil {
		return cli.NewExitError(err.Error(), exitFailure)
	}

	err = gcloud.GetClusterCredentials(
		clt.String("project"),
		clt.String("zone"),
		clt.String("cluster"),
	)
	if err != nil {
		return cli.NewExitError(err.Error(), exitFailure)
	}

	return nil
}
