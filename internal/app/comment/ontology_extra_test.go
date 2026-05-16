package comment

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
}

func TestReportStatusError(t *testing.T) {
	t.Parallel()

	t.Run("with failures", func(t *testing.T) {
		data := map[string][]*reportContent{
			"fail": {{Name: "dicty_pheno.obo"}},
		}
		err := reportStatusError(data)
		assert.Error(t, err)
	})

	t.Run("no failures", func(t *testing.T) {
		data := map[string][]*reportContent{
			"pass": {{Name: "dicty_assay.obo"}},
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
			expected: "dicty_assay",
		},
		{
			name:     "simple filename",
			input:    "dicty_flower.obo",
			expected: "dicty_flower",
		},
		{
			name:     "multiple dots",
			input:    "/path/to/dicty_pheno.v1.obo",
			expected: "dicty_pheno",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, baseNoSuffix(tt.input))
		})
	}
}
