package cmd

import (
	"fmt"

	"github.com/kasuganosora/bangumi.skill/cli/api"
	"github.com/kasuganosora/bangumi.skill/cli/internal/config"
	"github.com/kasuganosora/bangumi.skill/cli/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authStatusCmd)
	authCmd.AddCommand(authLogoutCmd)
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "管理 Bangumi 认证令牌",
	Long: `首次使用前需要设置个人访问令牌，令牌保存在 token.json 中。

⚠️ 没有令牌时，所有需要认证的命令会提示你先设置令牌。

令牌申请: 浏览器打开 https://next.bgm.tv/demo/access-token 获取。`,
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "设置并验证个人令牌",
	Long: `设置 Bangumi 个人令牌，会自动验证有效性后保存。

令牌保存在当前二进制文件同级目录下的 token.json 中，供所有命令自动读取。

示例:
  bangumi auth login --token "your_token_here"              # 直接指定令牌
  bangumi auth login --token-file ./my_token.txt             # 从文件读取令牌`,
	RunE: func(cmd *cobra.Command, args []string) error {
		tokenVal, _ := cmd.Flags().GetString("token")
		tokenFile, _ := cmd.Flags().GetString("token-file")
		token, err := FileOrValue(tokenFile, tokenVal, "token")
		if err != nil {
			return err
		}
		if token == "" {
			return fmt.Errorf("请通过 --token 或 --token-file 提供令牌\n申请地址: https://next.bgm.tv/demo/access-token")
		}

		// 验证 token 有效性
		client, err := api.NewClient(api.WithAccessToken(token))
		if err != nil {
			return err
		}
		status, err := client.GetTokenStatus(BackgroundCtx(), token)
		if err != nil {
			return fmt.Errorf("令牌无效: %w\n请确认令牌正确，或重新申请: https://next.bgm.tv/demo/access-token", err)
		}

		td := &config.TokenData{
			AccessToken: token,
			UserID:      status.UserID,
		}
		if err := config.SaveToken(td); err != nil {
			return err
		}

		// 获取用户名
		me, err := client.GetMe(BackgroundCtx())
		if err != nil {
			me = &api.UserDetail{} // 获取失败也不影响登录
		}

		out := LoadFormat()
		if outputFormat == "json" {
			_ = out.Print(map[string]interface{}{
				"status":   "ok",
				"user_id":  status.UserID,
				"username": me.Username,
				"nickname": me.Nickname,
				"expires":  status.Expires,
			})
		} else {
			fmt.Printf("✅ 令牌已保存 - %s (ID: %d)\n", me.Nickname, status.UserID)
		}
		log.Info("token saved", "user_id", status.UserID)
		return nil
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "查看令牌是否有效",
	Long:  "检查当前存储的令牌是否有效，显示关联的用户ID和过期时间。",
	RunE: func(cmd *cobra.Command, args []string) error {
		td, err := config.LoadToken()
		if err != nil {
			return err
		}
		if td == nil {
			out := LoadFormat()
			if outputFormat == "json" {
				_ = out.Print(map[string]string{"status": "no_token"})
			} else {
				fmt.Println("未设置令牌。")
				fmt.Println("请申请: https://next.bgm.tv/demo/access-token")
				fmt.Println("然后运行: bangumi auth login --token <令牌>")
			}
			return nil
		}

		client, err := api.NewClient(api.WithAccessToken(td.AccessToken))
		if err != nil {
			return err
		}
		status, err := client.GetTokenStatus(BackgroundCtx(), td.AccessToken)
		if err != nil {
			return fmt.Errorf("令牌无效或已过期: %w", err)
		}
		out := LoadFormat()
		if outputFormat == "json" {
			_ = out.Print(status)
		} else {
			fmt.Printf("令牌有效 ✅\n用户 ID: %d\n客户端 ID: %s\n过期时间戳: %d\n",
				status.UserID, status.ClientID, status.Expires)
		}
		return nil
	},
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "删除已保存的令牌",
	Long:  "删除 token.json 文件，清除本地存储的令牌。",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.DeleteToken(); err != nil {
			return err
		}
		fmt.Println("令牌已删除。")
		log.Info("token deleted")
		return nil
	},
}

func init() {
	authLoginCmd.Flags().String("token", "", "Access Token")
	authLoginCmd.Flags().String("token-file", "", "从文件读取 Access Token")
}
