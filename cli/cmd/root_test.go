package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// ---------------------------------------------------------------------------
// Basic assertions with assert & require
// ---------------------------------------------------------------------------

func TestNewRootCmd_Basic(t *testing.T) {
	info := BuildInfo{Version: "1.0.0", Commit: "abc1234", Date: "2026-05-05T00:00:00Z"}
	cmd := NewRootCmd(info)

	// assert: soft check, continues on failure
	assert.Equal(t, "bangumi", cmd.Use, "command name should be bangumi")
	assert.NotEmpty(t, cmd.Short, "short description should not be empty")

	// require: hard check, stops test on failure
	require.NotNil(t, cmd.RunE, "RunE should not be nil")
}

func TestNewRootCmd_VersionOutput(t *testing.T) {
	info := BuildInfo{Version: "2.3.1", Commit: "d9cd9d9", Date: "2026-01-15T12:30:00Z"}
	cmd := NewRootCmd(info)

	assert.Equal(t, "2.3.1", cmd.Version, "version should match build info")

	// Capture version template output
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--version"})
	err := cmd.Execute()

	require.NoError(t, err)
	assert.Contains(t, buf.String(), "bangumi 2.3.1")
	assert.Contains(t, buf.String(), "commit: d9cd9d9")
	assert.Contains(t, buf.String(), "built: 2026-01-15T12:30:00Z")
}

func TestNewRootCmd_HelpOnNoArgs(t *testing.T) {
	info := BuildInfo{Version: "dev", Commit: "none", Date: "unknown"}
	cmd := NewRootCmd(info)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	// 重置 args 确保不受之前测试影响
	cmd.SetArgs(nil)
	err := cmd.Execute()

	require.NoError(t, err)
	assert.Contains(t, buf.String(), "bangumi")
}

// ---------------------------------------------------------------------------
// Table-driven tests
// ---------------------------------------------------------------------------

func TestNewRootCmd_TableDriven(t *testing.T) {
	tests := []struct {
		name    string
		info    BuildInfo
		args    []string
		wantOut string
	}{
		{
			name:    "help flag",
			info:    BuildInfo{Version: "dev", Commit: "none", Date: "unknown"},
			args:    []string{"--help"},
			wantOut: "Usage:",
		},
		{
			name:    "verbose flag",
			info:    BuildInfo{Version: "dev", Commit: "none", Date: "unknown"},
			args:    []string{"--verbose"},
			wantOut: "Usage:",
		},
		{
			name:    "debug flag",
			info:    BuildInfo{Version: "dev", Commit: "none", Date: "unknown"},
			args:    []string{"--debug"},
			wantOut: "Usage:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewRootCmd(tt.info)
			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			require.NoError(t, err)
			assert.Contains(t, buf.String(), tt.wantOut)
		})
	}
}

// ---------------------------------------------------------------------------
// Suite-style tests with setup/teardown
// ---------------------------------------------------------------------------

type RootCmdSuite struct {
	suite.Suite
	cmd *cobra.Command
}

func (s *RootCmdSuite) SetupSuite() {
	// Runs once before all tests in the suite
}

func (s *RootCmdSuite) TearDownSuite() {
	// Runs once after all tests in the suite
}

func (s *RootCmdSuite) SetupTest() {
	// Runs before each test
	s.cmd = NewRootCmd(BuildInfo{Version: "test", Commit: "test", Date: "test"})
}

func (s *RootCmdSuite) TearDownTest() {
	// Runs after each test
}

func (s *RootCmdSuite) TestCompletionDisabled() {
	// Verify completion command is not in the command list
	for _, c := range s.cmd.Commands() {
		s.Assert().NotEqual("completion", c.Name(), "completion command should be disabled")
	}
}

func (s *RootCmdSuite) TestHasShortDescription() {
	s.Assert().NotEmpty(s.cmd.Short)
	s.Assert().NotEmpty(s.cmd.Long)
}

// ---------------------------------------------------------------------------
// 日志标志测试
// ---------------------------------------------------------------------------

func TestNewRootCmd_VerboseFlag(t *testing.T) {
	info := BuildInfo{Version: "dev", Commit: "none", Date: "unknown"}
	cmd := NewRootCmd(info)

	// 验证 --verbose 标志存在
	vFlag := cmd.PersistentFlags().Lookup("verbose")
	require.NotNil(t, vFlag, "--verbose flag should exist")
	assert.Equal(t, "v", vFlag.Shorthand)

	// 验证 --debug 标志存在
	dFlag := cmd.PersistentFlags().Lookup("debug")
	require.NotNil(t, dFlag, "--debug flag should exist")
}

func TestNewRootCmd_PersistentPreRunE_NoFlagSilent(t *testing.T) {
	info := BuildInfo{Version: "dev", Commit: "none", Date: "unknown"}
	cmd := NewRootCmd(info)

	// 默认不传 --verbose/--debug 时 PersistentPreRunE 应正常返回
	require.NotNil(t, cmd.PersistentPreRunE)
	err := cmd.PersistentPreRunE(cmd, nil)
	assert.NoError(t, err)
}

// Run the suite
func TestRootCmdSuite(t *testing.T) {
	suite.Run(t, new(RootCmdSuite))
}
