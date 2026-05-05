package cmd

import (
	"fmt"
	"strings"

	"github.com/kasuganosora/bangumi.skill/cli/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(episodeCmd)
	episodeCmd.AddCommand(episodeListCmd)
	episodeCmd.AddCommand(episodeGetCmd)
	episodeListCmd.Flags().Int("id", 0, "条目ID（精确查找）")
}

var episodeCmd = &cobra.Command{
	Use:   "episode",
	Short: "查看章节/集数信息",
	Long: `查看动画或剧集的章节列表和详情。

使用场景:
  - 追番前查看一共有多少话
  - 获取某一话的标题和简介
  - 标记观看进度前查找章节ID

子命令:
  list  - 列出条目的所有章节
  get   - 获取单个章节详情`,
}

var episodeListCmd = &cobra.Command{
	Use:   "list [条目ID]",
	Short: "列出条目所有章节",
	Long: `列出指定条目的所有章节（集数）。

位置参数输入作品名称自动搜索，用 --id 精确指定。

可选筛选:
  --type         章节类型: 0=本篇 1=SP 2=OP 3=ED
  --limit/--offset 分页

示例:
  bangumi episode list "AIR"                   # 名称搜索
  bangumi episode list --id 12                 # ID精确查找
  bangumi episode list --id 12 --type 0        # 仅本篇`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		sid, err := RequireIDOrName(cmd, args, client, "subject")
		if err != nil {
			return err
		}
		var epType *api.EpType
		if cmd.Flags().Changed("type") {
			t, _ := cmd.Flags().GetInt("type")
			et := api.EpType(t)
			epType = &et
		}
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

		eps, err := client.GetEpisodes(BackgroundCtx(), sid, epType, limit, offset)
		if err != nil {
			return err
		}
		return PrintOutput(eps, formatEpisodes(eps))
	},
}

var episodeGetCmd = &cobra.Command{
	Use:   "get <章节ID>",
	Short: "获取章节详情",
	Long:  "获取指定章节的详细信息（标题/时长/放送日期/简介等）。",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := parseInt(args[0])
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		ep, err := client.GetEpisodeByID(BackgroundCtx(), id)
		if err != nil {
			return err
		}
		return PrintOutput(ep, formatEpisode(ep))
	},
}

func init() {
	episodeListCmd.Flags().Int("type", -1, "章节类型 (0-6)")
	episodeListCmd.Flags().Int("limit", 100, "每页数量")
	episodeListCmd.Flags().Int("offset", 0, "偏移")
}

func formatEpisodes(r *api.Paged[api.EpisodeDetail]) fmt.Stringer {
	return stringerFunc(func() string {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("章节列表: 共 %d 话\n", r.Total))
		b.WriteString(strings.Repeat("─", 50) + "\n")
		for _, e := range r.Data {
			b.WriteString(fmt.Sprintf("第%.0f话 %s", e.Sort, e.Name))
			if e.NameCN != "" {
				b.WriteString(fmt.Sprintf(" (%s)", e.NameCN))
			}
			b.WriteString(fmt.Sprintf("  [ID:%d]", e.ID))
			if e.Airdate != "" {
				b.WriteString(fmt.Sprintf("  %s", e.Airdate))
			}
			b.WriteString(fmt.Sprintf("  %s\n", e.Duration))
		}
		return b.String()
	})
}

func formatEpisode(e *api.EpisodeDetail) fmt.Stringer {
	return stringerFunc(func() string {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("第%.0f话 %s\n", e.Sort, e.Name))
		if e.NameCN != "" {
			b.WriteString(fmt.Sprintf("中文标题: %s\n", e.NameCN))
		}
		b.WriteString(fmt.Sprintf("ID: %d | 条目ID: %d\n", e.ID, e.SubjectID))
		b.WriteString(fmt.Sprintf("类型: %s | 时长: %s\n", epTypeName(e.Type), e.Duration))
		if e.Airdate != "" {
			b.WriteString(fmt.Sprintf("放送日期: %s\n", e.Airdate))
		}
		if e.Desc != "" {
			b.WriteString(fmt.Sprintf("简介: %s\n", e.Desc))
		}
		return b.String()
	})
}

func epTypeName(t api.EpType) string {
	switch t {
	case api.EpMainStory:
		return "本篇"
	case api.EpSP:
		return "SP"
	case api.EpOP:
		return "OP"
	case api.EpED:
		return "ED"
	case api.EpPV:
		return "预告/广告"
	case api.EpMAD:
		return "MAD"
	case api.EpOtherType:
		return "其他"
	default:
		return "未知"
	}
}
