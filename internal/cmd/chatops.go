package cmd

import (
	"github.com/dictyBase-docker/github-actions/internal/app/chatops"
	"github.com/urfave/cli"
)

func ParseChatOpsDeploy() cli.Command {
	return cli.Command{
		Name:    "parse-chatops-deploy",
		Aliases: []string{"pcd"},
		Usage:   "extracts necessary values from chatops deploy commands and converts to expected outputs",
		Action:  chatops.ParseDeployCommand,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "payload-file,f",
				Required: true,
				Usage:    "Path to JSON payload",
			},
		},
	}
}
