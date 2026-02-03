package parser

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

type IssueData struct {
	BodyHTML string `json:"body_html"`
}

func TestParseTables(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	// Load the test data
	testDataPath := filepath.Join("..", "..", "testdata", "issue.json")
	data, err := os.ReadFile(testDataPath)
	assert.NoError(err, "should be able to read testdata/issue.json")

	// Parse JSON to extract body_html
	var issue IssueData
	err = json.Unmarshal(data, &issue)
	assert.NoError(err, "should be able to parse JSON")

	// Parse HTML string to *html.Node
	doc, err := html.Parse(strings.NewReader(issue.BodyHTML))
	assert.NoError(err, "should be able to parse HTML")

	// Parse tables from HTML node
	tables, err := ParseTables(doc)
	assert.NoError(err, "should be able to parse tables from HTML")

	// Verify we extracted 4 tables
	assert.Len(tables, 4, "should extract 4 tables from body_html")

	// Verify the table headers
	assert.Equal([]string{"Shipping address", "", "Billing address"}, tables[0].Headers, "first table should be shipping and billing information")
	assert.Equal([]string{"Item", "Quantity", "Unit price($)", "Total($)"}, tables[1].Headers, "second table should be stocks ordered")
	assert.Equal([]string{"ID", "Descriptor", "Name(s)", "Systematic Name", "Characteristics"}, tables[2].Headers, "third table should be strain information")
	assert.Equal([]string{"Name", "Stored as", "Location", "No. of vials", "Color"}, tables[3].Headers, "fourth table should be strain storage")
}

func TestExtractBillingEmail(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	// Load the test data
	testDataPath := filepath.Join("..", "..", "testdata", "issue.json")
	data, err := os.ReadFile(testDataPath)
	assert.NoError(err, "should be able to read testdata/issue.json")

	// Parse JSON to extract body_html
	var issue IssueData
	err = json.Unmarshal(data, &issue)
	assert.NoError(err, "should be able to parse JSON")

	// Parse HTML string to *html.Node
	doc, err := html.Parse(strings.NewReader(issue.BodyHTML))
	assert.NoError(err, "should be able to parse HTML")

	// Extract billing email
	email, err := ExtractBillingEmail(doc)
	assert.NoError(err, "should be able to extract billing email")
	assert.Equal("art@vandelayindustries.com", email, "should extract the correct email from billing address column")
}

func TestExtractOrderID(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	// Load the test data
	testDataPath := filepath.Join("..", "..", "testdata", "issue.json")
	data, err := os.ReadFile(testDataPath)
	assert.NoError(err, "should be able to read testdata/issue.json")

	// Parse JSON to extract body_html
	var issue IssueData
	err = json.Unmarshal(data, &issue)
	assert.NoError(err, "should be able to parse JSON")

	// Parse HTML string to *html.Node
	doc, err := html.Parse(strings.NewReader(issue.BodyHTML))
	assert.NoError(err, "should be able to parse HTML")

	// Extract order ID
	orderID, err := ExtractOrderID(doc)
	assert.NoError(err, "should be able to extract order ID")
	assert.Equal("37500885", orderID, "should extract the correct order ID from paragraph")
}
