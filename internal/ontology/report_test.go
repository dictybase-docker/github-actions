package ontology

import (
	"testing"

	"github.com/dictyBase-docker/github-actions/internal/fake"
	"github.com/stretchr/testify/require"
)

func TestParseEmptyViolations(t *testing.T) {
	assert := require.New(t)
	f, err := fake.OntoReportWithEmptyError()
	assert.NoError(
		err,
		"should not receive any error from getting error report file",
	)
	_, err = ParseViolations(f, "ERROR")
	assert.True(IsViolationNotFound(err), "should be violation not found error")
}

func TestParseViolations(t *testing.T) {
	assert := require.New(t)
	f, err := fake.OntoErrorFile()
	assert.NoError(
		err,
		"should not receive any error from getting error report file",
	)
	viol, err := ParseViolations(f, "ERROR")
	assert.NoError(
		err,
		"should not produce any error from parsing violations",
	)
	assert.Len(viol, 3, "should have 3 error violations")
	assert.Contains(
		viol,
		"missing ontology license",
		"should have missing ontology license violation",
	)
	assert.Contains(
		viol,
		"missing ontology title",
		"should have missing ontology title violation",
	)
	_, err = ParseViolations(f, "FATAL")
	assert.True(IsViolationNotFound(err), "should be violation not found error")
}
