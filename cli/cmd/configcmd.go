package cmd

import (
	"fmt"
	"strconv"

	"github.com/kasuganosora/bangumi.skill/cli/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configProxyCmd)
	configCmd.AddCommand(configTimeoutCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "管理 CLI 配置（代理/超时）",
	Long:  `config 命令用于管理 CLI 全局配置：HTTP 代理和请求超时。`,
}

var configProxyCmd = &cobra.Command{
	Use:   "proxy [代理地址]",
	Short: "查看或设置 HTTP 代理",
	Long: `查看或设置 HTTP 代理地址。

无参数时显示当前代理设置。
有参数时设置代理（如 http://127.0.0.1:1081），传 "none" 清除代理。

示例:
  bangumi config proxy                       # 查看当前代理
  bangumi config proxy http://127.0.0.1:1081  # 设置代理
  bangumi config proxy none                   # 清除代理`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}
		if len(args) == 0 {
			showConfig(cfg)
			return nil
		}
		val := args[0]
		if val == "none" || val == "" {
			cfg.Proxy = ""
			if err := config.SaveConfig(cfg); err != nil {
				return err
			}
			fmt.Println("代理已清除。")
		} else {
			cfg.Proxy = val
			if err := config.SaveConfig(cfg); err != nil {
				return err
			}
			fmt.Printf("代理已设置为: %s\n", val)
		}
		return nil
	},
}

var configTimeoutCmd = &cobra.Command{
	Use:   "timeout [秒数]",
	Short: "查看或设置请求超时",
	Long: `查看或设置 API 请求超时时间（秒），默认 60 秒。

无参数时显示当前超时设置。

示例:
  bangumi config timeout      # 查看当前超时
  bangumi config timeout 30   # 设置为 30 秒
  bangumi config timeout 120  # 设置为 120 秒`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}
		if len(args) == 0 {
			showConfig(cfg)
			return nil
		}
		sec, err := strconv.Atoi(args[0])
		if err != nil || sec <= 0 {
			return fmt.Errorf("无效的超时值: %s（需为正整数秒）", args[0])
		}
		cfg.Timeout = sec
		if err := config.SaveConfig(cfg); err != nil {
			return err
		}
		fmt.Printf("请求超时已设置为: %d 秒\n", sec)
		return nil
	},
}

func showConfig(cfg *config.ConfigData) {
	if cfg.Proxy != "" {
		fmt.Printf("代理: %s\n", cfg.Proxy)
	} else {
		fmt.Println("代理: 未设置")
	}
	if cfg.Timeout > 0 {
		fmt.Printf("超时: %d 秒\n", cfg.Timeout)
	} else {
		fmt.Println("超时: 60 秒（默认）")
	}
}
