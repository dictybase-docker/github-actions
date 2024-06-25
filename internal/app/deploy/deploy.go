package deploy

import (
	"context"
	"fmt"

	"github.com/dictyBase-docker/github-actions/internal/client"
	"github.com/dictyBase-docker/github-actions/internal/logger"
	"github.com/google/go-github/v62/github"
	"github.com/urfave/cli"
)

func Status(clt *cli.Context) error {
	logger := logger.GetLogger(clt)
	gclient, err := client.GetGithubClient(clt.GlobalString("token"))
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in getting github client %s", err),
			2,
		)
	}
	state := clt.String("state")
	url := clt.String("url")
	desc := fmt.Sprintf("setting deployment status %s", clt.String("state"))
	dsp, _, err := gclient.Repositories.CreateDeploymentStatus(
		context.Background(),
		clt.GlobalString("owner"),
		clt.GlobalString("repository"),
		clt.Int64("deployment_id"),
		&github.DeploymentStatusRequest{
			State:       &state,
			LogURL:      &url,
			Description: &desc,
		})
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf(
				"error in creating deployment status %s %s",
				state,
				err,
			),
			2,
		)
	}
	logger.Infof(
		"created deployment status %s with id %d",
		dsp.GetState(),
		dsp.GetID(),
	)

	return nil
}
