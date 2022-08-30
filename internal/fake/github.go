package fake

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	gh "github.com/google/go-github/v32/github"
)

func GithubCommitComparison() (*gh.CommitsComparison, error) {
	ccg := &gh.CommitsComparison{}
	dir, err := os.Getwd()
	if err != nil {
		return ccg, fmt.Errorf("unable to get current dir %s", err)
	}
	path := filepath.Join(
		filepath.Dir(dir), "../testdata", "commit-diff.json",
	)
	b, err := os.ReadFile(path)
	if err != nil {
		return ccg, errors.New("unable to read test file")
	}
	if err := json.Unmarshal(b, ccg); err != nil {
		return ccg, fmt.Errorf("error in decoding json %s", err)
	}

	return ccg, nil
}
