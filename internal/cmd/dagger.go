package cmd

import (
	"github.com/dictyBase-docker/github-actions/internal/app/dagger"
	"github.com/urfave/cli"
)

func SetupDaggerCmd() cli.Command {
	return cli.Command{
		Name:    "setup-dagger",
		Aliases: []string{"sd"},
		Usage:   "setup dagger command line",
		Action:  dagger.SetupDagger,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "dagger-file",
				Usage: "The suffix of the dagger tarball file which contains the binary",
				Value: "linux_amd64.tar.gz",
			},
			cli.StringFlag{
				Name:  "checksum-file",
				Usage: "The name of the checksum file for the dagger release",
				Value: "checksums.txt",
			},
		},
	}
}
