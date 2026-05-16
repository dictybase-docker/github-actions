package file

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func TestInputOutput(t *testing.T) {
	t.Parallel()

	t.Run("with input file and output file", func(t *testing.T) {
		tmpDir := t.TempDir()
		inpFile := filepath.Join(tmpDir, "input.json")
		outFile := filepath.Join(tmpDir, "output.json")

		require.NoError(t, os.WriteFile(inpFile, []byte(`{"key":"value"}`), 0o600))

		flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
		flagSet.String("payload-file", inpFile, "")
		flagSet.String("output", outFile, "")

		app := cli.NewApp()
		app.Flags = []cli.Flag{
			cli.StringFlag{Name: "payload-file"},
			cli.StringFlag{Name: "output"},
		}

		ctx := cli.NewContext(app, flagSet, nil)
		inp, out, err := InputOutput(ctx)
		require.NoError(t, err)
		assert.NotNil(t, inp)
		assert.NotNil(t, out)
		assert.NotEqual(t, os.Stdout, out)

		_ = inp.Close()
		_ = out.Close()
	})

	t.Run("with input file and stdout", func(t *testing.T) {
		tmpDir := t.TempDir()
		inpFile := filepath.Join(tmpDir, "input.json")

		require.NoError(t, os.WriteFile(inpFile, []byte(`{"key":"value"}`), 0o600))

		flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
		flagSet.String("payload-file", inpFile, "")
		flagSet.String("output", "", "")

		app := cli.NewApp()
		app.Flags = []cli.Flag{
			cli.StringFlag{Name: "payload-file"},
			cli.StringFlag{Name: "output"},
		}

		ctx := cli.NewContext(app, flagSet, nil)
		inp, out, err := InputOutput(ctx)
		require.NoError(t, err)
		assert.NotNil(t, inp)
		assert.Equal(t, os.Stdout, out)

		_ = inp.Close()
	})

	t.Run("nonexistent input file", func(t *testing.T) {
		flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
		flagSet.String("payload-file", "/nonexistent/file.json", "")
		flagSet.String("output", "", "")

		app := cli.NewApp()
		app.Flags = []cli.Flag{
			cli.StringFlag{Name: "payload-file"},
			cli.StringFlag{Name: "output"},
		}

		ctx := cli.NewContext(app, flagSet, nil)
		_, _, err := InputOutput(ctx)
		assert.Error(t, err)
	})
}
