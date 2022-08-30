package deploy

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/dictyBase-docker/github-actions/internal/logger"
	"github.com/google/go-github/v32/github"
	"github.com/sethvargo/go-githubactions"
	"github.com/urfave/cli"
)

type Payload struct {
	Cluster   string `json:"cluster"`
	Zone      string `json:"zone"`
	Chart     string `json:"chart"`
	Path      string `json:"path"`
	Namespace string `json:"namespace"`
	ImageTag  string `json:"image_tag"`
}

func GetPayload(data []byte) (*Payload, error) {
	var s string
	pld := new(Payload)
	if err := json.Unmarshal(data, &s); err != nil {
		return pld, fmt.Errorf("error in decoding json data to string %s", err)
	}
	if err := json.Unmarshal([]byte(s), pld); err != nil {
		return pld, fmt.Errorf("error in decoding string to structure %s", err)
	}

	return pld, nil
}

func ShareDeployPayload(clt *cli.Context) error {
	pld, err := os.Open(clt.String("payload-file"))
	if err != nil {
		return fmt.Errorf("error in reading content from file %s", err)
	}
	defer pld.Close()
	gdl := &github.Deployment{}
	if err := json.NewDecoder(pld).Decode(gdl); err != nil {
		return fmt.Errorf("error in decoding json %s", err)
	}
	pgdl, err := GetPayload(gdl.Payload)
	if err != nil {
		return err
	}
	act := githubactions.New()
	log := logger.GetLogger(clt)
	act.SetOutput("id", strconv.Itoa(int(gdl.GetID())))
	act.SetOutput("url", gdl.GetURL())
	act.SetOutput("cluster", pgdl.Cluster)
	act.SetOutput("zone", pgdl.Zone)
	act.SetOutput("chart", pgdl.Chart)
	act.SetOutput("namespace", pgdl.Namespace)
	act.SetOutput("image_tag", pgdl.ImageTag)
	act.SetOutput("path", pgdl.Path)
	log.Info("added all keys to the output")

	return nil
}
