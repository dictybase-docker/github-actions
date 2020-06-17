package main

import (
	"log"
	"os"

	"github.com/dictyBase-docker/github-actions/internal/command/issue"
	"github.com/dictyBase-docker/github-actions/internal/command/storage"
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
			Usage: "format of the logging out, either of json or text.",
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
		{
			Name:    "issue-comment-report",
			Aliases: []string{"icr"},
			Usage:   "reports no of comments for every issue",
			Action:  issue.CommentsReport,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "output",
					Usage: "file where csv format output is written, creates a timestamp based file by default",
				},
				cli.StringFlag{
					Name:  "state",
					Usage: "state of the issue for filtering",
					Value: "all",
				},
			},
		},
		{
			Name:    "store-report",
			Aliases: []string{"ur"},
			Usage:   "save report to s3 storage",
			Action:  storage.SaveInS3,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "s3-server",
					Usage:  "S3 server endpoint",
					Value:  "minio",
					EnvVar: "MINIO_SERVICE_HOST",
				},
				cli.StringFlag{
					Name:   "s3-server-port",
					Usage:  "S3 server port",
					EnvVar: "MINIO_SERVICE_PORT",
				},
				cli.StringFlag{
					Name:  "s3-bucket",
					Usage: "S3 bucket where the data will be uploaded",
					Value: "report",
				},
				cli.StringFlag{
					Name:  "access-key, akey",
					Usage: "access key for S3 server, required based on command run",
				},
				cli.StringFlag{
					Name:  "secret-key, skey",
					Usage: "secret key for S3 server, required based on command run",
				},
				cli.StringFlag{
					Name:  "input",
					Usage: "input file that will be uploaded",
				},
				cli.StringFlag{
					Name:  "upload-path,p",
					Usage: "full upload path inside the bucket",
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("error in running command %s", err)
	}
}
