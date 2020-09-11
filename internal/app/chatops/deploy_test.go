package chatops

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-github/v32/github"
	"github.com/stretchr/testify/require"
)

type mockPullRequestClient struct {
	resp *github.PullRequest
}

func (m *mockPullRequestClient) Get(ctx context.Context, owner string, repo string, number int) (*github.PullRequest, *github.Response, error) {
	return m.resp, nil, nil
}

type mockBranchClient struct {
	resp *github.Branch
}

func (m *mockBranchClient) GetBranch(ctx context.Context, owner string, repo string, branch string) (*github.Branch, *github.Response, error) {
	return m.resp, nil, nil
}

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

func TestParsePR(t *testing.T) {
	assert := require.New(t)
	mockSHA := "17f9184c165252d85994174b82fa86e7edf44b4f"
	mc := &mockPullRequestClient{
		resp: &github.PullRequest{
			Head: &github.PullRequestBranch{
				SHA: &mockSHA,
			},
		},
	}
	prc := &pullRequestClient{
		ctx:               context.Background(),
		pullRequestClient: mc,
	}
	// test output when given a commit
	i := &Inputs{
		Commit:      "f85f132b3a986c12eb0c2a61d60a5c3dd8347bf3",
		IssueNumber: "9",
	}
	o, err := parsePR(prc, i)
	assert.NoError(err, "should not have error from parsing pr")
	assert.Equal(o.ImageTag, "pr-9-f85f132", "should match pr image tag")
	assert.Equal(o.Ref, i.Commit, "should match ref value")

	// test output when not given a commit
	i2 := &Inputs{
		IssueNumber: "9",
	}
	o2, err := parsePR(prc, i2)
	assert.NoError(err, "should not have error from parsing pr")
	assert.Equal(o2.ImageTag, "pr-9-17f9184", "should match pr image tag")
	assert.Equal(o2.Ref, mockSHA, "should match ref value")
}

func TestParseIssue(t *testing.T) {
	assert := require.New(t)
	mockSHA := "17f9184c165252d85994174b82fa86e7edf44b4f"
	mc := &mockBranchClient{
		resp: &github.Branch{
			Commit: &github.RepositoryCommit{
				SHA: &mockSHA,
			},
		},
	}
	bc := &branchClient{
		ctx:          context.Background(),
		branchClient: mc,
	}
	// test when given a commit
	i := &Inputs{
		Commit:      "f85f132b3a986c12eb0c2a61d60a5c3dd8347bf3",
		IssueNumber: "9",
	}
	o, err := parseIssue(bc, i)
	assert.NoError(err, "should not have error from parsing issue")
	assert.Equal(o.ImageTag, "f85f132", "should match commit image tag")
	assert.Equal(o.Ref, i.Commit, "should match ref value")
	// test when given a branch
	i2 := &Inputs{
		Branch:      "feature/new-command",
		IssueNumber: "9",
	}
	o2, err := parseIssue(bc, i2)
	assert.NoError(err, "should not have error from parsing issue")
	assert.Equal(o2.ImageTag, "feature-new-command-17f9184", "should match branch image tag")
	assert.Equal(o2.Ref, mockSHA, "should match ref value")
}
