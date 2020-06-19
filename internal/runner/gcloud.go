package runner

import (
	"fmt"
	"os"
	"os/exec"
)

type Gcloud struct {
	Cmd string
}

func NewGcloud() (*Gcloud, error) {
	path, err := exec.LookPath("gcloud")
	if err != nil {
		return &Gcloud{}, fmt.Errorf("gcloud command not found %s", err)
	}
	return &Gcloud{Cmd: path}, nil
}

func NewGcloudWithPath(path string) (*Gcloud, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Gcloud{}, fmt.Errorf("file %s does not exist", path)
	}
	return &Gcloud{Cmd: path}, nil
}

func (g *Gcloud) GetClusterCredentials(project, zone, cluster string) error {
	cmd := exec.Command(
		g.Cmd,
		"containers",
		"clusters",
		"get-credentials",
		cluster,
		"--zone",
		zone,
		"--project",
		project,
	)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error %s in running command %s", err, cmd.String())
	}
	return nil
}
