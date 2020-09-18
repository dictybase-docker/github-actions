package github

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/dictyBase-docker/github-actions/internal/fake"

	"github.com/stretchr/testify/require"
)

func fakePullReqSyncPayload() (io.Reader, error) {
	var r io.Reader
	dir, err := os.Getwd()
	if err != nil {
		return r, fmt.Errorf("unable to get current dir %s", err)
	}
	path := filepath.Join(
		filepath.Dir(dir),
		"../testdata",
		"pull-request-sync.json",
	)
	return os.Open(path)
}

func fakePushPayload() (io.Reader, error) {
	var r io.Reader
	dir, err := os.Getwd()
	if err != nil {
		return r, fmt.Errorf("unable to get current dir %s", err)
	}
	path := filepath.Join(filepath.Dir(dir), "../testdata", "push.json")
	return os.Open(path)
}

func TestCommitedFilesInpush(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	r, err := fakePushPayload()
	assert.NoError(err, "should not receive any error from reading push payload")
	server, client := fake.GhServerClient()
	defer server.Close()
	b, err := NewGithubManager(client).CommittedFilesInPush(r)
	assert.NoError(
		err,
		"should not receive any error from getting a list of committed files",
	)
	files := b.FilterUniqueByName().List()
	assert.Len(files, 11, "should have committed 11 unique files")
	assert.Contains(
		FileNames(files),
		"dicty_assay.obo",
		"should have dicty_assay.obo file",
	)
	files = b.FilterDeleted(true).List()
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
	files = b.FilterSuffix("obo").List()
	assert.Len(files, 3, "should have committed 3 unique files")
	assert.Contains(
		FileNames(files),
		"dicty_anatomy.obo",
		"should have dicty_anatomy.obo file",
	)
	files = b.FilterSuffix("txt").FilterDeleted(true).FilterUniqueByName().List()
	assert.Len(files, 4, "should have committed 4 unique files")
	assert.Contains(
		FileNames(files),
		"GWDI_Strain_Annotation.txt",
		"should have GWDI_Strain_Annotation.txt file",
	)
}
