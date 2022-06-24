package cmd

import (
	"github.com/dictyBase-docker/github-actions/internal/app/analytics"
	"github.com/urfave/cli"
)

func AnalyticsReportCmd() cli.Command {
	return cli.Command{
		Name:    "analytics-report",
		Usage:   "generate google analytics report of sessions,users,pageviews in csv format.",
		Aliases: []string{"ar"},
		Action:  analytics.Report,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "view-id",
				Usage:    "view id of the project",
				Required: true,
			},
			cli.StringFlag{
				Name:     "credential-file, c",
				Usage:    "credential for google service account",
				Required: true,
			},
			cli.StringFlag{
				Name:     "start-date,s",
				Usage:    "start date for the query, should be in YYYY-MM-DD format",
				Required: true,
			},
			cli.StringFlag{
				Name:  "end-date,e",
				Usage: "end date for the query, should be in YYYY-MM-DD format, current date is assumed by default",
			},
			cli.StringFlag{
				Name:  "output,o",
				Usage: "output csv file name",
			},
		},
	}
}
