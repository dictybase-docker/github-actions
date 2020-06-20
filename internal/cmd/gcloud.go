package cmd

import (
	"github.com/dictyBase-docker/github-actions/internal/app/gcloud"
	"github.com/urfave/cli"
)

func GetK8sClusterCredentialsCmd() cli.Command {
	return cli.Command{
		Name:    "get-cluster-credentials",
		Aliases: []string{"gcre"},
		Usage:   "get kubernetes cluster credentials using gcloud",
		Action:  gcloud.K8sClusterCredentials,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "project",
				Required: true,
				Usage:    "Google cloud project id",
			},
			cli.StringFlag{
				Name:     "zone",
				Required: true,
				Usage:    "Compute zone for the cluster",
			},
			cli.StringFlag{
				Name:     "cluster",
				Required: true,
				Usage:    "Name of k8s cluster",
			},
		},
	}
}
