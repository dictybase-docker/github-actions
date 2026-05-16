package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOwnerRepo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected []*repo
	}{
		{
			name:  "single repo",
			input: "dictybase/foobar",
			expected: []*repo{
				{owner: "dictybase", name: "foobar"},
			},
		},
		{
			name:  "multiple repos",
			input: "dictybase/foobar\ndictybase/baz",
			expected: []*repo{
				{owner: "dictybase", name: "foobar"},
				{owner: "dictybase", name: "baz"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseOwnerRepo(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
