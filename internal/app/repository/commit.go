package repository

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dictyBase-docker/github-actions/internal/client"
	gh "github.com/dictyBase-docker/github-actions/internal/github"
	"github.com/dictyBase-docker/github-actions/internal/logger"
	"github.com/urfave/cli"
)

func FilesCommited(c *cli.Context) error {
	gclient, err := client.GetGithubClient(c.GlobalString("token"))
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in getting github client %s", err),
			2,
		)
	}
	in, err := os.Open(c.String("payload-file"))
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in reading content from file %s", err),
			2,
		)
	}
	defer in.Close()
	files, err := gh.FilterCommittedFiles(&gh.CommittedFilesParams{
		Client:      gclient,
		Input:       in,
		Event:       c.String("event-type"),
		FileSuffix:  c.String("include-file-suffix"),
		SkipDeleted: c.BoolT("skip-deleted"),
	})
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	log := logger.GetLogger(c)
	if len(files) == 0 {
		log.Warn("no committed file found matching the criteria")
	} else {
		var out io.Writer
		if len(c.String("output")) > 0 {
			w, err := os.Create(c.String("output"))
			if err != nil {
				return cli.NewExitError(
					fmt.Errorf("error in creating file %s %s", c.String("output"), err),
					2,
				)
			}
			out = w
		} else {
			out = os.Stdout
		}
		fmt.Fprint(out, strings.Join(files, "\n"))
		log.Infof("%d files has changed in the push", len(files))
	}
	return nil
}
