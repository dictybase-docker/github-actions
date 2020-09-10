package chatops

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func openTestJSON() (*os.File, error) {
	f := &os.File{}
	dir, err := os.Getwd()
	if err != nil {
		return f, fmt.Errorf("unable to get current dir %s", err)
	}
	path := filepath.Join(
		filepath.Dir(dir), "../../testdata", "chatops_event.json",
	)
	r, err := os.Open(path)
	if err != nil {
		return f, fmt.Errorf("error in reading content from file %s", err)
	}
	return r, nil
}

func TestGetWorkflowInputsFromJSON(t *testing.T) {
	assert := require.New(t)
	file, err := openTestJSON()
	assert.NoError(err, "should not receive error from opening test data")
	i, err := getWorkflowInputsFromJSON(file)
	assert.NoError(err, "should not receive error from extracting workflow inputs")
	assert.Equal(i.Cluster, "erickube", "should match cluster")
	assert.Equal(i.URL, "https://github.com/dictybase-playground/github-actions-experiments/pull/18#issuecomment-690700284", "should match html-url")
	assert.Equal(i.IssueNumber, "18", "should match issue number")
	assert.Equal(i.RepositoryName, "github-actions-experiments", "should match repository name")
	assert.Equal(i.RepositoryOwner, "dictybase-playground", "should match repository owner")
	assert.Empty(i.Commit, "should have empty commit value")
	assert.Empty(i.Branch, "should have empty branch value")
}
