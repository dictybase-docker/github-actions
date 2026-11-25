package ontology

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViolationNotFoundError_Error(t *testing.T) {
	t.Parallel()

	err := &ViolationNotFoundError{Level: "ERROR"}
	assert.Equal(t, "violation ERROR is not found", err.Error())
}

func TestIsViolationNotFound(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "violation not found error",
			err:      &ViolationNotFoundError{Level: "WARN"},
			expected: true,
		},
		{
			name:     "generic error",
			err:      errors.New("something else"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsViolationNotFound(tt.err))
		})
	}
}
