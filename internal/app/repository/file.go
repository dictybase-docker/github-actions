package repository

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/dictyBase-docker/github-actions/internal/client"
	"github.com/dictyBase-docker/github-actions/internal/logger"
	"github.com/google/go-github/v32/github"
	"github.com/urfave/cli"
)

type repo struct {
	name, owner string
}

func AddBatchFile(c *cli.Context) error {
	log := logger.GetLogger(c)
	gclient, err := client.GetGithubClient(c.GlobalString("token"))
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in getting github client %s", err),
			2,
		)
	}
	rl, err := ioutil.ReadFile(c.String("repository-list"))
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf(
				"error in opening repository list %s %s",
				c.String("repository-list"),
				err,
			),
			2,
		)
	}
	wc, err := ioutil.ReadFile(c.String("workflow-file"))
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf(
				"error in opening workflow file %s %s",
				c.String("workflow-file"),
				err,
			),
			2,
		)
	}
	path := fmt.Sprintf(
		"%s/%s",
		c.String("repository-path"),
		filepath.Base(c.String("input-file")),
	)
	msg := github.String(
		fmt.Sprintf(
			"adding %s workflow file",
			filepath.Base(c.String("workflow-file")),
		),
	)
	for _, r := range parseOwerRepo(string(rl)) {
		_, _, err := gclient.Repositories.CreateFile(
			context.Background(),
			r.owner,
			r.name,
			path,
			&github.RepositoryContentFileOptions{
				Message: msg,
				Content: wc,
				Branch:  github.String(c.String("branch")),
			},
		)
		if err != nil {
			return cli.NewExitError(
				fmt.Sprintf(
					"error in adding file %s to repository %s %s",
					path,
					fmt.Sprintf("%s/%s", r.owner, r.name),
					err,
				),
				2,
			)
		}
		log.Debugf(
			"uploaded file %s to repository %s",
			path,
			fmt.Sprintf("%s/%s", r.owner, r.name),
		)
	}
	return nil
}

func parseOwerRepo(str string) []*repo {
	var r []*repo
	for _, f := range strings.Split(str, "\n") {
		v := strings.Split(f, "/")
		r = append(r, &repo{owner: v[0], name: v[1]})
	}
	return r
}
