package github

import (
	"testing"

	"github.com/dictyBase-docker/github-actions/internal/fake"
	"github.com/stretchr/testify/require"
)

func TestFilterUnique(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	cc, err := fake.GithubCommitComparison()
	assert.NoError(err, "should not receive any error for parsing push event data")
	files := CommittedFiles(cc).FilterUniqueByName().List()
	assert.Len(files, 11, "should have committed 11 unique files")
	assert.Contains(
		FileNames(files),
		"dicty_assay.obo",
		"should have dicty_assay.obo file",
	)
}

func TestFilterDeleted(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	cc, err := fake.GithubCommitComparison()
	assert.NoError(err, "should not receive any error for parsing push event data")
	files := CommittedFiles(cc).FilterDeleted(true).List()
	assert.Len(files, 14, "should have committed 14 unique files")
	assert.Contains(
		FileNames(files),
		"dicty_assay.obo",
		"should have dicty_assay.obo file",
	)
}

func TestFilterSuffix(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	cc, err := fake.GithubCommitComparison()
	assert.NoError(err, "should not receive any error for parsing push event data")
	files := CommittedFiles(cc).FilterSuffix("obo").List()
	assert.Len(files, 3, "should have committed 3 unique files")
	assert.Contains(FileNames(files), "dicty_anatomy.obo", "should have dicty_anatomy.obo file")
}

func TestCommitedFiles(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	cc, err := fake.GithubCommitComparison()
	assert.NoError(err, "should not receive any error for parsing push event data")
	assert.Equal(cc.GetStatus(), "ahead", "should match the status")
	assert.Equal(cc.GetAheadBy(), 31, "should match ahead by value")
	assert.Equal(
		cc.GetTotalCommits(),
		cc.GetAheadBy(),
		"total commits and ahead by should match",
	)
	files := CommittedFiles(cc).List()
	assert.Len(files, 14, "should have committed 14 unique files")
	assert.Contains(FileNames(files), "navbar.json", "should have navbar.json file")
}

func TestFilterChain(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	cc, err := fake.GithubCommitComparison()
	assert.NoError(err, "should not receive any error for parsing push event data")
	files := CommittedFiles(cc).FilterSuffix("txt").FilterDeleted(true).FilterUniqueByName().List()
	assert.Len(files, 4, "should have committed 4 unique files")
	assert.Contains(
		FileNames(files),
		"GWDI_Strain_Annotation.txt",
		"should have GWDI_Strain_Annotation.txt file",
	)
}
