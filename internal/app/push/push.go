package push

import (
	"fmt"
	"strings"

	"github.com/dictyBase-docker/github-actions/internal/file"
	gh "github.com/dictyBase-docker/github-actions/internal/github"
	"github.com/dictyBase-docker/github-actions/internal/logger"
	"github.com/urfave/cli"
)

func PushFileCommited(c *cli.Context) error {
	in, out, err := file.InputOutput(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	defer in.Close()
	defer out.Close()
	files, err := gh.FilterCommittedFiles(c, in, "push")
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
