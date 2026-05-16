package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseChatOpsDeploy(t *testing.T) {
	t.Parallel()

	cmd := ParseChatOpsDeploy()
	assert.Equal(t, "parse-chatops-deploy", cmd.Name)
	assert.Equal(t, []string{"pcd"}, cmd.Aliases)
	assert.NotNil(t, cmd.Action)
	assert.Len(t, cmd.Flags, 2)
}

func TestOntoReportOnPullComment(t *testing.T) {
	t.Parallel()

	cmd := OntoReportOnPullComment()
	assert.Equal(t, "report-as-comment", cmd.Name)
	assert.Equal(t, []string{"rac"}, cmd.Aliases)
	assert.NotNil(t, cmd.Action)
	assert.Len(t, cmd.Flags, 4)
}

func TestSetupDaggerChecksumCmd(t *testing.T) {
	t.Parallel()

	cmd := SetupDaggerChecksumCmd()
	assert.Equal(t, "setup-dagger-checksum", cmd.Name)
	assert.NotNil(t, cmd.Action)
	assert.Len(t, cmd.Flags, 3)
}

func TestSetupDaggerBinCmd(t *testing.T) {
	t.Parallel()

	cmd := SetupDaggerBinCmd()
	assert.Equal(t, "setup-dagger-bin", cmd.Name)
	assert.NotNil(t, cmd.Action)
	assert.Len(t, cmd.Flags, 3)
}
