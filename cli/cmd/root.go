package cmd

import (
	"fmt"

	"github.com/kasuganosora/bangumi.skill/cli/log"
	"github.com/spf13/cobra"
)

// BuildInfo holds version metadata injected at build time.
type BuildInfo struct {
	Version string
	Commit  string
	Date    string
}

var rootCmd = &cobra.Command{
	Use:   "bangumi",
	Short: "Bangumi CLI - 番组管理命令行工具",
	Long: `Bangumi CLI 是一个番组管理命令行工具,
提供番剧信息查询、观看进度跟踪、评分管理等功能。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var (
	flagVerbose bool
	flagDebug   bool
)

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{Use: "no-help", Hidden: true})

	rootCmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "显示详细日志")
	rootCmd.PersistentFlags().BoolVar(&flagDebug, "debug", false, "显示调试日志（含详细日志）")
	AddGlobalFlags(rootCmd)
}

// NewRootCmd 配置 rootCmd 并返回
func NewRootCmd(info BuildInfo) *cobra.Command {
	rootCmd.Version = info.Version
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		switch {
		case flagDebug:
			log.SetLevel(log.LevelDebug)
		case flagVerbose:
			log.SetLevel(log.LevelInfo)
		}
		return nil
	}

	rootCmd.SetVersionTemplate(fmt.Sprintf(
		"bangumi %s (commit: %s, built: %s)\n",
		info.Version, info.Commit, info.Date,
	))

	return rootCmd
}
