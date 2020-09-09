package push

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/google/go-github/v32/github"

	"github.com/stretchr/testify/require"
)

func testDataForGitPush() (*github.CommitsComparison, error) {
	cc := &github.CommitsComparison{}
	dir, err := os.Getwd()
	if err != nil {
		return cc, fmt.Errorf("unable to get current dir %s", err)
	}
	path := filepath.Join(
		filepath.Dir(dir), "../../testdata", "event.json",
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

func TestCommitedFiles(t *testing.T) {
	assert := require.New(t)
	cc, err := testDataForGitPush()
	assert.NoError(err, "should not receive any error for parsing push event data")
	assert.Equal(cc.GetStatus(), "ahead", "should match the status")
	assert.Equal(cc.GetAheadBy(), 31, "should match ahead by value")
	assert.Equal(
		cc.GetTotalCommits(),
		cc.GetAheadBy(),
		"total commits and ahead by should match",
	)
	files := committedFiles(cc, false)
	assert.Len(files, 11, "should have committed 11 unique files")
	assert.Contains(toFileNames(files), "navbar.json", "should have navbar.json file")
}

func toFileNames(s []string) []string {
	var a []string
	for _, f := range s {
		a = append(a, path.Base(f))
	}
	return a
}
