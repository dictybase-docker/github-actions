package cmd

import (
	"github.com/dictyBase-docker/github-actions/internal/app/chart"
	"github.com/urfave/cli"
)

func DeployChartCmd() cli.Command {
	return cli.Command{
		Name:    "deploy-chart",
		Usage:   "deploy helm chart",
		Aliases: []string{"dc"},
		Action:  chart.DeployChart,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "name",
				Usage:    "Name of the chart",
				Required: true,
			},
			cli.StringFlag{
				Name:     "namespace",
				Usage:    "Kubernetes namespace",
				Required: true,
			},
			cli.StringFlag{
				Name:     "image-tag",
				Usage:    "Docker image tag",
				Required: true,
			},
			cli.StringFlag{
				Name:     "path",
				Usage:    "Relative chart path from the root of the repo",
				Required: true,
			},
		},
	}
}
