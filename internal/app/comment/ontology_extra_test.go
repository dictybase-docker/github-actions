package comment

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func TestReadHTMLContent(t *testing.T) {
	t.Parallel()

	t.Run("existing file", func(t *testing.T) {
		tmpDir := t.TempDir()
		htmlFile := filepath.Join(tmpDir, "report.html")
		content := "<html><body>test</body></html>"
		require.NoError(t, os.WriteFile(htmlFile, []byte(content), 0o600))

		result, err := readHTMLContent(htmlFile)
		require.NoError(t, err)
		assert.Equal(t, content, result)
	})

	t.Run("nonexistent file", func(t *testing.T) {
		result, err := readHTMLContent("/nonexistent/file.html")
		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("directory instead of file", func(t *testing.T) {
		tmpDir := t.TempDir()
		_, err := readHTMLContent(tmpDir)
		assert.Error(t, err)
	})
}

func TestReportStatusError(t *testing.T) {
	t.Parallel()

	t.Run("with failures", func(t *testing.T) {
		data := map[string][]*reportContent{
			"fail": {{Name: dictyPhenoObo}},
		}
		err := reportStatusError(data)
		assert.Error(t, err)
	})

	t.Run("no failures", func(t *testing.T) {
		data := map[string][]*reportContent{
			"pass": {{Name: dictyAssayObo}},
		}
		err := reportStatusError(data)
		assert.NoError(t, err)
	})

	t.Run("empty", func(t *testing.T) {
		data := map[string][]*reportContent{}
		err := reportStatusError(data)
		assert.NoError(t, err)
	})
}

func TestBaseNoSuffix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "obo file",
			input:    "/onto/dicty_assay.obo",
			expected: dictyAssay,
		},
		{
			name:     "simple filename",
			input:    dictyFlowerObo,
			expected: dictyFlower,
		},
		{
			name:     "multiple dots",
			input:    "/path/to/dicty_pheno.v1.obo",
			expected: dictyPheno,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, baseNoSuffix(tt.input))
		})
	}
}

func TestOntoReport(t *testing.T) {
	t.Parallel()

	t.Run("with fail and pass", func(t *testing.T) {
		reportDir := t.TempDir()

		failJSON := `[{"level": "ERROR", "violations": [{"missing_title": [{"subject": "test"}]}]}]`
		require.NoError(t,
			os.WriteFile(
				filepath.Join(reportDir, "dicty_pheno.json"),
				[]byte(failJSON),
				0o600,
			),
		)
		require.NoError(t,
			os.WriteFile(
				filepath.Join(reportDir, "dicty_pheno.html"),
				[]byte("<html>fail</html>"),
				0o600,
			),
		)

		passJSON := `[{"level": "WARN", "violations": [{"missing_label": [{"subject": "test"}]}]}]` //nolint:gosec
		require.NoError(t,
			os.WriteFile(
				filepath.Join(reportDir, "dicty_assay.json"),
				[]byte(passJSON),
				0o600,
			),
		)

		flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
		flagSet.String(reportDirFlag, reportDir, "")

		app := cli.NewApp()
		app.Flags = []cli.Flag{cli.StringFlag{Name: reportDirFlag}}

		ctx := cli.NewContext(app, flagSet, nil)
		result, err := ontoReport(ctx, []string{dictyPheno, dictyAssay})
		require.NoError(t, err)
		assert.Contains(t, result, "fail")
		assert.Contains(t, result, "pass")
		assert.Len(t, result["fail"], 1)
		assert.Len(t, result["pass"], 1)
		assert.Equal(t, dictyPhenoObo, result["fail"][0].Name)
		assert.Equal(t, dictyAssayObo, result["pass"][0].Name)
	})

	t.Run("with only failures", func(t *testing.T) {
		reportDir := t.TempDir()

		failJSON := `[{"level": "ERROR", "violations": [{"missing_title": [{"subject": "test"}]}]}]`
		require.NoError(t,
			os.WriteFile(
				filepath.Join(reportDir, "dicty_env.json"),
				[]byte(failJSON),
				0o600,
			),
		)

		flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
		flagSet.String(reportDirFlag, reportDir, "")

		app := cli.NewApp()
		app.Flags = []cli.Flag{cli.StringFlag{Name: reportDirFlag}}

		ctx := cli.NewContext(app, flagSet, nil)
		result, err := ontoReport(ctx, []string{dictyEnv})
		require.NoError(t, err)
		assert.Contains(t, result, "fail")
		assert.NotContains(t, result, "pass")
		assert.Len(t, result["fail"], 1)
	})

	t.Run("with bad json returns error", func(t *testing.T) {
		reportDir := t.TempDir()

		require.NoError(t,
			os.WriteFile(
				filepath.Join(reportDir, "dicty_bad.json"),
				[]byte("not-json"),
				0o600,
			),
		)

		flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
		flagSet.String(reportDirFlag, reportDir, "")

		app := cli.NewApp()
		app.Flags = []cli.Flag{cli.StringFlag{Name: reportDirFlag}}

		ctx := cli.NewContext(app, flagSet, nil)
		_, err := ontoReport(ctx, []string{"dicty_bad"})
		assert.Error(t, err)
	})
}
