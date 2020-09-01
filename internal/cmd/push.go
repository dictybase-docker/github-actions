package cmd

import (
	"github.com/dictyBase-docker/github-actions/internal/app/push"
	"github.com/urfave/cli"
)

func PushFileCommited() cli.Command {
	return cli.Command{
		Name:    "push-file-committed",
		Usage:   "outputs list of file committed in a git push",
		Aliases: []string{"pfc"},
		Action:  push.PushFileCommited,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "payload-file,f",
				Required: true,
				Usage:    "Full path to the file that contain the push payload",
			},
			cli.StringFlag{
				Name:  "include-file-suffix,ifs",
				Usage: "file with the given suffix will only be reported",
			},
			cli.BoolTFlag{
				Name:  "skip-deleted,sd",
				Usage: "skip deleted files in the commit",
			},
			cli.StringFlag{
				Name:  "output,o",
				Usage: "Name of output file, defaults to stdout",
			},
		},
	}
}
