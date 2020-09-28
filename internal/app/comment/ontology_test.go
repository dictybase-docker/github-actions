package comment

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func failData() map[string][]*reportContent {
	data := make(map[string][]*reportContent)
	data["fail"] = []*reportContent{
		{
			Name: "dicty_pheno.obo",
			Violations: []string{
				"best of the best",
				"error in fun",
				"exceptionally fun",
			},
		},
		{
			Name: "dicty_env.obo",
			Violations: []string{
				"no env",
				"green is good",
			},
		},
	}
	return data
}

func passData() map[string][]*reportContent {
	data := make(map[string][]*reportContent)
	data["pass"] = []*reportContent{
		{Name: "dicty_assay.obo"},
		{Name: "dicty_flower.obo"},
		{Name: "foobar.obo"},
	}
	return data
}

func failAndPassData() map[string][]*reportContent {
	data := make(map[string][]*reportContent)
	data["pass"] = []*reportContent{
		{Name: "dicty_assay.obo"},
		{Name: "dicty_flower.obo"},
		{Name: "foobar.obo"},
	}
	data["fail"] = []*reportContent{
		{
			Name: "dicty_pheno.obo",
			Violations: []string{
				"best of the best",
				"error in fun",
				"exceptionally fun",
			},
		},
		{
			Name: "dicty_env.obo",
			Violations: []string{
				"no env",
				"green is good",
			},
		},
	}
	return data
}

func TestMkdownOutput(t *testing.T) {
	assert := require.New(t)
	b, err := mkdownOutput(failAndPassData())
	assert.NoError(err, "should not produce any error from template execution")
	subslice := []string{
		"dicty_env",
		"dicty_pheno",
		"dicty_flower",
		"best of the best",
		"green is good",
	}
	for _, n := range subslice {
		assert.True(strings.Contains(b.String(), n))
	}
	b, err = mkdownOutput(failData())
	assert.NoError(err, "should not produce any error from template execution")
	subslice = []string{
		"dicty_env",
		"dicty_pheno",
		"best of the best",
		"green is good",
	}
	for _, n := range subslice {
		assert.True(strings.Contains(b.String(), n))
	}
	assert.False(strings.Contains(b.String(), "dicty_assay"))
	b, err = mkdownOutput(passData())
	assert.NoError(err, "should not produce any error from template execution")
	subslice = []string{
		"dicty_assay",
		"dicty_flower",
	}
	for _, n := range subslice {
		assert.True(strings.Contains(b.String(), n))
	}
	assert.False(strings.Contains(b.String(), "dicty_pheno"))
	assert.False(strings.Contains(b.String(), "best of the best"))
}

func TestListCommittedFiles(t *testing.T) {
	assert := require.New(t)
	tmpf, err := ioutil.TempFile("", "jxt")
	assert.NoError(
		err,
		"should not throw error from creating a temp file",
	)
	defer os.Remove(tmpf.Name())
	content := []string{"/onto/dicty_assay.obo", "/pronto/dicty_flower.obo"}
	for _, line := range content {
		if _, err := fmt.Fprintf(tmpf, "%s\n", line); err != nil {
			assert.NoError(
				err,
				"should not throw error from writing to the temp file",
			)
		}
	}
	files, err := listCommittedFiles(tmpf.Name())
	assert.NoError(err, "should not throw error from getting the list")
	assert.ElementsMatch(
		files,
		[]string{"dicty_assay", "dicty_flower"},
		"should match the contents of test file",
	)
}
