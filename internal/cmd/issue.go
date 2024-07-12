package cmd

import (
	"github.com/dictyBase-docker/github-actions/internal/app/issue"
	"github.com/urfave/cli"
)

func CommentsCountByDateCmds() cli.Command {
	return cli.Command{
		Name:    "issue-comment-count",
		Aliases: []string{"icc"},
		Usage:   "reports total no of issues and comments since a particular date",
		Action:  issue.CommentsCountByDate,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "since",
				Usage: "date value should be given in YYYY-MM-DD format",
			},
		},
	}
}

func IssueCommentCmds() cli.Command {
	return cli.Command{
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
	}
}
