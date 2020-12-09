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

type poll struct {
	repo          *gh.Repository
	client        *gh.Client
	repoShare     chan *gh.Repository
	pollThreshold context.Context
	pollInterval  time.Duration
	log           *logrus.Entry
}

func (p *poll) forRepo() error {
	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()
OUTER:
	for {
		select {
		case <-ticker.C:
			r, _, err := p.client.Repositories.Get(
				context.Background(),
				p.repo.GetOwner().GetLogin(),
				p.repo.GetName(),
			)
			if err == nil {
				p.repoShare <- r
				p.log.Debugf(
					"polling finished for repo %s/%s",
					r.GetName(), r.GetOwner().GetLogin(),
				)
				break OUTER
			}
			errResp, ok := err.(*gh.ErrorResponse)
			if !ok {
				return fmt.Errorf("unexpected github error %s", err)
			}
			if errResp.Response.StatusCode != http.StatusNotFound {
				return fmt.Errorf("unexpected github error %s", err)
			}
		case <-p.pollThreshold.Done():
			return fmt.Errorf("polling timed out for repository %s", p.repo.GetName())
		}
	}
	return nil
}

type migration struct {
	repositories  []string
	client        *gh.Client
	from          string
	to            string
	repoShare     chan *gh.Repository
	repoNameShare chan string
	pollThreshold context.Context
	pollInterval  time.Duration
	log           *logrus.Entry
}

func (m *migration) createFork() error {
	defer close(m.repoShare)
	rgr := new(errgroup.Group)
	for _, repo := range m.repositories {
		r, _, err := m.client.Repositories.CreateFork(
			context.Background(),
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
			return fmt.Errorf("error in creating fork for repo %s %v", repo, err)
		}
		m.log.Debugf(
			"created fork for repo %s on organization %s\n",
			repo, r.GetOwner().GetLogin(),
		)
		p := &poll{
			repo:          r,
			log:           m.log,
			client:        m.client,
			repoShare:     m.repoShare,
			pollThreshold: m.pollThreshold,
			pollInterval:  m.pollInterval,
		}
		rgr.Go(p.forRepo)
	}
	return rgr.Wait()
}

func (m *migration) makeArchive() error {
	defer close(m.repoNameShare)
	for repo := range m.repoShare {
		repo.Archived = gh.Bool(true)
		_, _, err := m.client.Repositories.Edit(
			context.Background(),
			repo.GetOwner().GetLogin(),
			repo.GetName(),
			repo,
		)
		if err != nil {
			return fmt.Errorf("error in setting archive status %s", err)
		}
		m.repoNameShare <- repo.GetName()
		m.log.Debugf(
			"created archive for repo %s/%s",
			repo.GetName(), repo.GetOwner().GetLogin(),
		)
	}
	return nil
}

func (m *migration) delRepo() error {
	for repo := range m.repoNameShare {
		_, err := m.client.Repositories.Delete(
			context.Background(),
			m.from,
			repo,
		)
		if err != nil {
			return fmt.Errorf("error in deleting repo %s %s", repo, err)
		}
		m.log.Debugf("deleted repo %s", repo)
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
	log := logger.GetLogger(c)
	deadline := time.Now().Add(time.Duration(c.Int64("poll-for")) * time.Second)
	ctx, cancelFn := context.WithDeadline(context.Background(), deadline)
	defer cancelFn()
	m := &migration{
		repositories:  c.StringSlice("repo-to-move"),
		from:          c.GlobalString("owner"),
		to:            c.String("owner-to-migrate"),
		pollInterval:  time.Duration(c.Int64("poll-interval")) * time.Second,
		pollThreshold: ctx,
		client:        gclient,
		log:           log,
		repoShare:     make(chan *gh.Repository),
		repoNameShare: make(chan string),
	}
	fgr := new(errgroup.Group)
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
