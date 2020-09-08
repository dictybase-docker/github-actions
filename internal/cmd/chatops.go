package cmd

import (
	"github.com/dictyBase-docker/github-actions/internal/app/chatops"
	"github.com/urfave/cli"
)

func ChatOpsPullRequestDeploy() cli.Command {
	return cli.Command{
		Name:    "chatops-pr-deploy",
		Aliases: []string{"cpd"},
		Usage:   "extracts necessary values from chatops deploy commands in pull requests",
		Action:  chatops.PullRequestDeploy,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "cluster",
				Required: true,
				Usage:    "k8s cluster to deploy to",
			},
			cli.StringFlag{
				Name:  "commit",
				Usage: "commit id to deploy",
			},
			cli.StringFlag{
				Name:  "pr-id",
				Usage: "id (number) of given pull request",
			},
			cli.StringFlag{
				Name:  "head-sha",
				Usage: "head commit id of the pull request",
			},
		},
	}
}
