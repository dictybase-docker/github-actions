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

func BatchMultiRepo(c *cli.Context) error {
	gclient, err := client.GetGithubClient(c.GlobalString("token"))
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in getting github client %s", err),
			2,
		)
	}
	rl, wc, err := readFiles(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	path := fmt.Sprintf(
		"%s/%s", c.String("repository-path"),
		filepath.Base(c.String("input-file")),
	)
	msg := github.String(
		fmt.Sprintf("adding %s file",
			filepath.Base(c.String("input-file")),
		),
	)
	for _, r := range parseOwnerRepo(string(rl)) {
		_, _, err := gclient.Repositories.CreateFile(
			context.Background(),
			r.owner, r.name, path,
			&github.RepositoryContentFileOptions{
				Message: msg,
				Content: wc,
				Branch:  github.String(c.String("branch")),
			},
		)
		if err != nil {
			return cli.NewExitError(
				fmt.Sprintf("error in adding file %s to repository %s %s",
					path, fmt.Sprintf("%s/%s", r.owner, r.name), err,
				),
				2,
			)
		}
		logger.GetLogger(c).Debugf(
			"uploaded file %s to repository %s", path,
			fmt.Sprintf("%s/%s", r.owner, r.name),
		)
	}
	return nil
}

func readFiles(c *cli.Context) ([]byte, []byte, error) {
	var b []byte
	rl, err := ioutil.ReadFile(c.String("repository-list"))
	if err != nil {
		return rl, b, fmt.Errorf(
			"error in opening repository list %s %s",
			c.String("repository-list"),
			err,
		)
	}
	wc, err := ioutil.ReadFile(c.String("input-file"))
	if err != nil {
		return rl, wc, fmt.Errorf(
			"error in opening input file %s %s",
			c.String("input-file"),
			err,
		)
	}
	return rl, wc, nil
}

func parseOwnerRepo(str string) []*repo {
	var r []*repo
	for _, f := range strings.Split(str, "\n") {
		v := strings.Split(f, "/")
		r = append(r, &repo{owner: v[0], name: v[1]})
	}
	return r
}
