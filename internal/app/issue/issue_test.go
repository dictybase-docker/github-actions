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

const (
	shippedLabel       = "shipped"
	orderDataErrorText = "error extracting order data"
)

func TestIssueOpts(t *testing.T) {
	t.Parallel()

	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	flagSet.String("state", "open", "")

	app := cli.NewApp()
	app.Flags = []cli.Flag{cli.StringFlag{Name: "state"}}
	ctx := cli.NewContext(app, flagSet, nil)

	opts := issueOpts(ctx)
	require.Equal(t, "open", opts.State)
	require.Equal(t, "comments", opts.Sort)
	require.Equal(t, 30, opts.PerPage)
}

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
	assert.Equal(
		"Order ID:10283618 kevin.tun@northwestern.edu",
		issue.GetTitle(),
		"should have correct title",
	)
	assert.Equal("open", issue.GetState(), "should have correct state")
	assert.NotEmpty(issue.GetBody(), "should have non-empty body")
	assert.Contains(issue.GetBody(), "**Order ID:** 10283618", "body should contain order ID")
	assert.Contains(issue.GetBody(), "Billing address", "body should contain billing address")

	doc, err := parser.MarkdownToHTML(issue.GetBody())
	assert.NoError(err)
	order, err := extractAndValidateOrderData(doc, shippedLabel)
	assert.NoError(err)
	assert.Equal("10283618", order.OrderID)
	assert.Equal("kevin.tun@northwestern.edu", order.RecipientEmail)
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
				Body: new("This is the issue body with **Order ID:** 12345"),
			},
			want:    "This is the issue body with **Order ID:** 12345",
			wantErr: false,
		},
		{
			name: "empty body",
			issue: &github.Issue{
				Body: new(""),
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
				Body: new("   "),
			},
			want:    "   ",
			wantErr: false,
		},
		{
			name: "multiline body",
			issue: &github.Issue{
				Body: new("Line 1\nLine 2\n**Order ID:** 99999"),
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
		name: "embedded shipping HTML",
		markdown: `**Order ID:** 12345

## Ship To:

<div>Jane Doe<br/>jane@example.com</div>`,
		label:       shippedLabel,
		wantOrderID: "12345",
		wantEmail:   "jane@example.com",
		wantLabel:   shippedLabel,
	},
	{
		name: "billing email takes precedence over shipping column",
		markdown: `**Order ID:** 12345

| Shipping address | Billing address |
|------------------|-----------------|
| shipping@example.com | billing@example.com |`,
		label:       shippedLabel,
		wantOrderID: "12345",
		wantEmail:   "billing@example.com",
		wantLabel:   shippedLabel,
	},
	{
		name: "valid order data",
		markdown: `**Order ID:** 12345

| Shipping address | | Billing address |
|-----------------|---|-----------------|
| John Doe | | jane@example.com |
| 123 Main St | | 456 Elm St |`,
		label:       shippedLabel,
		wantOrderID: "12345",
		wantEmail:   "jane@example.com",
		wantLabel:   shippedLabel,
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
		errContains: orderDataErrorText,
	},
	{
		name: "missing email",
		markdown: `**Order ID:** 99999

| Shipping address | | Billing address |
|-----------------|---|-----------------|
| John Doe | | No email here |`,
		label:       shippedLabel,
		wantErr:     true,
		errContains: orderDataErrorText,
	},
	{
		name: "empty order ID after extraction",
		markdown: `**Order ID:**

| Shipping address | | Billing address |
|-----------------|---|-----------------|
| John Doe | | test@example.com |`,
		label:       shippedLabel,
		wantErr:     true,
		errContains: orderDataErrorText,
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
