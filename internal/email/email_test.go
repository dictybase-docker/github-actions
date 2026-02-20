package email

import (
	"os"
	"reflect"
	"regexp"
	"testing"

	"github.com/dictyBase-docker/github-actions/internal/parser"
	"github.com/stretchr/testify/require"
)

func TestNewEmailClient(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	domain := "test.mailgun.org"
	apiKey := "test-api-key-123"
	fromEmail := "test@example.com"

	client := NewEmailClient(domain, fromEmail, apiKey)

	assert.NotNil(client, "client should not be nil")
	assert.NotNil(client.mg, "mailgun client should not be nil")
	assert.Equal(domain, client.config.Domain, "domain should match")
	assert.Equal(fromEmail, client.config.From, "from email should match")
}

func TestTemplateFieldsMatchStruct(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	// Read the template file
	templatePath := "order_update.tmpl"
	templateContent, err := os.ReadFile(templatePath)
	if os.IsNotExist(err) {
		t.Skip("Skipping test: order_update.tmpl not found in current directory")
	}
	assert.NoError(err, "should read template file")

	content := string(templateContent)

	// Remove content within {{range}}...{{end}} blocks to exclude iteration-scoped fields
	rangePattern := regexp.MustCompile(`(?s)\{\{range\s+[^}]+\}\}.*?\{\{end\}\}`)
	contentWithoutRanges := rangePattern.ReplaceAllString(content, "")

	// Extract all template field references from non-range content
	// Matches {{.Field}}, {{.Field.Nested}}, {{if .Field}}, etc.
	fieldPattern := regexp.MustCompile(`\{\{[^}]*?\.(\w+)(?:\.\w+)*[^}]*?\}\}`)
	matches := fieldPattern.FindAllStringSubmatch(contentWithoutRanges, -1)

	// Collect unique top-level field names used in template
	templateFields := make(map[string]bool)
	for _, match := range matches {
		if len(match) > 1 {
			// Extract only the first field name after the dot
			templateFields[match[1]] = true
		}
	}

	assert.NotEmpty(templateFields, "template should use at least one field")

	// Get actual struct fields using reflection
	structType := reflect.TypeOf(OrderEmailData{})
	structFields := make(map[string]bool)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if field.IsExported() {
			structFields[field.Name] = true
		}
	}

	// Verify every top-level template field exists in the struct
	for templateField := range templateFields {
		assert.True(structFields[templateField],
			"template uses field %q which doesn't exist in OrderEmailData struct", templateField)
	}

	// Note: We don't require all struct fields to be used in the template,
	// as RecipientEmail is used for addressing but not in the email body
}

//nolint:funlen // Test cases table
func getEmailHTMLTestCases() []struct {
	name     string
	data     OrderEmailData
	wantErr  bool
	contains []string
} {
	return []struct {
		name     string
		data     OrderEmailData
		wantErr  bool
		contains []string
	}{
		{
			name: "valid order data",
			data: OrderEmailData{
				RecipientEmail: "test@example.com",
				OrderID:        "ORD-12345",
				Label:          "shipped",
			},
			wantErr: false,
			contains: []string{
				"Order Number: <span class=\"order-id-value\">ORD-12345</span>",
				"Current Status",
				"shipped",
				"Dicty Stock Center",
				"dictystocks@northwestern.edu",
			},
		},
		{
			name: "order with different status",
			data: OrderEmailData{
				RecipientEmail: "user@test.com",
				OrderID:        "ORD-99999",
				Label:          "processing",
			},
			wantErr: false,
			contains: []string{
				"Order Number: <span class=\"order-id-value\">ORD-99999</span>",
				"processing",
			},
		},
		{
			name: "empty order data",
			data: OrderEmailData{
				RecipientEmail: "",
				OrderID:        "",
				Label:          "",
			},
			wantErr: false,
			contains: []string{
				"Order Number:",
				"Current Status",
			},
		},
		{
			name: "order with strain data",
			data: OrderEmailData{
				RecipientEmail: "test@example.com",
				OrderID:        "ORD-11111",
				Label:          "shipped",
				StockData: parser.StockData{
					StrainInfo: []parser.StrainInfo{
						{ID: "DBS0351362", Descriptor: "HL16/HL106"},
						{ID: "DBS0351363", Descriptor: "HL84/XM101"},
					},
				},
			},
			wantErr: false,
			contains: []string{
				"Order Number: <span class=\"order-id-value\">ORD-11111</span>",
				"shipped",
				"Strains Ordered",
				"DBS0351362",
				"HL16/HL106",
				"DBS0351363",
				"HL84/XM101",
			},
		},
		{
			name: "order with plasmid data",
			data: OrderEmailData{
				RecipientEmail: "test@example.com",
				OrderID:        "ORD-22222",
				Label:          "processing",
				StockData: parser.StockData{
					PlasmidInfo: []parser.PlasmidInfo{
						{ID: "DBP0001064", Name: "pDDB_G0279361/lacZ"},
						{ID: "DBP0001065", Name: "pDDB_G0279362/lacZ"},
					},
				},
			},
			wantErr: false,
			contains: []string{
				"Order Number: <span class=\"order-id-value\">ORD-22222</span>",
				"processing",
				"Plasmids Ordered",
				"DBP0001064",
				"pDDB_G0279361/lacZ",
				"DBP0001065",
				"pDDB_G0279362/lacZ",
			},
		},
		{
			name: "order with both strains and plasmids",
			data: OrderEmailData{
				RecipientEmail: "test@example.com",
				OrderID:        "ORD-33333",
				Label:          "shipped",
				StockData: parser.StockData{
					StrainInfo: []parser.StrainInfo{
						{ID: "DBS0351362", Descriptor: "HL16/HL106"},
					},
					PlasmidInfo: []parser.PlasmidInfo{
						{ID: "DBP0001064", Name: "pDDB_G0279361/lacZ"},
					},
				},
			},
			wantErr: false,
			contains: []string{
				"Order Number: <span class=\"order-id-value\">ORD-33333</span>",
				"Strains Ordered",
				"DBS0351362",
				"Plasmids Ordered",
				"DBP0001064",
			},
		},
	}
}

func TestCreateEmailHTML(t *testing.T) {
	t.Parallel()

	for _, testCase := range getEmailHTMLTestCases() {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			assert := require.New(t)

			html, err := createEmailHTML(testCase.data)

			if testCase.wantErr {
				assert.Error(err, "expected error for test case: %s", testCase.name)
				return
			}

			assert.NoError(err, "unexpected error for test case: %s", testCase.name)
			assert.NotEmpty(html, "HTML should not be empty")

			// Verify HTML contains expected content
			for _, expectedContent := range testCase.contains {
				assert.Contains(html, expectedContent,
					"HTML should contain %q in test case: %s", expectedContent, testCase.name)
			}

			// Verify it's valid HTML structure
			assert.Contains(html, "<!DOCTYPE html", "should start with HTML doctype")
			assert.Contains(html, "</html>", "should end with closing html tag")
		})
	}
}
