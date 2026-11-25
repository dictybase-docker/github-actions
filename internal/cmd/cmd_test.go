package cmd

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
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

func TestAnalyticsReportCmd(t *testing.T) {
	t.Parallel()

	cmd := AnalyticsReportCmd()
	assert.Equal(t, "analytics-report", cmd.Name)
	assert.Equal(t, []string{"ar"}, cmd.Aliases)
	assert.NotNil(t, cmd.Action)
	assert.Len(t, cmd.Flags, 5)
}

func TestDeployStatusCmd(t *testing.T) {
	t.Parallel()

	cmd := DeployStatusCmd()
	assert.Equal(t, "deploy-status", cmd.Name)
	assert.Equal(t, []string{"ds"}, cmd.Aliases)
	assert.NotNil(t, cmd.Action)
	assert.Len(t, cmd.Flags, 3)
}

func TestShareDeployPayloadCmd(t *testing.T) {
	t.Parallel()

	cmd := ShareDeployPayloadCmd()
	assert.Equal(t, "share-deploy-payload", cmd.Name)
	assert.Equal(t, []string{"sdp"}, cmd.Aliases)
	assert.NotNil(t, cmd.Action)
	assert.Len(t, cmd.Flags, 1)
}

func TestGenerateMkdownCmd(t *testing.T) {
	t.Parallel()

	cmd := GenerateMkdownCmd()
	assert.Equal(t, "doc", cmd.Name)
	assert.NotNil(t, cmd.Action)
	assert.Empty(t, cmd.Flags)

	app := cli.NewApp()
	app.Name = "test-app"
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	ctx := cli.NewContext(app, flagSet, nil)
	err := cmd.Action.(func(*cli.Context) error)(ctx)
	assert.NoError(t, err)
}

func TestGetK8sClusterCredentialsCmd(t *testing.T) {
	t.Parallel()

	cmd := GetK8sClusterCredentialsCmd()
	assert.Equal(t, "get-cluster-credentials", cmd.Name)
	assert.Equal(t, []string{"gcre"}, cmd.Aliases)
	assert.NotNil(t, cmd.Action)
	assert.Len(t, cmd.Flags, 3)
}

func TestCommentsCountByDateCmds(t *testing.T) {
	t.Parallel()

	cmd := CommentsCountByDateCmds()
	assert.Equal(t, "issue-comment-count", cmd.Name)
	assert.Equal(t, []string{"icc"}, cmd.Aliases)
	assert.NotNil(t, cmd.Action)
	assert.Len(t, cmd.Flags, 1)
}

func TestIssueCommentCmds(t *testing.T) {
	t.Parallel()

	cmd := IssueCommentCmds()
	assert.Equal(t, "issue-comment-report", cmd.Name)
	assert.Equal(t, []string{"icr"}, cmd.Aliases)
	assert.NotNil(t, cmd.Action)
	assert.Len(t, cmd.Flags, 2)
}

func TestMigrateRepositories(t *testing.T) {
	t.Parallel()

	cmd := MigrateRepositories()
	assert.Equal(t, "migrate-repos", cmd.Name)
	assert.Equal(t, []string{"mr"}, cmd.Aliases)
	assert.NotNil(t, cmd.Action)
	assert.Len(t, cmd.Flags, 4)
}

func TestFilesCommited(t *testing.T) {
	t.Parallel()

	cmd := FilesCommited()
	assert.Equal(t, "files-committed", cmd.Name)
	assert.Equal(t, []string{"fc"}, cmd.Aliases)
	assert.NotNil(t, cmd.Action)
	assert.Len(t, cmd.Flags, 5)
}

func TestBatchMultiRepo(t *testing.T) {
	t.Parallel()

	cmd := BatchMultiRepo()
	assert.Equal(t, "batch-multi-repo", cmd.Name)
	assert.Equal(t, []string{"bmr"}, cmd.Aliases)
	assert.NotNil(t, cmd.Action)
	assert.Len(t, cmd.Flags, 4)
}

func TestStoreReportCmd(t *testing.T) {
	t.Parallel()

	cmd := StoreReportCmd()
	assert.Equal(t, "store-report", cmd.Name)
	assert.Equal(t, []string{"ur"}, cmd.Aliases)
	assert.NotNil(t, cmd.Action)
	assert.Len(t, cmd.Flags, 7)
}

func TestDeployChartCmd(t *testing.T) {
	t.Parallel()

	cmd := DeployChartCmd()
	assert.Equal(t, "deploy-chart", cmd.Name)
	assert.Equal(t, []string{"dc"}, cmd.Aliases)
	assert.NotNil(t, cmd.Action)
	assert.Len(t, cmd.Flags, 4)
}
