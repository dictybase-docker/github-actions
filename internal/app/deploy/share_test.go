package deploy

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPayload(t *testing.T) {
	t.Parallel()

	pld := &Payload{
		Cluster:   "erickube",
		Zone:      "us-central1-a",
		Chart:     "dicty-chart",
		Path:      "./helm",
		Namespace: "default",
		ImageTag:  "abc123",
	}
	encoded, err := json.Marshal(pld)
	require.NoError(t, err)

	wrapper, err := json.Marshal(string(encoded))
	require.NoError(t, err)

	result, err := GetPayload(wrapper)
	require.NoError(t, err)
	assert.Equal(t, pld, result)
}

func TestGetPayloadInvalidJSON(t *testing.T) {
	t.Parallel()

	_, err := GetPayload([]byte("not-json"))
	assert.Error(t, err)
}
