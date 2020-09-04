package chatops

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dictyBase-docker/github-actions/internal/logger"
	"github.com/google/go-github/v32/github"
	"github.com/sethvargo/go-githubactions"
	"github.com/urfave/cli"
)

type Payload struct {
	Event github.RepositoryDispatchEvent `json:"event"`
}

type SlashCommand struct {
	Args Args `json:"args"`
}

type Args struct {
	All     string      `json:"all"`
	Named   NamedArgs   `json:"name"`
	Unnamed UnnamedArgs `json:"unnamed"`
}

type NamedArgs struct {
	Cluster string `json:"cluster"`
	Branch  string `json:"branch"`
	PR      string `json:"pr"`
	Commit  string `json:"commit"`
}

type UnnamedArgs struct {
	All string `json:"all"`
}

func GetSlashCommand(data []byte) (*SlashCommand, error) {
	var s string
	o := new(SlashCommand)
	if err := json.Unmarshal(data, &s); err != nil {
		return o, fmt.Errorf("error in decoding json data to string %s", err)
	}
	if err := json.Unmarshal([]byte(s), o); err != nil {
		return o, fmt.Errorf("error in decoding string to structure %s", err)
	}
	return o, nil
}

func ShareChatOpsPayload(c *cli.Context) error {
	r, err := os.Open(c.String("payload-file"))
	if err != nil {
		return fmt.Errorf("error in reading content from file %s", err)
	}
	defer r.Close()
	d := &Payload{}
	if err := json.NewDecoder(r).Decode(d); err != nil {
		return fmt.Errorf("error in decoding json %s", err)
	}
	p, err := GetSlashCommand(d.Event.ClientPayload)
	if err != nil {
		return err
	}
	a := githubactions.New()
	log := logger.GetLogger(c)
	a.SetOutput("cluster", p.Args.Named.Cluster)
	a.SetOutput("ref", "")
	a.SetOutput("image_tag", "")
	log.Info("added all keys to the output")
	return nil
}
