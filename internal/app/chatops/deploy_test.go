package chatops

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func openTestJSON(filename string) (*os.File, error) {
	f := &os.File{}
	dir, err := os.Getwd()
	if err != nil {
		return f, fmt.Errorf("unable to get current dir %s", err)
	}
	path := filepath.Join(
		filepath.Dir(dir), "../../testdata", filename,
	)
	r, err := os.Open(path)
	if err != nil {
		return f, fmt.Errorf("error in reading content from file %s", err)
	}
	return r, nil
}

func TestGetWorkflowInputsFromJSON(t *testing.T) {
	assert := require.New(t)
	// check json payload from pull request
	file, err := openTestJSON("chatops_event.json")
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
	// check json payload with branch
	file2, err := openTestJSON("chatops_event_branch.json")
	assert.NoError(err, "should not receive error from opening test data")
	i2, err := getWorkflowInputsFromJSON(file2)
	assert.NoError(err, "should not receive error from extracting workflow inputs")
	assert.Empty(i2.Commit, "should have empty commit value")
	assert.Equal(i2.Branch, "feature/new-command", "should have empty branch value")
	// check json payload for commits
	file3, err := openTestJSON("chatops_event_commit.json")
	assert.NoError(err, "should not receive error from opening test data")
	i3, err := getWorkflowInputsFromJSON(file3)
	assert.NoError(err, "should not receive error from extracting workflow inputs")
	assert.Empty(i3.Branch, "should have empty branch value")
	assert.Equal(i3.Commit, "f85f132b3a986c12eb0c2a61d60a5c3dd8347bf3", "should match commit value")
}
