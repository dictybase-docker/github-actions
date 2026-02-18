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

func IssueLabelEmailCmds() cli.Command {
	return cli.Command{
		Name:    "issue-label-email",
		Aliases: []string{"ile"},
		Usage:   "sends an email to a recipient of an order when certain labels are added to the issue",
		Action:  issue.SendIssueLabelEmail,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "apiKey",
				Usage: "API key for mailgun",
			},
			cli.StringFlag{
				Name:  "label",
				Usage: "The label that was added to the issue",
			},
			cli.IntFlag{
				Name:  "issueid",
				Usage: "The id of the issue",
			},
			cli.StringFlag{
				Name:  "domain",
				Usage: "Domain of mailgun endpoint",
			},
			cli.StringFlag{
				Name:  "fromEmail",
				Usage: "Email address of the sender",
			},
		},
	}
}
