package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLegacyGithubClient(t *testing.T) {
	client, err := GetLegacyGithubClient("fake-token")
	require.NoError(t, err)
	assert.NotNil(t, client)
}

func TestGetGithubClient(t *testing.T) {
	client, err := GetGithubClient("fake-token")
	require.NoError(t, err)
	assert.NotNil(t, client)
}
