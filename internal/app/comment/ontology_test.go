package comment

import (
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
