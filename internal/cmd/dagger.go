package cmd

import (
	"github.com/dictyBase-docker/github-actions/internal/app/dagger"
	"github.com/urfave/cli"
)

func SetupDaggerChecksumCmd() cli.Command {
	return cli.Command{
		Name:    "setup-dagger-checksum",
		Aliases: []string{"sc"},
		Usage:   "setup checksum of dagger binary",
		Action:  dagger.SetupDaggerCheckSum,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "checksum-file",
				Usage: "The name of the checksum file for the dagger release",
				Value: "checksums.txt",
			},
			cli.StringFlag{
				Name:  "dagger-file",
				Usage: "The suffix of the dagger tarball file which contains the binary",
				Value: "linux_amd64.tar.gz",
			},
		},
	}
}

func SetupDaggerBinCmd() cli.Command {
	return cli.Command{
		Name:    "setup-dagger-bin",
		Aliases: []string{"sd"},
		Usage:   "setup dagger command line",
		Action:  dagger.SetupDaggerBin,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "dagger-file",
				Usage: "The suffix of the dagger tarball file which contains the binary",
				Value: "linux_amd64.tar.gz",
			},
			cli.StringFlag{
				Name:     "dagger-version",
				Usage:    "version of the dagger release",
				Required: true,
			},
		},
	}
}
