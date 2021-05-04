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
	p := new(Payload)
	if err := json.Unmarshal(data, &s); err != nil {
		return p, fmt.Errorf("error in decoding json data to string %s", err)
	}
	if err := json.Unmarshal([]byte(s), p); err != nil {
		return p, fmt.Errorf("error in decoding string to structure %s", err)
	}
	return p, nil
}

func ShareDeployPayload(c *cli.Context) error {
	r, err := os.Open(c.String("payload-file"))
	if err != nil {
		return fmt.Errorf("error in reading content from file %s", err)
	}
	defer r.Close()
	d := &github.Deployment{}
	if err := json.NewDecoder(r).Decode(d); err != nil {
		return fmt.Errorf("error in decoding json %s", err)
	}
	p, err := GetPayload(d.Payload)
	if err != nil {
		return err
	}
	a := githubactions.New()
	log := logger.GetLogger(c)
	a.SetOutput("id", strconv.Itoa(int(d.GetID())))
	a.SetOutput("url", d.GetURL())
	a.SetOutput("cluster", p.Cluster)
	a.SetOutput("zone", p.Zone)
	a.SetOutput("chart", p.Chart)
	a.SetOutput("namespace", p.Namespace)
	a.SetOutput("image_tag", p.ImageTag)
	a.SetOutput("path", p.Path)
	log.Info("added all keys to the output")
	return nil
}
