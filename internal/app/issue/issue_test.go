package issue

import (
	"flag"
	"testing"

	"github.com/dictyBase-docker/github-actions/internal/fake"
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
	set.Int("issue", 193, "issue number")
	set.Parse([]string{})

	ctx := cli.NewContext(app, set, nil)

	// Call getIssue
	issue, err := getIssue(client, ctx)
	assert.NoError(err, "should not return error when fetching issue")
	assert.NotNil(issue, "should return a non-nil issue")

	// Assert that the issue data matches issue.json
	assert.Equal(193, issue.GetNumber(), "should have correct issue number")
	assert.Equal("Order ID:37500885 art@vandelayindustries.com", issue.GetTitle(), "should have correct title")
	assert.Equal("open", issue.GetState(), "should have correct state")
	assert.NotEmpty(issue.GetBody(), "should have non-empty body")
	assert.Contains(issue.GetBody(), "**Order ID:** 37500885", "body should contain order ID")
	assert.Contains(issue.GetBody(), "Billing address", "body should contain billing address")
}
