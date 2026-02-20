package issue

import (
	"flag"
	"testing"

	"github.com/dictyBase-docker/github-actions/internal/fake"
	parser "github.com/dictyBase-docker/github-actions/internal/parser"
	"github.com/google/go-github/v62/github"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func TestGetIssue(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	// Set up fake server and client
	server, client := fake.GhServerClient()
	defer server.Close()

	// Create CLI context with required flags
	app := cli.NewApp()
	set := flag.NewFlagSet("test", 0)
	set.String("owner", "dictybase-playground", "repository owner")
	set.String("repository", "learn-github-action", "repository name")
	set.Int("issueid", 122, "issue number")
	err := set.Parse([]string{})
	assert.NoError(err, "should parse flags without error")

	ctx := cli.NewContext(app, set, nil)

	// Call getIssue
	issue, issueErr := getIssue(client, ctx)
	assert.NoError(issueErr, "should not return error when fetching issue")
	assert.NotNil(issue, "should return a non-nil issue")

	// Assert that the issue data matches issue.json
	assert.Equal(122, issue.GetNumber(), "should have correct issue number")
	assert.Equal("Order ID:10283618 kevin.tun@northwestern.edu", issue.GetTitle(), "should have correct title")
	assert.Equal("open", issue.GetState(), "should have correct state")
	assert.NotEmpty(issue.GetBody(), "should have non-empty body")
	assert.Contains(issue.GetBody(), "**Order ID:** 10283618", "body should contain order ID")
	assert.Contains(issue.GetBody(), "Billing address", "body should contain billing address")
}

func TestGetIssueBody(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		issue   *github.Issue
		want    string
		wantErr bool
	}{
		{
			name: "valid body",
			issue: &github.Issue{
				Body: github.String("This is the issue body with **Order ID:** 12345"),
			},
			want:    "This is the issue body with **Order ID:** 12345",
			wantErr: false,
		},
		{
			name: "empty body",
			issue: &github.Issue{
				Body: github.String(""),
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "nil body pointer",
			issue: &github.Issue{
				Body: nil,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "body with whitespace only",
			issue: &github.Issue{
				Body: github.String("   "),
			},
			want:    "   ",
			wantErr: false,
		},
		{
			name: "multiline body",
			issue: &github.Issue{
				Body: github.String("Line 1\nLine 2\n**Order ID:** 99999"),
			},
			want:    "Line 1\nLine 2\n**Order ID:** 99999",
			wantErr: false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			assert := require.New(t)

			body, err := getIssueBody(testCase.issue)

			if testCase.wantErr {
				assert.Error(err, "expected error for test case: %s", testCase.name)
				assert.Contains(err.Error(), "issue body is empty",
					"error message should indicate empty body")
				return
			}

			assert.NoError(err, "unexpected error for test case: %s", testCase.name)
			assert.Equal(testCase.want, body, "body should match expected value")
		})
	}
}

var extractAndValidateOrderDataTests = []struct {
	name        string
	markdown    string
	label       string
	wantOrderID string
	wantEmail   string
	wantLabel   string
	wantErr     bool
	errContains string
}{
	{
		name: "valid order data",
		markdown: `**Order ID:** 12345

| Shipping address | | Billing address |
|-----------------|---|-----------------|
| John Doe | | jane@example.com |
| 123 Main St | | 456 Elm St |`,
		label:       "shipped",
		wantOrderID: "12345",
		wantEmail:   "jane@example.com",
		wantLabel:   "shipped",
		wantErr:     false,
	},
	{
		name: "missing order ID",
		markdown: `Some text without order ID

| Shipping address | | Billing address |
|-----------------|---|-----------------|
| John Doe | | jane@example.com |`,
		label:       "processing",
		wantErr:     true,
		errContains: "error extracting order data",
	},
	{
		name: "missing email",
		markdown: `**Order ID:** 99999

| Shipping address | | Billing address |
|-----------------|---|-----------------|
| John Doe | | No email here |`,
		label:       "shipped",
		wantErr:     true,
		errContains: "error extracting order data",
	},
	{
		name: "empty order ID after extraction",
		markdown: `**Order ID:**

| Shipping address | | Billing address |
|-----------------|---|-----------------|
| John Doe | | test@example.com |`,
		label:       "shipped",
		wantErr:     true,
		errContains: "error extracting order data",
	},
	{
		name: "different label",
		markdown: `**Order ID:** 54321

| Shipping address | | Billing address |
|-----------------|---|-----------------|
| John Doe | | admin@test.com |`,
		label:       "cancelled",
		wantOrderID: "54321",
		wantEmail:   "admin@test.com",
		wantLabel:   "cancelled",
		wantErr:     false,
	},
}

func TestExtractAndValidateOrderData(t *testing.T) {
	t.Parallel()

	for _, testCase := range extractAndValidateOrderDataTests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			assert := require.New(t)

			// Convert markdown to HTML
			htmlNode, err := parser.MarkdownToHTML(testCase.markdown)
			assert.NoError(err, "should convert markdown to HTML")

			// Call the function under test
			emailData, err := extractAndValidateOrderData(htmlNode, testCase.label)

			if testCase.wantErr {
				assert.Error(err, "expected error for test case: %s", testCase.name)
				if testCase.errContains != "" {
					assert.Contains(err.Error(), testCase.errContains,
						"error should contain %q", testCase.errContains)
				}
				return
			}

			assert.NoError(err, "unexpected error for test case: %s", testCase.name)
			assert.Equal(testCase.wantOrderID, emailData.OrderID, "order ID should match")
			assert.Equal(testCase.wantEmail, emailData.RecipientEmail, "email should match")
			assert.Equal(testCase.wantLabel, emailData.Label, "label should match")
		})
	}
}
