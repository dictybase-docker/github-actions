package cmd

import (
	"github.com/dictyBase-docker/github-actions/internal/app/repository"
	"github.com/urfave/cli"
)

func MigrateRepositories() cli.Command {
	return cli.Command{
		Name:    "migrate-repos",
		Usage:   "fork and migrate repositories to a different owner or organization",
		Aliases: []string{"mr"},
		Action:  repository.MigrateRepositories,
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name:     "repo-to-move,m",
				Usage:    "repositories that will be migrated",
				Required: true,
			},
			cli.StringFlag{
				Name:     "owner-to-migrate,om",
				Usage:    "owner name where the repositories will be migrated",
				Required: true,
			},
		},
	}
}
