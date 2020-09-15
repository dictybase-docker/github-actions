package push

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

func PushFileCommited(c *cli.Context) error {
	in, out, err := inputOutput(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	defer in.Close()
	defer out.Close()
	files, err := changedFiles(c, in)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	log := logger.GetLogger(c)
	if len(files) == 0 {
		log.Warn("no committed file found matching the criteria")
	} else {
		fmt.Fprint(out, strings.Join(files, "\n"))
		log.Infof("%d files has changed in the push", len(files))
	}
	return nil
}

func inputOutput(c *cli.Context) (*os.File, *os.File, error) {
	var in *os.File
	var out *os.File
	r, err := os.Open(c.String("payload-file"))
	if err != nil {
		return in, out, fmt.Errorf("error in reading content from file %s", err)
	}
	in = r
	if len(c.String("output")) > 0 {
		w, err := os.Create(c.String("output"))
		if err != nil {
			return in, out, fmt.Errorf("error in creating file %s %s", c.String("output"), err)
		}
		out = w
	} else {
		out = os.Stdout
	}
	return in, out, nil
}

func changedFiles(c *cli.Context, in io.Reader) ([]string, error) {
	gclient, err := client.GetGithubClient(c.GlobalString("token"))
	if err != nil {
		return []string{}, fmt.Errorf("error in getting github client %s", err)
	}
	fb, err := gh.NewGithubManager(gclient).CommitedFilesInPush(in)
	if err != nil {
		return []string{}, err
	}
	return fb.FilterUniqueByName().
		FilterDeleted(c.BoolT("skip-deleted")).
		FilterSuffix(c.String("include-file-suffix")).
		List(), nil
}
