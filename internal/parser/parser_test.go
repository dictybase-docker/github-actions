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

var markdownToHTMLTests = []struct {
	name            string
	markdown        string
	wantErr         bool
	containsHTML    []string
	notContainsHTML []string
}{
	{
		name:     "basic paragraph",
		markdown: "This is a paragraph.",
		wantErr:  false,
		containsHTML: []string{
			"<p>This is a paragraph.</p>",
		},
	},
	{
		name:     "headers",
		markdown: "# Header 1\n## Header 2",
		wantErr:  false,
		containsHTML: []string{
			"<h1>Header 1</h1>",
			"<h2>Header 2</h2>",
		},
	},
	{
		name:     "bold text",
		markdown: "**Order ID:** 12345",
		wantErr:  false,
		containsHTML: []string{
			"<strong>Order ID:</strong>",
			"12345",
		},
	},
	{
		name: "GFM table",
		markdown: `| Column 1 | Column 2 |
|----------|----------|
| Value 1  | Value 2  |`,
		wantErr: false,
		containsHTML: []string{
			"<table>",
			"<thead>",
			"<tbody>",
			"<th>Column 1</th>",
			"<th>Column 2</th>",
			"<td>Value 1</td>",
			"<td>Value 2</td>",
		},
	},
	{
		name:     "links",
		markdown: "[GitHub](https://github.com)",
		wantErr:  false,
		containsHTML: []string{
			"<a href=\"https://github.com\">GitHub</a>",
		},
	},
	{
		name:     "mixed content",
		markdown: "**Order ID:** 37500885\n\nSome text with [link](http://example.com)",
		wantErr:  false,
		containsHTML: []string{
			"<strong>Order ID:</strong>",
			"37500885",
			"<a href=\"http://example.com\">link</a>",
		},
	},
	{
		name:         "empty string",
		markdown:     "",
		wantErr:      false,
		containsHTML: []string{},
	},
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
	assert.Len(tables, 5, "should extract 5 tables from body_html")

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
	assert.Equal("kevin.tun@northwestern.edu", email, "should extract the correct email from billing address column")
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
	assert.Equal("10283618", orderID, "should extract the correct order ID from paragraph")
}

func verifyHTMLContent(t *testing.T, htmlNode *html.Node, containsHTML, notContainsHTML []string) {
	t.Helper()
	assert := require.New(t)

	// Convert HTML node back to string for verification
	var buf strings.Builder
	err := html.Render(&buf, htmlNode)
	assert.NoError(err, "should render HTML node to string")

	htmlString := buf.String()

	// Verify expected HTML is present
	for _, expected := range containsHTML {
		assert.Contains(htmlString, expected, "HTML should contain %q", expected)
	}

	// Verify unexpected HTML is not present
	for _, notExpected := range notContainsHTML {
		assert.NotContains(htmlString, notExpected, "HTML should not contain %q", notExpected)
	}
}

func TestMarkdownToHTML(t *testing.T) {
	t.Parallel()

	for _, testCase := range markdownToHTMLTests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			assert := require.New(t)

			htmlNode, err := MarkdownToHTML(testCase.markdown)

			if testCase.wantErr {
				assert.Error(err, "expected error for test case: %s", testCase.name)
				return
			}

			assert.NoError(err, "unexpected error for test case: %s", testCase.name)
			assert.NotNil(htmlNode, "HTML node should not be nil")

			verifyHTMLContent(t, htmlNode, testCase.containsHTML, testCase.notContainsHTML)
		})
	}
}
