package main

import (
	"os"

	"github.com/dictyBase-docker/github-actions/internal/command/issue"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Name = "gh-action"
	app.Usage = "run github action"
	app.Flags = []cli.Flag{
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
		{
			Name:    "issue-comment-report",
			Aliases: []string{"icr"},
			Usage:   "reports no of comments for every issue",
			Action:  issue.CommentsReport,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "output",
					Usage: "report output, goes to stdout by default",
				},
				cli.StringFlag{
					Name:  "state",
					Usage: "state of the issue for filtering",
					Value: "all",
				},
				cli.BoolTFlag{
					Name:  "attach",
					Usage: "Attach the output to a new issue or in a comment",
				},
			},
		},
	}
	app.Run(os.Args)
}
