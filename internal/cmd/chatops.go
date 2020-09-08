package cmd

import (
	"github.com/dictyBase-docker/github-actions/internal/app/chatops"
	"github.com/urfave/cli"
)

func ShareChatOpsPayload() cli.Command {
	return cli.Command{
		Name:    "share-chatops-payload",
		Aliases: []string{"scp"},
		Usage:   "share chatops payload data in github workflow",
		Action:  chatops.ShareChatOpsPayload,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "payload-file,f",
				Required: true,
				Usage:    "Full path to the file that contains the chatops payload",
			},
			cli.StringFlag{
				Name:     "cluster,c",
				Required: true,
				Usage:    "k8s cluster to deploy to",
			},
		},
	}
}
