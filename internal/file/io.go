package file

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func InputOutput(clt *cli.Context) (*os.File, *os.File, error) {
	var inp *os.File
	var out *os.File
	r, err := os.Open(clt.String("payload-file"))
	if err != nil {
		return inp, out, fmt.Errorf("error in reading content from file %s", err)
	}
	inp = r
	if len(clt.String("output")) > 0 {
		w, err := os.Create(clt.String("output"))
		if err != nil {
			return inp, out, fmt.Errorf("error in creating file %s %s", clt.String("output"), err)
		}
		out = w
	} else {
		out = os.Stdout
	}

	return inp, out, nil
}
