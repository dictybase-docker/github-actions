package repository

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/dictyBase-docker/github-actions/internal/client"
	"github.com/dictyBase-docker/github-actions/internal/logger"
	gh "github.com/google/go-github/v32/github"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

type migration struct {
	repositories  []string
	client        *gh.Client
	from          string
	to            string
	ctx           context.Context
	repoShare     chan *gh.Repository
	repoNameShare chan string
	logger        *logrus.Entry
	pollThreshold time.Duration
	pollInterval  time.Duration
}

func (m *migration) pollForRepo(repo *gh.Repository) error {
	until := time.Now().Add(m.pollThreshold)
	ticker := time.NewTicker(m.pollInterval)
	defer ticker.Stop()
	for t := range ticker.C {
		if t.After(until) {
			return fmt.Errorf("polling timed out for repository %s", repo.GetName())
		}
		r, _, err := m.client.Repositories.Get(
			m.ctx,
			repo.GetOwner().GetLogin(),
			repo.GetName(),
		)
		if err == nil {
			m.logger.Debugf(
				"repository %s forked at %s", r.GetName(), r.GetOwner().GetLogin(),
			)
			m.repoShare <- r
			return nil
		}
		errResp, ok := err.(*gh.ErrorResponse)
		if !ok {
			return fmt.Errorf("unexpected github error %s", err)
		}
		if errResp.Response.StatusCode != http.StatusNotFound {
			return fmt.Errorf("unexpected github error %s", err)
		}
	}
	return nil
}

func (m *migration) createFork() error {
	defer close(m.repoShare)
	rgr := new(errgroup.Group)
	for _, repo := range m.repositories {
		r, _, err := m.client.Repositories.CreateFork(
			m.ctx,
			m.from,
			repo,
			&gh.RepositoryCreateForkOptions{
				Organization: m.to,
			},
		)
		if err == nil {
			return fmt.Errorf(
				"error in creating fork for %s, did not get accepted response",
				repo,
			)
		}
		if _, ok := err.(*gh.AcceptedError); !ok {
			return fmt.Errorf("error in creating fork %s", err)
		}
		m.logger.Debugf("started forking of repository %s", repo)
		func(repo *gh.Repository) {
			rgr.Go(func() error {
				return m.pollForRepo(repo)
			})
		}(r)
	}
	return rgr.Wait()
}

func (m *migration) makeArchive() error {
	defer close(m.repoNameShare)
	for repo := range m.repoShare {
		repo.Archived = gh.Bool(true)
		_, _, err := m.client.Repositories.Edit(
			m.ctx,
			repo.GetOwner().GetLogin(),
			repo.GetName(),
			repo,
		)
		if err != nil {
			return fmt.Errorf("error in setting archive status %s", err)
		}
		m.logger.Debugf("archived repository %s", repo.GetName())
		m.repoNameShare <- repo.GetName()
	}
	return nil
}

func (m *migration) delRepo() error {
	for repo := range m.repoNameShare {
		_, err := m.client.Repositories.Delete(
			m.ctx,
			m.from,
			repo,
		)
		if err != nil {
			return fmt.Errorf("error in deleting repo %s %s", repo, err)
		}
		m.logger.Debugf("deleted repository %s", repo)
	}
	return nil
}

func MigrateRepositories(c *cli.Context) error {
	gclient, err := client.GetGithubClient(c.GlobalString("token"))
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in getting github client %s", err),
			2,
		)
	}
	nc := make(chan string)
	rc := make(chan *gh.Repository)
	log := logger.GetLogger(c)
	fgr, ctx := errgroup.WithContext(context.Background())
	m := &migration{
		repositories:  c.StringSlice("repo-to-move"),
		from:          c.GlobalString("owner"),
		to:            c.String("owner-to-migrate"),
		client:        gclient,
		ctx:           ctx,
		repoShare:     rc,
		repoNameShare: nc,
		logger:        log,
	}
	fgr.Go(m.createFork)
	fgr.Go(m.makeArchive)
	fgr.Go(m.delRepo)
	if err := fgr.Wait(); err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in migrating repository %s", err),
			2,
		)
	}
	log.Infof("migrated %d repositories", len(c.StringSlice("repo-to-move")))
	return nil
}
