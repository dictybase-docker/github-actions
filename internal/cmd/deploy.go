package cmd

import (
	"github.com/dictyBase-docker/github-actions/internal/app/deploy"
	"github.com/urfave/cli"
)

func DeployStatusCmd() cli.Command {
	return cli.Command{
		Name:    "deploy-status",
		Aliases: []string{"ds"},
		Usage:   "create a github deployment status",
		Action:  deploy.DeployStatus,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "state",
				Usage: "The state of the deployment status",
			},
			cli.Int64Flag{
				Name:  "deployment_id",
				Usage: "Deployment identifier",
			},
			cli.StringFlag{
				Name:  "url",
				Usage: "The url that is associated with this status",
			},
		},
	}
}
