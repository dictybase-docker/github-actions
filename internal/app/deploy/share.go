package deploy

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/sethvargo/go-githubactions"
	"github.com/urfave/cli"
)

var okey []string = []string{
	"id",
	"url",
	"payload.cluster",
	"payload.zone",
	"payload.chart",
	"payload.namespace",
	"image_tag",
	"path",
}

func getKeys(s string) []string {
	if strings.HasPrefix(s, "payload") {
		return strings.Split(s, ":")
	}
	return []string{s}
}

func ShareDeployPayload(c *cli.Context) error {
	b, err := ioutil.ReadFile(c.String("payload-file"))
	if err != nil {
		return fmt.Errorf("error in reading content from file %s", err)
	}
	a := githubactions.New()
	for _, k := range okey {
		keys := getKeys(k)
		val, err := jsonparser.GetString(b, keys...)
		if err != nil {
			return fmt.Errorf("error in reading payload value %s", err)
		}
		a.SetOutput(keys[0], val)
	}
	return nil
}
