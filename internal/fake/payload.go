package fake

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func PullReqPayload(name string) (io.Reader, error) {
	var r io.Reader
	dir, err := os.Getwd()
	if err != nil {
		return r, fmt.Errorf("unable to get current dir %s", err)
	}
	path := filepath.Join(
		filepath.Dir(dir),
		"../testdata",
		name,
	)
	return os.Open(path)
}

func PushPayload() (io.Reader, error) {
	var r io.Reader
	dir, err := os.Getwd()
	if err != nil {
		return r, fmt.Errorf("unable to get current dir %s", err)
	}
	path := filepath.Join(filepath.Dir(dir), "../testdata", "push.json")
	return os.Open(path)
}
