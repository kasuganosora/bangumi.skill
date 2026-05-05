package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/kasuganosora/bangumi.skill/cli/api"
	"github.com/kasuganosora/bangumi.skill/cli/internal/config"
	"github.com/kasuganosora/bangumi.skill/cli/internal/output"
	"github.com/kasuganosora/bangumi.skill/cli/log"
	"github.com/spf13/cobra"
)

// 全局标志
var (
	outputFormat string // --output json 时设为 "json"
	tokenFlag    string // --token 直接指定 token
	proxyFlag    string // --proxy 指定 HTTP 代理
)

// AddGlobalFlags 添加全局标志
func AddGlobalFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&outputFormat, "output", "txt",
		`输出格式: txt (默认, 适合AI阅读) 或 json (结构化数据)`)
	cmd.PersistentFlags().StringVar(&tokenFlag, "token", "",
		"手动指定 Access Token（覆盖 token.json）")
	cmd.PersistentFlags().StringVar(&proxyFlag, "proxy", "",
		"HTTP 代理地址 (如 http://127.0.0.1:1081)")
}

// AddFileValueFlags 添加 xxx 和 xxx-file 互斥标志
func AddFileValueFlags(cmd *cobra.Command, name, desc string) (*string, *string) {
	val := cmd.Flags().String(name, "", desc)
	file := cmd.Flags().String(name+"-file", "", desc+"（从文件读取，适合多行/长文本）")
	return val, file
}

// LoadFormat 返回当前格式的 Writer
func LoadFormat() *output.Writer {
	return output.NewWriter(loadFormat())
}
func loadFormat() output.Format {
	if outputFormat == "json" {
		return output.FormatJSON
	}
	return output.FormatTxt
}

// NewAPIClient 创建带 token 和 proxy 的 API 客户端
func NewAPIClient() (*api.HTTPClient, *config.TokenData, error) {
	var opts []api.ClientOption

	// --token 标志优先
	if tokenFlag != "" {
		opts = append(opts, api.WithAccessToken(tokenFlag))
	} else {
		td, err := config.LoadToken()
		if err != nil {
			return nil, nil, err
		}
		if td == nil || td.AccessToken == "" {
			client, _ := api.NewClient(opts...)
			return client, nil, fmt.Errorf(
				"未设置个人令牌\n\n请先申请个人令牌: https://next.bgm.tv/demo/access-token\n然后运行: bangumi auth login --token <你的令牌>",
			)
		}
		opts = append(opts, api.WithAccessToken(td.AccessToken))
	}

	// --proxy 标志优先，否则从 config.json 读取
	if proxyFlag != "" {
		opts = append(opts, api.WithProxy(proxyFlag))
	} else {
		cfg, err := config.LoadConfig()
		if err == nil && cfg.Proxy != "" {
			opts = append(opts, api.WithProxy(cfg.Proxy))
		}
	}

	// 从 config.json 读取超时设置，默认 60s
	cfg, err := config.LoadConfig()
	if err == nil && cfg.Timeout > 0 {
		opts = append(opts, api.WithTimeout(cfg.Timeout))
	} else {
		opts = append(opts, api.WithTimeout(60))
	}

	opts = append(opts, api.WithOnUnauthorized(func() {
		_ = config.DeleteToken()
		fmt.Fprintln(os.Stderr, "\n⚠️ 令牌已失效，已清除本地令牌。")
		fmt.Fprintln(os.Stderr, "请重新设置: bangumi auth login --token <新令牌>")
		fmt.Fprintln(os.Stderr, "令牌申请: https://next.bgm.tv/demo/access-token")
	}))

	client, err := api.NewClient(opts...)
	return client, &config.TokenData{}, err
}

// FileOrValue 取 file 内容或 val
func FileOrValue(filePath, val, flagName string) (string, error) {
	if filePath != "" {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("读取 --%s-file 文件失败: %w", flagName, err)
		}
		return string(data), nil
	}
	return val, nil
}

// BackgroundCtx 返回带超时的 context
func BackgroundCtx() context.Context {
	return context.Background()
}

// Fmt 格式化为 AI 友好文本
func Fmt(s string) string { return s }

// LogErr 输出错误并退出
func LogErr(err error) {
	fmt.Fprintln(os.Stderr, "错误:", err)
	log.Error("command failed", "error", err)
}

// PrintOutput json 时输出原始 API 数据，txt 时输出格式化数据
func PrintOutput(raw, formatted interface{}) error {
	w := output.NewWriter(output.FormatTxt)
	if outputFormat == "json" {
		w = output.NewWriter(output.FormatJSON)
		return w.Print(raw)
	}
	return w.Print(formatted)
}

// ResolveSubjectID 根据名称搜索条目并返回第一个匹配 ID
func ResolveSubjectID(client *api.HTTPClient, name string) (int, error) {
	r, err := client.SearchSubjects(BackgroundCtx(), api.SearchSubjectRequest{Keyword: name, Sort: "match"}, 3, 0)
	if err != nil {
		return 0, fmt.Errorf("搜索 '%s' 失败: %w", name, err)
	}
	if len(r.Data) == 0 {
		return 0, fmt.Errorf("未找到作品 '%s'", name)
	}
	return r.Data[0].ID, nil
}

// ResolveCharacterID 根据名称搜索角色并返回第一个匹配 ID
func ResolveCharacterID(client *api.HTTPClient, name string) (int, error) {
	r, err := client.SearchCharacters(BackgroundCtx(), api.SearchCharacterRequest{Keyword: name}, 3, 0)
	if err != nil {
		return 0, fmt.Errorf("搜索角色 '%s' 失败: %w", name, err)
	}
	if len(r.Data) == 0 {
		return 0, fmt.Errorf("未找到角色 '%s'", name)
	}
	return r.Data[0].ID, nil
}

// ResolvePersonID 根据名称搜索人物并返回第一个匹配 ID
func ResolvePersonID(client *api.HTTPClient, name string) (int, error) {
	r, err := client.SearchPersons(BackgroundCtx(), api.SearchPersonRequest{Keyword: name}, 3, 0)
	if err != nil {
		return 0, fmt.Errorf("搜索人物 '%s' 失败: %w", name, err)
	}
	if len(r.Data) == 0 {
		return 0, fmt.Errorf("未找到人物 '%s'", name)
	}
	return r.Data[0].ID, nil
}

// RequireIDOrName 位置参数默认作为名称搜索，--id 做精确查找
// 用法: bangumi character get "神尾观铃"  或  bangumi character get --id 303
// entity: "subject"/"character"/"person"
func RequireIDOrName(cmd *cobra.Command, args []string, client *api.HTTPClient, entity string) (int, error) {
	// --id 精确查找
	if cmd.Flags().Changed("id") {
		id, _ := cmd.Flags().GetInt("id")
		if id <= 0 {
			return 0, fmt.Errorf("--id 需要正整数")
		}
		return id, nil
	}
	// 位置参数作为名称搜索
	if len(args) == 1 {
		name := args[0]
		switch entity {
		case "subject":
			return ResolveSubjectID(client, name)
		case "character":
			return ResolveCharacterID(client, name)
		case "person":
			return ResolvePersonID(client, name)
		default:
			return 0, fmt.Errorf("不支持的实体类型: %s", entity)
		}
	}
	return 0, fmt.Errorf("请指定名称（如 'AIR'）或通过 --id 指定 ID")
}
