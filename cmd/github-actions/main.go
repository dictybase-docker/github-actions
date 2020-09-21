package main

import (
	"log"
	"os"

	"github.com/dictyBase-docker/github-actions/internal/cmd"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Name = "gh-action"
	app.Usage = "run github action"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "log-format",
			Usage: "format of the log, either of json or text.",
			Value: "json",
		},
		cli.StringFlag{
			Name:  "log-level",
			Usage: "log level for the application",
			Value: "error",
		},
		cli.StringFlag{
			Name:   "token,t",
			Usage:  "github personal access token",
			EnvVar: "GITHUB_TOKEN",
		},
		cli.StringFlag{
			Name:  "repository, r",
			Usage: "Github repository",
		},
		cli.StringFlag{
			Name:  "owner",
			Usage: "Github repository owner",
			Value: "dictyBase",
		},
	}
	app.Commands = []cli.Command{
		cmd.IssueCommentCmds(),
		cmd.StoreReportCmd(),
		cmd.DeployStatusCmd(),
		cmd.ShareDeployPayloadCmd(),
		cmd.GetK8sClusterCredentialsCmd(),
		cmd.GenerateMkdownCmd(),
		cmd.DeployChartCmd(),
		cmd.FilesCommited(),
		cmd.BatchMultiRepo(),
		cmd.ParseChatOpsDeploy(),
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("error in running command %s", err)
	}
}
