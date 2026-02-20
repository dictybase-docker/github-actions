package repository

import (
	"fmt"
	"os"
	"strings"

	"github.com/dictyBase-docker/github-actions/internal/client"
	gh "github.com/dictyBase-docker/github-actions/internal/github"
	"github.com/dictyBase-docker/github-actions/internal/logger"
	"github.com/urfave/cli"
)

func FilesCommited(clt *cli.Context) error {
	gclient, err := client.GetGithubClient(clt.GlobalString("token"))
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in getting github client %s", err),
			2,
		)
	}
	inp, err := os.Open(clt.String("payload-file"))
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in reading content from file %s", err),
			2,
		)
	}
	defer inp.Close()
	files, err := gh.FilterCommittedFiles(&gh.CommittedFilesParams{
		Client:      gclient,
		Input:       inp,
		Event:       clt.String("event-type"),
		FileSuffix:  clt.String("include-file-suffix"),
		SkipDeleted: clt.BoolT("skip-deleted"),
	})
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	log := logger.GetLogger(clt)
	if len(files) == 0 {
		log.Warn("no committed file found matching the criteria")

		return nil
	}
	out := os.Stdout
	if len(clt.String("output")) > 0 {
		wout, err := os.Create(clt.String("output"))
		if err != nil {
			return cli.NewExitError(
				fmt.Errorf(
					"error in creating file %s %s",
					clt.String("output"),
					err), 2)
		}
		out = wout
	}
	fmt.Fprint(out, strings.Join(files, "\n"))
	log.Infof("%d files has changed in the push", len(files))

	return nil
}
