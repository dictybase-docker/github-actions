package issue

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestIssueOpts(t *testing.T) {
	t.Parallel()

	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	flagSet.String("state", "open", "")

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "state"},
	}

	ctx := cli.NewContext(app, flagSet, nil)
	opts := issueOpts(ctx)
	assert.Equal(t, "open", opts.State)
	assert.Equal(t, "comments", opts.Sort)
	assert.Equal(t, 30, opts.PerPage)
}
