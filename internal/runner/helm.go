package runner

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type Helm struct {
	Cmd string
}

type ChartParams struct {
	Name      string
	Namespace string
	ImageTag  string
	ChartPath string
}

func NewHelm() (*Helm, error) {
	path, err := exec.LookPath("helm")
	if err != nil {
		return &Helm{},
			fmt.Errorf("helm command not found %s", err)
	}
	return &Helm{Cmd: path}, nil
}

func NewHelmWithPath(path string) (*Helm, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Helm{},
			fmt.Errorf("helm command %s does not exist %s", path, err)
	}
	return &Helm{Cmd: path}, nil
}

func (h *Helm) InstallChart(args *ChartParams) error {
	cmd := exec.Command(
		h.Cmd,
		"install",
		"--name",
		args.Name,
		"--namespace",
		args.Namespace,
		"--set",
		"image.tag=",
		args.ImageTag,
		args.ChartPath,
	) //nolint:gosec
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error %s in installing helm chart %s", err, cmd.String())
	}
	return nil
}

func (h *Helm) UpgradeChart(args *ChartParams) error {
	cmd := exec.Command(
		h.Cmd,
		"upgrade",
		args.Name,
		"--namespace",
		args.Namespace,
		"--set",
		"image.tag=",
		args.ImageTag,
		args.ChartPath,
	) // nolint:gosec
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error %s in upgrading helm chart %s", err, cmd.String())
	}
	return nil
}

func (h *Helm) IsChartDeployed(chart string) (bool, error) {
	cmd := exec.Command(
		h.Cmd,
		"ls",
		fmt.Sprintf("^%s$", chart),
		"--short",
	)
	output, err := cmd.Output()
	if err != nil {
		return false,
			fmt.Errorf("error %s in running command %s", err, cmd.String())
	}
	trimmed := bytes.TrimSpace(output)
	if chart == string(trimmed) {
		return true, nil
	}
	return false, nil
}

func (h *Helm) IsConnected() error {
	_, err := h.ServerVersion()
	return err
}

func (h *Helm) ServerVersion() (string, error) {
	cmd := exec.Command(
		h.Cmd,
		"version",
		"--server",
		"--short",
	)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("error in getting helm version %s", err)
	}
	return out.String(), nil
}