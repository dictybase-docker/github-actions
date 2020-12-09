package repository

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	"github.com/dictyBase-docker/github-actions/internal/fake"
	gh "github.com/google/go-github/v32/github"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestMigrateRepositories(t *testing.T) {
	assert := require.New(t)
	server, client := fake.GhServerClient()
	defer server.Close()
	gr := new(errgroup.Group)
	deadline := time.Now().Add(6 * time.Second)
	ctx, cancelFn := context.WithDeadline(context.Background(), deadline)
	defer cancelFn()
	m := &migration{
		repositories:  []string{"abc", "cde", "efg"},
		from:          "vandeley",
		to:            "varnsen",
		client:        client,
		repoShare:     make(chan *gh.Repository),
		repoNameShare: make(chan string),
		pollInterval:  time.Second * 1,
		pollThreshold: ctx,
		log: logrus.NewEntry(&logrus.Logger{
			Out:   ioutil.Discard,
			Level: logrus.ErrorLevel,
		}),
	}
	gr.Go(m.createFork)
	gr.Go(m.makeArchive)
	gr.Go(m.delRepo)
	err := gr.Wait()
	assert.NoError(err, "expect to have no error from migration")
}
