package push

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/google/go-github/v32/github"

	"github.com/dictyBase-docker/github-actions/internal/client"
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
	logger.GetLogger(c).Infof("%d files has changed in the push", len(files))
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

func changedFiles(c *cli.Context, in io.Reader) ([]string, error) {
	var files []string
	pe := &github.PushEvent{}
	if err := json.NewDecoder(in).Decode(pe); err != nil {
		return files, fmt.Errorf("error in decoding json %s", err)
	}
	gclient, err := client.GetGithubClient(c.GlobalString("token"))
	if err != nil {
		return files, fmt.Errorf("error in getting github client %s", err)
	}
	comc, _, err := gclient.Repositories.CompareCommits(
		context.Background(),
		pe.GetRepo().GetOwner().GetLogin(),
		pe.GetRepo().GetName(),
		pe.GetBefore(),
		pe.GetAfter(),
	)
	if err != nil {
		return files, fmt.Errorf("error in comparing commits %s", err)
	}
	return screenFiles(c, comc)
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

func screenFiles(c *cli.Context, event *github.CommitsComparison) ([]string, error) {
	files := committedFiles(event, c.BoolT("skip-deleted"))
	if len(files) == 0 {
		return files,
			errors.New("no committed file found matching the criteria")
	}
	if len(c.String("include-file-suffix")) > 0 {
		files = suffixFilter(c, files)
		if len(files) == 0 {
			return files,
				fmt.Errorf(
					"no committed file found after filtering though include-file-suffix %s",
					c.String("include-file-suffix"),
				)
		}
	}
	return files, nil
}

func committedFiles(event *github.CommitsComparison, skipDeleted bool) []string {
	var files []string
	for _, f := range event.Files {
		if skipDeleted {
			if f.GetStatus() == "deleted" {
				continue
			}
		}
		files = append(files, f.GetFilename())
	}
	return uniqueFiles(files)
}

func uniqueFiles(sl []string) []string {
	if len(sl) == 1 {
		return sl
	}
	m := make(map[string]int)
	var a []string
	for _, v := range sl {
		n := path.Base(v)
		if _, ok := m[n]; !ok {
			a = append(a, v)
			m[n] = 1
		}
	}
	return a
}
