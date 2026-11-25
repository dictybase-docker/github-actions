package repository

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func TestReadFiles(t *testing.T) {
	t.Parallel()

	t.Run("successful read", func(t *testing.T) {
		tmpDir := t.TempDir()
		repoList := filepath.Join(tmpDir, "repos.txt")
		inputFile := filepath.Join(tmpDir, "input.txt")

		require.NoError(t, os.WriteFile(repoList, []byte("owner/repo1"), 0o600))
		require.NoError(t, os.WriteFile(inputFile, []byte("content"), 0o600))

		flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
		flagSet.String("repository-list", repoList, "")
		flagSet.String("input-file", inputFile, "")

		app := cli.NewApp()
		app.Flags = []cli.Flag{
			cli.StringFlag{Name: "repository-list"},
			cli.StringFlag{Name: "input-file"},
		}

		ctx := cli.NewContext(app, flagSet, nil)
		rpl, wnc, err := readFiles(ctx)
		require.NoError(t, err)
		assert.Equal(t, []byte("owner/repo1"), rpl)
		assert.Equal(t, []byte("content"), wnc)
	})

	t.Run("nonexistent repo list", func(t *testing.T) {
		flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
		flagSet.String("repository-list", "/nonexistent/repos.txt", "")
		flagSet.String("input-file", "/dev/null", "")

		app := cli.NewApp()
		app.Flags = []cli.Flag{
			cli.StringFlag{Name: "repository-list"},
			cli.StringFlag{Name: "input-file"},
		}

		ctx := cli.NewContext(app, flagSet, nil)
		_, _, err := readFiles(ctx)
		assert.Error(t, err)
	})
}
