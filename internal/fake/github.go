package fake

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	gh "github.com/google/go-github/v32/github"
)

func GithubCommitComparison() (*gh.CommitsComparison, error) {
	cc := &gh.CommitsComparison{}
	dir, err := os.Getwd()
	if err != nil {
		return cc, fmt.Errorf("unable to get current dir %s", err)
	}
	path := filepath.Join(
		filepath.Dir(dir), "../testdata", "commit-diff.json",
	)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return cc, errors.New("unable to read test file")
	}
	if err := json.Unmarshal(b, cc); err != nil {
		return cc, fmt.Errorf("error in decoding json %s", err)
	}
	return cc, nil
}
