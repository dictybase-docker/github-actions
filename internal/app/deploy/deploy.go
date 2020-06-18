package deploy

import (
	"context"
	"fmt"

	"github.com/dictyBase-docker/github-actions/internal/client"
	"github.com/google/go-github/v32/github"
	"github.com/urfave/cli"
)

func DeployStatus(c *cli.Context) error {
	gclient, err := client.GetGithubClient(c.GlobalString("token"))
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in getting github client %s", err),
			2,
		)
	}
	state := c.String("state")
	url := c.String("url")
	desc := fmt.Sprintf("setting deployment status %s", c.String("state"))
	_, _, err = gclient.Repositories.CreateDeploymentStatus(
		context.Background(),
		c.GlobalString("owner"),
		c.GlobalString("repository"),
		c.Int64("deployment_id"),
		&github.DeploymentStatusRequest{
			State:       &state,
			LogURL:      &url,
			Description: &desc,
		})
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in creating deployment status %s %s", state, err),
			2,
		)
	}
	return nil
}
