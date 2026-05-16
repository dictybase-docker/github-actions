package fake

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGithubCommitComparison(t *testing.T) {
	t.Parallel()

	cc, err := GithubCommitComparison()
	require.NoError(t, err)
	assert.NotNil(t, cc)
	assert.NotEmpty(t, cc.Files)
}

func TestOntoReportWithEmptyError(t *testing.T) {
	t.Parallel()

	path, err := OntoReportWithEmptyError()
	require.NoError(t, err)
	assert.Contains(t, path, "pheno_report.json")
}

func TestOntoErrorFile(t *testing.T) {
	t.Parallel()

	path, err := OntoErrorFile()
	require.NoError(t, err)
	assert.Contains(t, path, "report.json")
}

func TestPullReqPayload(t *testing.T) {
	t.Parallel()

	r, err := PullReqPayload("pull-request-sync.json")
	require.NoError(t, err)

	data, err := io.ReadAll(r)
	require.NoError(t, err)
	assert.NotEmpty(t, data)
}

func TestPushPayload(t *testing.T) {
	t.Parallel()

	r, err := PushPayload()
	require.NoError(t, err)

	data, err := io.ReadAll(r)
	require.NoError(t, err)
	assert.NotEmpty(t, data)
}

func TestGhServerClient(t *testing.T) {
	t.Parallel()

	server, _ := GhServerClient()
	defer server.Close()

	assert.NotNil(t, server)
}
