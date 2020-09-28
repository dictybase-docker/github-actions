package cmd

import "github.com/urfave/cli"

func OntoReportOnPullComment() cli.Command {
	return cli.Command{
		Name:    "report-as-comment",
		Usage:   "generate ontology report in pull request comment",
		Aliases: []string{"rac"},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "commit-list-file,c",
				Required: true,
				Usage:    "file that contain a list of committed file in the push event",
			},
			cli.StringFlag{
				Name:     "report-dir",
				Required: true,
				Usage:    "folder containing ontology reports",
			},
			cli.IntFlag{
				Name:     "pull-request-id,id",
				Required: true,
				Usage:    "id of a pull-request where the comment should be made",
			},
		},
	}
}
