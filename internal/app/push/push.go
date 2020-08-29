package push

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/google/go-github/github"

	"github.com/dictyBase-docker/github-actions/internal/logger"
	"github.com/urfave/cli"
)

func PushFileCommited(c *cli.Context) error {
	logger := logger.GetLogger(c)
	in, out, err := inputOutput(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	pe := &github.PushEvent{}
	if err := json.NewDecoder(in).Decode(pe); err != nil {
		return fmt.Errorf("error in decoding json %s", err)
	}
	files := committedFiles(c, pe)
	if len(files) == 0 {
		logger.Warn("no committed file found matching the criteria")
		return nil
	}
	if len(c.String("skip-file-suffix")) > 0 {
		files = prefixFilter(c, files)
		if len(files) == 0 {
			logger.Warnf(
				"no committed file found after filtering though skip-file-suffix",
				c.String("skip-file-suffix"),
			)
			return nil
		}
	}
	logger.Infof("%d files has changed in the push", len(files))
	fmt.Fprintf(out, strings.Join(files, "\n"))
	return nil
}

func inputOutput(c *cli.Context) (io.Reader, io.Writer, error) {
	var in io.Reader
	var out io.Writer
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

func prefixFilter(c *cli.Context, sl []string) []string {
	var a []string
	for _, v := range sl {
		if strings.HasSuffix(v, c.String("skip-file-suffix")) {
			continue
		}
		a = append(a, v)
	}
	return a
}

func committedFiles(c *cli.Context, event *github.PushEvent) []string {
	var files []string
	for _, commit := range event.Commits {
		files = append(files, commit.Added...)
		files = append(files, commit.Modified...)
		if !c.BoolT("skip-deleted") {
			files = append(files, commit.Removed...)
		}
	}
	return UniqueString(files)
}

// UniqueString remove duplicates from string slice
func UniqueString(sl []string) []string {
	if len(sl) == 1 {
		return sl
	}
	m := make(map[string]int)
	var a []string
	for _, v := range sl {
		if _, ok := m[v]; !ok {
			a = append(a, v)
			m[v] = 1
		}
	}
	return a
}
