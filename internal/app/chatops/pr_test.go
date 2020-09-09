package chatops

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToOutput(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	// test output when given a commit ID
	i := &Input{
		Commit:        "c675a18e9c395e59960ddbd0ab0a023d962749d6",
		SHA:           "a19fc461a48ff70230c4d440da1f35db075e33e2",
		PullRequestID: "3",
	}
	o := convertToOutput(i)
	assert.Exactly(o.Ref, i.Commit, "ref should match commit id")
	assert.Exactly(o.ImageTag, "pr-3-c675a18", "should match expected pr tag format")

	// test output when not given a commit ID
	i2 := &Input{
		SHA:           "a19fc461a48ff70230c4d440da1f35db075e33e2",
		PullRequestID: "9",
	}
	o2 := convertToOutput(i2)
	assert.Exactly(o2.Ref, i2.SHA, "ref should match head SHA")
	assert.Exactly(o2.ImageTag, "pr-9-a19fc46", "should match expected pr tag format")
}
