package repository

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dictyBase-docker/github-actions/internal/client"
	"github.com/dictyBase-docker/github-actions/internal/logger"
	"github.com/google/go-github/v62/github"
	"github.com/urfave/cli"
)

type repo struct {
	name, owner string
}

func BatchMultiRepo(clt *cli.Context) error {
	gclient, err := client.GetGithubClient(clt.GlobalString("token"))
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in getting github client %s", err), 2)
	}
	rfl, wbc, err := readFiles(clt)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	path := fmt.Sprintf(
		"%s/%s", clt.String("repository-path"),
		filepath.Base(clt.String("input-file")),
	)
	msg := github.String(
		fmt.Sprintf("adding %s file",
			filepath.Base(clt.String("input-file")),
		),
	)
	for _, rpo := range parseOwnerRepo(string(rfl)) {
		_, _, err := gclient.Repositories.CreateFile(
			context.Background(),
			rpo.owner, rpo.name, path,
			&github.RepositoryContentFileOptions{
				Message: msg,
				Content: wbc,
				Branch:  github.String(clt.String("branch")),
			},
		)
		if err != nil {
			return cli.NewExitError(
				fmt.Sprintf("error in adding file %s to repository %s %s",
					path, fmt.Sprintf("%s/%s", rpo.owner, rpo.name), err,
				),
				2,
			)
		}
		logger.GetLogger(clt).Debugf(
			"uploaded file %s to repository %s", path,
			fmt.Sprintf("%s/%s", rpo.owner, rpo.name),
		)
	}

	return nil
}

func readFiles(clt *cli.Context) ([]byte, []byte, error) {
	var byr []byte
	rpl, err := os.ReadFile(clt.String("repository-list"))
	if err != nil {
		return rpl, byr, fmt.Errorf(
			"error in opening repository list %s %s",
			clt.String("repository-list"),
			err,
		)
	}
	wnc, err := os.ReadFile(clt.String("input-file"))
	if err != nil {
		return rpl, wnc, fmt.Errorf(
			"error in opening input file %s %s",
			clt.String("input-file"),
			err,
		)
	}

	return rpl, wnc, nil
}

func parseOwnerRepo(str string) []*repo {
	rpo := make([]*repo, 0)
	for _, f := range strings.Split(str, "\n") {
		v := strings.Split(f, "/")
		rpo = append(rpo, &repo{owner: v[0], name: v[1]})
	}

	return rpo
}
