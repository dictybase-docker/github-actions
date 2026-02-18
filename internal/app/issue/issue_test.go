package issue

import (
	"flag"
	"testing"

	"github.com/dictyBase-docker/github-actions/internal/fake"
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
	set.Int("issueid", 193, "issue number")
	err := set.Parse([]string{})
	assert.NoError(err, "should parse flags without error")

	ctx := cli.NewContext(app, set, nil)

	// Call getIssue
	issue, issueErr := getIssue(client, ctx)
	assert.NoError(issueErr, "should not return error when fetching issue")
	assert.NotNil(issue, "should return a non-nil issue")

	// Assert that the issue data matches issue.json
	assert.Equal(193, issue.GetNumber(), "should have correct issue number")
	assert.Equal("Order ID:37500885 art@vandelayindustries.com", issue.GetTitle(), "should have correct title")
	assert.Equal("open", issue.GetState(), "should have correct state")
	assert.NotEmpty(issue.GetBody(), "should have non-empty body")
	assert.Contains(issue.GetBody(), "**Order ID:** 37500885", "body should contain order ID")
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
