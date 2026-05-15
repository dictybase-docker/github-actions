package chart

import (
	"github.com/dictyBase-docker/github-actions/internal/runner"
	"github.com/urfave/cli"
)

const exitFailure = 2

func DeployChart(clt *cli.Context) error {
	helm, err := runner.NewHelm()
	if err != nil {
		return cli.NewExitError(err.Error(), exitFailure)
	}

	if err := helm.IsConnected(); err != nil {
		return cli.NewExitError(err.Error(), exitFailure)
	}

	return installOrUpgrade(clt, helm)
}

func installOrUpgrade(clt *cli.Context, helm *runner.Helm) error {
	isok, err := helm.IsChartDeployed(clt.String("name"))
	if err != nil {
		return cli.NewExitError(err.Error(), exitFailure)
	}

	prc := &runner.ChartParams{
		Name:      clt.String("name"),
		Namespace: clt.String("namespace"),
		ImageTag:  clt.String("image-tag"),
		ChartPath: clt.String("path"),
	}
	if isok {
		if err := helm.UpgradeChart(prc); err != nil {
			return cli.NewExitError(err.Error(), exitFailure)
		}
	} else {
		if err := helm.InstallChart(prc); err != nil {
			return cli.NewExitError(err.Error(), exitFailure)
		}
	}

	return nil
}
