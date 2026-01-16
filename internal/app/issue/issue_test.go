package issue

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func TestGetOrderDataFromIssueBodyBothFields(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	app := cli.NewApp()
	set := flag.NewFlagSet("test", 0)
	set.String("order-id", "ORD-12345", "order id")
	set.String("email", "user@example.com", "recipient email")
	set.Parse([]string{})
	ctx := cli.NewContext(app, set, nil)

	result, err := getOrderDataFromIssueBody(ctx)
	assert.NoError(err, "should not return error when both fields are provided")
	assert.NotNil(result, "should return a non-nil OrderData object")
	assert.Equal("ORD-12345", result.orderID, "should extract order ID from context")
	assert.Equal("user@example.com", result.recipientEmail, "should extract email from context")
}

func TestGetOrderDataFromIssueBodyDifferentValues(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	app := cli.NewApp()
	set := flag.NewFlagSet("test", 0)
	set.String("order-id", "ORD-99999", "order id")
	set.String("email", "admin@example.com", "recipient email")
	set.Parse([]string{})
	ctx := cli.NewContext(app, set, nil)

	result, err := getOrderDataFromIssueBody(ctx)
	assert.NoError(err, "should not return error when both fields are provided")
	assert.NotNil(result, "should return a non-nil OrderData object")
	assert.Equal("ORD-99999", result.orderID, "should extract order ID from context")
	assert.Equal("admin@example.com", result.recipientEmail, "should extract email from context")
}

func TestGetOrderDataFromIssueBodyMissingOrderID(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	app := cli.NewApp()
	set := flag.NewFlagSet("test", 0)
	set.String("email", "user@example.com", "recipient email")
	set.Parse([]string{})
	ctx := cli.NewContext(app, set, nil)

	result, err := getOrderDataFromIssueBody(ctx)
	assert.Error(err, "should return error when order-id is missing")
	assert.Nil(result, "should return nil when order-id is missing")
}

func TestGetOrderDataFromIssueBodyMissingEmail(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	app := cli.NewApp()
	set := flag.NewFlagSet("test", 0)
	set.String("order-id", "ORD-12345", "order id")
	set.Parse([]string{})
	ctx := cli.NewContext(app, set, nil)

	result, err := getOrderDataFromIssueBody(ctx)
	assert.Error(err, "should return error when email is missing")
	assert.Nil(result, "should return nil when email is missing")
}

func TestGetOrderDataFromIssueBodyMissingBoth(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	app := cli.NewApp()
	set := flag.NewFlagSet("test", 0)
	ctx := cli.NewContext(app, set, nil)

	result, err := getOrderDataFromIssueBody(ctx)
	assert.Error(err, "should return error when both fields are missing")
	assert.Nil(result, "should return nil when both fields are missing")
}
