package file

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func InputOutput(c *cli.Context) (*os.File, *os.File, error) {
	var in *os.File
	var out *os.File
	r, err := os.Open(c.String("payload-file"))
	if err != nil {
		return in, out, fmt.Errorf("error in reading content from file %s", err)
	}
	in = r
	if len(c.String("output")) > 0 {
		w, err := os.Create(c.String("output"))
		if err != nil {
			return in, out, fmt.Errorf("error in creating file %s %s", c.String("output"), err)
		}
		out = w
	} else {
		out = os.Stdout
	}
	return in, out, nil
}
