package file

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

const (
	payloadFileFlag = "payload-file"
	outputFlag      = "output"
)

func InputOutput(clt *cli.Context) (*os.File, *os.File, error) {
	var (
		inp *os.File
		out *os.File
	)

	r, err := os.Open(clt.String(payloadFileFlag))
	if err != nil {
		return inp, out, fmt.Errorf("error in reading content from file %s", err)
	}

	inp = r

	if len(clt.String(outputFlag)) > 0 {
		w, err := os.Create(clt.String(outputFlag))
		if err != nil {
			return inp, out, fmt.Errorf("error in creating file %s %s", clt.String(outputFlag), err)
		}

		out = w
	} else {
		out = os.Stdout
	}

	return inp, out, nil
}
