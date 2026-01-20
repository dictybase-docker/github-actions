package issue

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractFromTemplate(t *testing.T) {
	assert := require.New(t)

	// Read the template file
	templatePath := filepath.Join("issue.tmpl")
	templateContent, err := os.ReadFile(templatePath)
	assert.NoError(err, "should read template file")

	// Read the JSON test data file
	testDataPath := filepath.Join("..", "..", "..", "testdata", "issue.json")
	testData, err := os.ReadFile(testDataPath)
	assert.NoError(err, "should read test data file")

	// Parse JSON to extract the body field
	var issueData struct {
		Body string `json:"body"`
	}
	err = json.Unmarshal(testData, &issueData)
	assert.NoError(err, "should parse JSON")

	// Call extractFromTemplate
	result, err := extractFromTemplate(issueData.Body, string(templateContent))
	assert.NoError(err, "should extract data from template")

	// Output the result
	t.Logf("Extracted data:\n")
	for key, value := range result {
		t.Logf("  %s: %q\n", key, value)
	}
}

func TestExtractFromTitle(t *testing.T) {
	assert := require.New(t)

	// Read the JSON test data file
	testDataPath := filepath.Join("..", "..", "..", "testdata", "issue.json")
	testData, err := os.ReadFile(testDataPath)
	assert.NoError(err, "should read test data file")

	// Parse JSON to extract the title field
	var issueData struct {
		Title string `json:"title"`
	}
	err = json.Unmarshal(testData, &issueData)
	assert.NoError(err, "should parse JSON")

	// Call extractWithRegex on the title
	// orderID, emailAddress, err := extractWithRegex(issueData.Title)
	orderID, emailAddress, err := extractFromTitle(issueData.Title)
	assert.NoError(err, "should extract data from title")

	// Output the result
	t.Logf("Extracted data:\n")
	t.Logf("  %s: %q\n", "order id", orderID)
	t.Logf("  %s: %q\n", "email address", emailAddress)
}

func TestExtractFromBody(t *testing.T) {
	assert := require.New(t)

	// Read the JSON test data file
	testDataPath := filepath.Join("..", "..", "..", "testdata", "issue.json")
	testData, err := os.ReadFile(testDataPath)
	assert.NoError(err, "should read test data file")

	// Parse JSON to extract the body field
	var issueData struct {
		Body string `json:"body"`
	}
	err = json.Unmarshal(testData, &issueData)
	assert.NoError(err, "should parse JSON")

	// Call extractFromBody on the body
	orderID, emailAddress, err := extractFromBody(issueData.Body)
	assert.NoError(err, "should extract data from body")

	// Output the result
	t.Logf("Extracted data from body:\n")
	t.Logf("  %s: %q\n", "order id", orderID)
	t.Logf("  %s: %q\n", "billing email", emailAddress)
}
