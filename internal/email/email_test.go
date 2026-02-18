package email

import (
	"os"
	"reflect"
	"regexp"
	"testing"

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

	// Extract all template field references ({{.FieldName}})
	fieldPatestCaseern := regexp.MustCompile(`\{\{\.(\w+)\}\}`)
	matches := fieldPatestCaseern.FindAllStringSubmatch(string(templateContent), -1)

	// Collect unique field names used in template
	templateFields := make(map[string]bool)
	for _, match := range matches {
		if len(match) > 1 {
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

	// Verify every template field exists in the struct
	for templateField := range templateFields {
		assert.True(structFields[templateField],
			"template uses field %q which doesn't exist in OrderEmailData struct", templateField)
	}

	// Note: We don't require all struct fields to be used in the template,
	// as RecipientEmail is used for addressing but not in the email body
}

func TestCreateEmailHTML(t *testing.T) {
	t.Parallel()

	tests := []struct {
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
				"Order Update # ORD-12345",
				"Your order status: shipped",
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
				"Order Update # ORD-99999",
				"Your order status: processing",
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
				"Order Update # ",
				"Your order status: ",
			},
		},
	}

	for _, testCase := range tests {
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
