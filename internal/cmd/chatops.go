package cmd

import (
	"github.com/dictyBase-docker/github-actions/internal/app/chatops"
	"github.com/urfave/cli"
)

const payloadFileFlag = "payload-file,f"

func ParseChatOpsDeploy() cli.Command {
	return cli.Command{
		Name:    "parse-chatops-deploy",
		Aliases: []string{"pcd"},
		Usage:   "parses chatops deploy command and extracts ref and image tag values",
		Action:  chatops.ParseDeployCommand,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     payloadFileFlag,
				Required: true,
				Usage:    "path to JSON payload",
			},
			cli.BoolFlag{
				Name:  "frontend",
				Usage: "used if deploying frontend web app (needed for updating image-tag correctly)",
			},
		},
	}
}
