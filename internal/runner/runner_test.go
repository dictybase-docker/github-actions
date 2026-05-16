package runner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGcloudWithPath(t *testing.T) {
	t.Parallel()

	t.Run("existing file", func(t *testing.T) {
		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "gcloud")
		require.NoError(t, os.WriteFile(path, []byte("fake"), 0o600))

		gcloud, err := NewGcloudWithPath(path)
		require.NoError(t, err)
		assert.Equal(t, path, gcloud.Cmd)
	})

	t.Run("nonexistent file", func(t *testing.T) {
		_, err := NewGcloudWithPath("/nonexistent/gcloud")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist")
	})
}

func TestNewHelmWithPath(t *testing.T) {
	t.Parallel()

	t.Run("existing file", func(t *testing.T) {
		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "helm")
		require.NoError(t, os.WriteFile(path, []byte("fake"), 0o600))

		helm, err := NewHelmWithPath(path)
		require.NoError(t, err)
		assert.Equal(t, path, helm.Cmd)
	})

	t.Run("nonexistent file", func(t *testing.T) {
		_, err := NewHelmWithPath("/nonexistent/helm")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist")
	})
}
