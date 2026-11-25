package comment

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	fakeHTML = `
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

const (
	dictyPheno     = "dicty_pheno"
	dictyAssay     = "dicty_assay"
	dictyFlower    = "dicty_flower"
	dictyEnv       = "dicty_env"
	dictyPhenoObo  = "dicty_pheno.obo"
	dictyAssayObo  = "dicty_assay.obo"
	dictyFlowerObo = "dicty_flower.obo"
	bestOfTheBest  = "best of the best"
	greenIsGood    = "green is good"
)

func failData() map[string][]*reportContent {
	data := make(map[string][]*reportContent)
	data["fail"] = []*reportContent{
		{
			Name: dictyPhenoObo,
			Violations: []string{
				bestOfTheBest,
				"error in fun",
				"exceptionally fun",
			},
		},
		{
			Name: "dicty_env.obo",
			Violations: []string{
				"no env",
				greenIsGood,
			},
			HTML: fakeHTML,
		},
	}

	return data
}

func passData() map[string][]*reportContent {
	data := make(map[string][]*reportContent)
	data["pass"] = []*reportContent{
		{Name: dictyAssayObo},
		{Name: dictyFlowerObo, HTML: fakeHTML},
		{Name: "foobar.obo"},
	}

	return data
}

func failAndPassData() map[string][]*reportContent {
	data := make(map[string][]*reportContent)
	data["pass"] = []*reportContent{
		{Name: dictyAssayObo, HTML: fakeHTML},
		{Name: dictyFlowerObo},
		{Name: "foobar.obo"},
	}
	data["fail"] = []*reportContent{
		{
			Name: dictyPhenoObo,
			Violations: []string{
				bestOfTheBest,
				"error in fun",
				"exceptionally fun",
			},
		},
		{
			Name: "dicty_env.obo",
			Violations: []string{
				"no env",
				greenIsGood,
			},
		},
	}

	return data
}

func TestMkdownOutput(t *testing.T) {
	t.Parallel()

	htmlStr := []string{
		"Full report",
		"bootstrap",
		"missing_definition",
		"missing_obsolete_label",
		"missing_superclass",
		"abolished vacuolation",
	}
	assert := require.New(t)
	bout, err := mkdownOutput(failAndPassData())
	assert.NoError(err, "should not produce any error from template execution")

	subslice := []string{
		dictyEnv,
		dictyPheno,
		dictyFlower,
		bestOfTheBest,
		greenIsGood,
	}
	for _, n := range subslice {
		assert.Contains(bout.String(), n)
	}

	for _, s := range htmlStr {
		assert.Containsf(bout.String(), s, "should have the string %s", s)
	}

	bout, err = mkdownOutput(failData())
	assert.NoError(err, "should not produce any error from template execution")

	subslice = []string{
		dictyEnv,
		dictyPheno,
		bestOfTheBest,
		greenIsGood,
	}
	for _, n := range subslice {
		assert.Contains(bout.String(), n)
	}

	assert.NotContains(bout.String(), dictyAssay)

	for _, s := range htmlStr {
		assert.Containsf(bout.String(), s, "should have the string %s", s)
	}

	bout, err = mkdownOutput(passData())
	assert.NoError(err, "should not produce any error from template execution")

	subslice = []string{
		dictyAssay,
		dictyFlower,
	}
	for _, n := range subslice {
		assert.Contains(bout.String(), n)
	}

	assert.NotContains(bout.String(), dictyPheno)
	assert.NotContains(bout.String(), bestOfTheBest)

	for _, s := range htmlStr {
		assert.Containsf(bout.String(), s, "should have the string %s", s)
	}
}

func TestListCommittedFiles(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	tmpf, err := os.CreateTemp("", "jxt")
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
		[]string{dictyAssay, dictyFlower},
		"should match the contents of test file",
	)
}
