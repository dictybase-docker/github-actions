package chart

import (
	"github.com/dictyBase-docker/github-actions/internal/runner"
	"github.com/urfave/cli"
)

func DeployChart(c *cli.Context) error {
	helm, err := runner.NewHelm()
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	if err := helm.IsConnected(); err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	return installOrUpgrade(c, helm)
}

func installOrUpgrade(c *cli.Context, helm *runner.Helm) error {
	ok, err := helm.IsChartDeployed(c.String("name"))
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	p := &runner.ChartParams{
		Name:      c.String("name"),
		Namespace: c.String("namespace"),
		ImageTag:  c.String("image-tag"),
		ChartPath: c.String("path"),
	}
	if ok {
		if err := helm.UpgradeChart(p); err != nil {
			return cli.NewExitError(err.Error(), 2)
		}
	} else {
		if err := helm.InstallChart(p); err != nil {
			return cli.NewExitError(err.Error(), 2)
		}
	}
	return nil
}
