package cmd

import (
	"github.com/dictyBase-docker/github-actions/internal/app/storage"
	"github.com/urfave/cli"
)

func StoreReportCmd() cli.Command {
	return cli.Command{
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
	}
}
