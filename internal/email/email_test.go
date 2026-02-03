package email

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewEmailClient(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	domain := "test.mailgun.org"
	apiKey := "test-api-key-123"
	fromEmail := "test@example.com"

	client := NewEmailClient(domain, apiKey, fromEmail)

	assert.NotNil(client, "client should not be nil")
	assert.NotNil(client.mg, "mailgun client should not be nil")
	assert.Equal(domain, client.config.Domain, "domain should match")
	assert.Equal(apiKey, client.config.APIKey, "API key should match")
	assert.Equal(fromEmail, client.config.From, "from email should match")
}
