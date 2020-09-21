package cmd

import (
	"github.com/dictyBase-docker/github-actions/internal/app/repository"
	"github.com/urfave/cli"
)

func FilesCommited() cli.Command {
	return cli.Command{
		Name:    "files-committed",
		Usage:   "outputs list of file committed in a git push or pull-request",
		Aliases: []string{"pfc"},
		Action:  repository.FilesCommited,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "payload-file,f",
				Required: true,
				Usage:    "Full path to the file that contain the event payload",
			},
			cli.StringFlag{
				Name:  "include-file-suffix,ifs",
				Usage: "file with the given suffix will only be reported",
			},
			cli.BoolTFlag{
				Name:  "skip-deleted,sd",
				Usage: "skip deleted files in the commit(s)",
			},
			cli.StringFlag{
				Name:  "output,o",
				Usage: "Name of output file, defaults to stdout",
			},
			cli.StringFlag{
				Name:  "event-type,e",
				Usage: "Name of the event, either or push or pull-request",
				Value: "push",
			},
		},
	}
}

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
