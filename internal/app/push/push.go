package push

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
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
	defer in.Close()
	defer out.Close()
	pe := &github.PushEvent{}
	if err := json.NewDecoder(in).Decode(pe); err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in decoding json %s", err),
			2,
		)
	}
	files, msg, ok := screenFiles(c, pe)
	if !ok {
		logger.Warn(msg)
		return cli.NewExitError(msg, 2)
	}
	logger.Infof("%d files has changed in the push", len(files))
	fmt.Fprint(out, strings.Join(files, "\n"))
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

func suffixFilter(c *cli.Context, sl []string) []string {
	var a []string
	for _, v := range sl {
		if strings.HasSuffix(v, c.String("include-file-suffix")) {
			a = append(a, v)
			continue
		}
	}
	return a
}

func screenFiles(c *cli.Context, event *github.PushEvent) ([]string, string, bool) {
	files := committedFiles(event, c.BoolT("skip-deleted"))
	if len(files) == 0 {
		return files,
			"no committed file found matching the criteria",
			false
	}
	if len(c.String("include-file-suffix")) > 0 {
		files = suffixFilter(c, files)
		if len(files) == 0 {
			return files,
				fmt.Sprintf(
					"no committed file found after filtering though include-file-suffix %s",
					c.String("include-file-suffix"),
				),
				false
		}
	}
	return files, "", true
}

func committedFiles(event *github.CommitsComparison, skipDeleted bool) []string {
	var files []string
	for _, f := range event.Files {
		if skipDeleted {
			if f.GetStatus() == "deleted" {
				continue
			}
		}
		files = append(files, path.Base(f.GetFilename()))
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
