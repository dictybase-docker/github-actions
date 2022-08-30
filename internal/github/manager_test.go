package github

import (
	"testing"

	"github.com/dictyBase-docker/github-actions/internal/fake"

	"github.com/stretchr/testify/require"
)

func TestCommittedFilesInPullSync(t *testing.T) {
	t.Parallel()
	testPull(t, "pull-request-sync.json")
}

func TestCommittedFilesInPullCreate(t *testing.T) {
	t.Parallel()
	testPull(t, "pull-request-create.json")
}

func TestCommittedFilesInpush(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	r, err := fake.PushPayload()
	assert.NoError(err, "should not receive any error from reading push payload")
	server, client := fake.GhServerClient()
	defer server.Close()
	b, err := NewGithubManager(client).CommittedFilesInPush(r)
	assert.NoError(
		err,
		"should not receive any error from getting a list of committed files",
	)
	testCommitList(t, b)
}

func testCommitList(t *testing.T, bfl *ChangedFilesBuilder) {
	t.Helper()
	assert := require.New(t)
	files := bfl.FilterUniqueByName().List()
	assert.Len(files, 11, "should have committed 11 unique files")
	assert.Contains(
		FileNames(files),
		"dicty_assay.obo",
		"should have dicty_assay.obo file",
	)
	files = bfl.FilterDeleted(true).List()
	assert.Len(
		files,
		14,
		"should have committed 14 unique files",
	)
	assert.Contains(
		FileNames(files),
		"dicty_assay.obo",
		"should have dicty_assay.obo file",
	)
	files = bfl.FilterSuffix("obo").List()
	assert.Len(files, 3, "should have committed 3 unique files")
	assert.Contains(
		FileNames(files),
		"dicty_anatomy.obo",
		"should have dicty_anatomy.obo file",
	)
	files = bfl.FilterSuffix("txt").FilterDeleted(true).FilterUniqueByName().List()
	assert.Len(files, 4, "should have committed 4 unique files")
	assert.Contains(
		FileNames(files),
		"GWDI_Strain_Annotation.txt",
		"should have GWDI_Strain_Annotation.txt file",
	)
}

func testPull(t *testing.T, name string) {
	t.Helper()
	assert := require.New(t)
	reqp, err := fake.PullReqPayload(name)
	assert.NoError(
		err,
		"should not receive any error from reading payload for push",
	)
	server, client := fake.GhServerClient()
	defer server.Close()
	b, err := NewGithubManager(client).CommittedFilesInPull(reqp)
	assert.NoError(
		err,
		"should not receive any error from getting a list of committed files",
	)
	testCommitList(t, b)
}
