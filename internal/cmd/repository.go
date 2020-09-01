package cmd

import (
	"github.com/dictyBase-docker/github-actions/internal/app/repository"
	"github.com/urfave/cli"
)

func BatchMultiRepo() cli.Command {
	return cli.Command{
		Name:    "batch-multi-repo",
		Usage:   "Commit a file to multiple repositories",
		Aliases: []string{"bmr"},
		Action:  repository.BatchMultiRepo,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "input-file,i",
				Usage:    "file that will be committed to repository",
				Required: true,
			},
			cli.StringFlag{
				Name:     "repository-path,rp",
				Usage:    "relative path(from root) in the repository for the input file",
				Required: true,
			},
			cli.StringFlag{
				Name:  "branch,b",
				Usage: "repository branch(should exist before committing)",
				Value: "develop",
			},
			cli.StringFlag{
				Name:     "repository-list,l",
				Usage:    "file with list of repositories name, one repository in every line",
				Required: true,
			},
		},
	}
}
