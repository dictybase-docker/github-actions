package comment

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	fakeHtml = `
			<head>
			  <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css">
			</head>
			<body>
			<table class="table table-bordered table-striped">
			<thead class="bg-dark text-white header-row">
			<tr>
			  <th>Level</th>
			  <th>Rule Name</th>
			  <th>Subject</th>
			  <th>Property</th>
			  <th>Value</th>
			</tr>
			</thead>
				<tr class="table-warning">
					<td>WARN</td>
					<td>missing_definition</td>
					<td><a href="http://purl.obolibrary.org/obo/DDSTRAINCHAR_0000001">obo:DDSTRAINCHAR_0000001</a></td>
					<td><a href="http://purl.obolibrary.org/obo/IAO_0000115">IAO:0000115</a></td>
					<td></td>
				</tr>
				<tr class="table-info">
					<td>INFO</td>
					<td>missing_superclass</td>
					<td><a href="http://purl.obolibrary.org/obo/DDSTRAINCHAR_0000001">obo:DDSTRAINCHAR_0000001</a></td>
					<td><a href="http://www.w3.org/2000/01/rdf-schema#subClassOf">rdfs:subClassOf</a></td>
					<td></td>
				</tr>
				<tr class="table-warning">
					<td>WARN</td>
					<td>missing_obsolete_label</td>
					<td><a href="http://purl.obolibrary.org/obo/DDPHENO_0000388">DDPHENO:0000388</a></td>
					<td><a href="http://www.w3.org/2000/01/rdf-schema#label">rdfs:label</a></td>
					<td>abolished vacuolation</td>
				</tr>
			</table>
			</body>
`
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
			Html: fakeHtml,
		},
	}
	return data
}

func passData() map[string][]*reportContent {
	data := make(map[string][]*reportContent)
	data["pass"] = []*reportContent{
		{Name: "dicty_assay.obo"},
		{Name: "dicty_flower.obo", Html: fakeHtml},
		{Name: "foobar.obo"},
	}
	return data
}

func failAndPassData() map[string][]*reportContent {
	data := make(map[string][]*reportContent)
	data["pass"] = []*reportContent{
		{Name: "dicty_assay.obo", Html: fakeHtml},
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
	htmlStr := []string{
		"Full report",
		"bootstrap",
		"missing_definition",
		"missing_obsolete_label",
		"missing_superclass",
		"abolished vacuolation",
	}
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
	for _, s := range htmlStr {
		assert.Containsf(b.String(), s, "should have the string %s", s)
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
	for _, s := range htmlStr {
		assert.Containsf(b.String(), s, "should have the string %s", s)
	}
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
	for _, s := range htmlStr {
		assert.Containsf(b.String(), s, "should have the string %s", s)
	}
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
