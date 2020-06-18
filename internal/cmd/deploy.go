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
				Name:     "state",
				Required: true,
				Usage:    "The state of the deployment status",
			},
			cli.Int64Flag{
				Name:     "deployment_id",
				Required: true,
				Usage:    "Deployment identifier",
			},
			cli.StringFlag{
				Name:     "url",
				Required: true,
				Usage:    "The url that is associated with this status",
			},
		},
	}
}

func ShareDeployPayloadCmd() cli.Command {
	return cli.Command{
		Name:    "share-deploy-payload",
		Aliases: []string{"sdp"},
		Usage:   "share deployment payload data in github workflow",
		Action:  deploy.ShareDeployPayload,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "payload-file,f",
				Required: true,
				Usage:    "Full path to the file that contain the deploy payload",
			},
		},
	}
}
