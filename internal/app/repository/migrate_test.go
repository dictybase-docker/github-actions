package repository

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/dictyBase-docker/github-actions/internal/fake"
	gh "github.com/google/go-github/v32/github"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestMigrateRepositories(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	server, client := fake.GhServerClient()
	defer server.Close()
	grp := new(errgroup.Group)
	deadline := time.Now().Add(4 * time.Second)
	ctx, cancelFn := context.WithDeadline(context.Background(), deadline)
	defer cancelFn()
	mgn := &migration{
		repositories:  []string{"abc", "cde", "efg"},
		from:          "vandeley",
		to:            "varnsen",
		client:        client,
		repoShare:     make(chan *gh.Repository),
		repoNameShare: make(chan string),
		pollInterval:  time.Millisecond * 500,
		pollThreshold: ctx,
		log: logrus.NewEntry(&logrus.Logger{
			Out:   io.Discard,
			Level: logrus.ErrorLevel,
		}),
	}
	grp.Go(mgn.createFork)
	grp.Go(mgn.makeArchive)
	grp.Go(mgn.delRepo)
	err := grp.Wait()
	assert.NoError(err, "expect to have no error from migration")
}
