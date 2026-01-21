package issue

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

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

	// Call extractFromTitle on the title
	orderID, emailAddress, err := extractFromTitle(issueData.Title)
	assert.NoError(err, "should extract data from title")
	assert.Equal("37500885", orderID, "should extract correct order ID from title")
	assert.Equal("art@vandelayindustries.com", emailAddress, "should extract correct email address from title")
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
	assert.Equal("37500885", orderID, "should extract correct order ID from body")
	assert.Equal("art@vandelayindustries.com", emailAddress, "should extract correct billing email from body")
}
