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
	assert.Len(tables, 4, "should extract 4 tables from body_html")

	// Verify the table headers
	assert.Equal([]string{"Item", "Quantity", "Unit price($)", "Total($)"}, tables[0].Headers, "first table should be stocks ordered")
	assert.Equal([]string{"ID", "Descriptor", "Name(s)", "Systematic Name", "Characteristics"}, tables[1].Headers, "second table should be strain information")
	assert.Equal([]string{"Name", "Stored as", "Location", "No. of vials", "Color"}, tables[2].Headers, "third table should be strain storage")
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

func getEmailExtractionTestCases() []struct {
	name     string
	input    string
	expected string
} {
	return []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple email",
			input:    "user@example.com",
			expected: "user@example.com",
		},
		{
			name:     "email with name",
			input:    "John Doe <john@example.com>",
			expected: "john@example.com",
		},
		{
			name:     "email in HTML",
			input:    "Contact: <a href=\"mailto:support@example.com\">support@example.com</a>",
			expected: "support@example.com",
		},
		{
			name:     "email with special characters",
			input:    "user.name+tag@example.co.uk",
			expected: "user.name+tag@example.co.uk",
		},
		{
			name:     "email in multiline text",
			input:    "Art Vandelay\nVandelay Industries\nart@vandelayindustries.com\nPhone: 123-456",
			expected: "art@vandelayindustries.com",
		},
		{
			name:     "invalid - consecutive dots",
			input:    "user..name@example.com",
			expected: "",
		},
		{
			name:     "dot at start - accepted by RFC 5322",
			input:    ".user@example.com",
			expected: "user@example.com", // mail.ParseAddress strips leading dot
		},
		{
			name:     "invalid - dot at end of username",
			input:    "user.@example.com",
			expected: "",
		},
		{
			name:     "hyphen at domain start - accepted by RFC 5322 but invalid DNS",
			input:    "user@-example.com",
			expected: "user@-example.com", // mail.ParseAddress allows this
		},
		{
			name:     "invalid - no domain",
			input:    "user@",
			expected: "",
		},
		{
			name:     "no email present",
			input:    "Just some text without an email",
			expected: "",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "email with subdomain",
			input:    "admin@mail.example.com",
			expected: "admin@mail.example.com",
		},
	}
}

func TestExtractEmailFromText(t *testing.T) {
	t.Parallel()

	for _, testCase := range getEmailExtractionTestCases() {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			assert := require.New(t)

			result := extractEmailFromText(testCase.input)
			assert.Equal(testCase.expected, result,
				"extractEmailFromText(%q) = %q, want %q",
				testCase.input, result, testCase.expected)
		})
	}
}

func TestExtractStockData(t *testing.T) {
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

	// Extract stock data
	stockData, err := ExtractStockData(doc)
	assert.NoError(err, "should be able to extract stock data")

	// Verify strain information
	assert.Len(stockData.StrainInfo, 2, "should extract 2 strain info entries")

	// Verify first strain entry
	assert.Equal("DBS0351362", stockData.StrainInfo[0].ID, "first strain should have correct ID")
	assert.Equal("HL16/HL106", stockData.StrainInfo[0].Descriptor, "first strain should have correct descriptor")

	// Verify second strain entry
	assert.Equal("DBS0351363", stockData.StrainInfo[1].ID, "second strain should have correct ID")
	assert.Equal("HL84/XM101", stockData.StrainInfo[1].Descriptor, "second strain should have correct descriptor")

	// Verify plasmid information
	assert.Len(stockData.PlasmidInfo, 4, "should extract 4 plasmid info entries")

	// Verify all plasmid entries have the same ID and Name
	for i, plasmid := range stockData.PlasmidInfo {
		assert.Equal("DBP0001064", plasmid.ID, "plasmid %d should have correct ID", i)
		assert.Equal("pDDB_G0279361/lacZ", plasmid.Name, "plasmid %d should have correct name", i)
	}
}

func TestExtractOrderData(t *testing.T) {
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

	// Extract all order data
	orderData, err := ExtractOrderData(doc)
	assert.NoError(err, "should be able to extract order data")

	// Verify order ID
	assert.Equal("10283618", orderData.OrderID, "should extract correct order ID")

	// Verify recipient email
	assert.Equal("kevin.tun@northwestern.edu", orderData.RecipientEmail, "should extract correct recipient email")

	// Verify strain data
	assert.Len(orderData.StockData.StrainInfo, 2, "should extract 2 strain info entries")
	assert.Equal("DBS0351362", orderData.StockData.StrainInfo[0].ID, "first strain should have correct ID")
	assert.Equal("HL16/HL106", orderData.StockData.StrainInfo[0].Descriptor, "first strain should have correct descriptor")

	// Verify plasmid data
	assert.Len(orderData.StockData.PlasmidInfo, 4, "should extract 4 plasmid info entries")
	assert.Equal("DBP0001064", orderData.StockData.PlasmidInfo[0].ID, "first plasmid should have correct ID")
	assert.Equal("pDDB_G0279361/lacZ", orderData.StockData.PlasmidInfo[0].Name, "first plasmid should have correct name")
}
