package cmd

import (
	"fmt"
	"strings"

	"github.com/kasuganosora/bangumi.skill/cli/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.AddCommand(searchSubjectsCmd)
	searchCmd.AddCommand(searchCharactersCmd)
	searchCmd.AddCommand(searchPersonsCmd)
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "搜索条目/角色/人物",
	Long: `搜索 Bangumi 数据库中的动画、角色或现实人物。

子命令:
  subjects   - 搜索动画/书籍/音乐/游戏/三次元条目
  characters - 搜索虚拟角色
  persons    - 搜索现实人物（声优/导演/作家等）`,
}

// --- subjects ---

var searchSubjectsCmd = &cobra.Command{
	Use:   "subjects [关键词]",
	Short: "搜索条目",
	Long: `根据关键词搜索 Bangumi 条目，关键词可选，支持多种筛选和排序。

筛选参数（可组合使用，所有条件为"且"关系）:
  --sort           排序方式: match(相关度)/heat(热度)/rank(排名)/score(评分)，默认 match
  --filter-type    条目类型(可多次指定): 1=书籍 2=动画 3=音乐 4=游戏 6=三次元
  --filter-tag     标签筛选(可多次指定)，如 "原创" "科幻"
  --filter-rating  评分范围，如 ">=6" "<8" (可多次指定)
  --filter-rank    排名范围，如 "<=100" (可多次指定)
  --filter-air-date 播出日期，如 ">=2020-01-01" (可多次指定)
  --limit/--offset  分页参数

使用场景:
  - 用户想看某部动画的详细信息，先搜索获取条目ID
  - 查找高分动画推荐
  - 搜索特定年份的动画

示例:
  bangumi search subjects "AIR"                                  # 基础搜索
  bangumi search subjects --filter-type 2 --filter-air-date ">=2026-05-01" --sort rank  # 新番按排名
  bangumi search subjects --filter-type 2 --sort rank --limit 5  # 高排名动画
  bangumi search subjects "科幻" --filter-rating ">=8"           # 高分科幻作品`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		keyword := ""
		if len(args) == 1 {
			keyword = args[0]
		}
		req := api.SearchSubjectRequest{Keyword: keyword}
		req.Sort, _ = cmd.Flags().GetString("sort")
		types, _ := cmd.Flags().GetIntSlice("filter-type")
		for _, t := range types {
			req.Filter.Type = append(req.Filter.Type, api.SubjectType(t))
		}
		tags, _ := cmd.Flags().GetStringSlice("filter-tag")
		req.Filter.Tag = tags
		airDates, _ := cmd.Flags().GetStringSlice("filter-air-date")
		req.Filter.AirDate = airDates
		ratings, _ := cmd.Flags().GetStringSlice("filter-rating")
		req.Filter.Rating = ratings
		ranks, _ := cmd.Flags().GetStringSlice("filter-rank")
		req.Filter.Rank = ranks
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

		result, err := client.SearchSubjects(BackgroundCtx(), req, limit, offset)
		if err != nil {
			return err
		}
		return PrintOutput(result, formatSubjects(result))
	},
}

// --- characters ---

var searchCharactersCmd = &cobra.Command{
	Use:   "characters <关键词>",
	Short: "搜索角色",
	Long: `根据关键词搜索虚拟角色。

使用场景:
  - 查找某个角色信息前先搜索获取角色ID
  - 确认角色存在于数据库中

示例:
  bangumi search characters "神尾观铃"             # 搜索角色
  bangumi search characters "神尾观铃" --limit 5   # 搜索角色`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		req := api.SearchCharacterRequest{Keyword: args[0]}
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

		result, err := client.SearchCharacters(BackgroundCtx(), req, limit, offset)
		if err != nil {
			return err
		}
		return PrintOutput(result, formatCharacters(result))
	},
}

// --- persons ---

var searchPersonsCmd = &cobra.Command{
	Use:   "persons <关键词>",
	Short: "搜索人物",
	Long: `根据关键词搜索现实人物（声优/导演/漫画家/作家等）。

使用场景:
  - 查找声优配音了哪些角色
  - 查找导演参与过哪些作品

示例:
  bangumi search persons "神尾观铃"                # 搜索声优
  bangumi search persons "神尾观铃" --limit 10      # 搜索导演`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		req := api.SearchPersonRequest{Keyword: args[0]}
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

		result, err := client.SearchPersons(BackgroundCtx(), req, limit, offset)
		if err != nil {
			return err
		}
		return PrintOutput(result, formatPersons(result))
	},
}

func init() {
	// 通用标志
	for _, c := range []*cobra.Command{searchSubjectsCmd, searchCharactersCmd, searchPersonsCmd} {
		c.Flags().Int("limit", 20, "每页数量")
		c.Flags().Int("offset", 0, "分页偏移")
	}
	searchSubjectsCmd.Flags().String("sort", "match", "排序方式: match/heat/rank/score")
	searchSubjectsCmd.Flags().IntSlice("filter-type", nil, "条目类型 (1=book 2=anime 3=music 4=game 6=real)")
	searchSubjectsCmd.Flags().StringSlice("filter-tag", nil, "标签筛选")
	searchSubjectsCmd.Flags().StringSlice("filter-air-date", nil, `播出日期条件，如 ">=2020-01-01"`)
	searchSubjectsCmd.Flags().StringSlice("filter-rating", nil, `评分条件，如 ">=6"`)
	searchSubjectsCmd.Flags().StringSlice("filter-rank", nil, `排名条件，如 "<=100"`)
}

// --- 格式化 ---

type fmtSubjects struct{ *api.Paged[api.Subject] }

func (f fmtSubjects) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("搜索结果: 共 %d 条 (显示 %d 条)\n", f.Total, len(f.Data)))
	b.WriteString(strings.Repeat("─", 60) + "\n")
	for i, s := range f.Data {
		b.WriteString(fmt.Sprintf("%d. %s", i+1, s.Name))
		if s.NameCN != "" {
			b.WriteString(fmt.Sprintf(" (%s)", s.NameCN))
		}
		b.WriteString(fmt.Sprintf("\n   ID: %d | 评分: %.1f | 话数: %d | 日期: %s\n",
			s.ID, s.Rating.Score, s.Eps, s.Date))
	}
	return b.String()
}

type fmtCharacters struct {
	*api.Paged[api.CharacterDetail]
}

func (f fmtCharacters) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("搜索结果: 共 %d 条 (显示 %d 条)\n", f.Total, len(f.Data)))
	b.WriteString(strings.Repeat("─", 60) + "\n")
	for i, c := range f.Data {
		b.WriteString(fmt.Sprintf("%d. %s [ID: %d]", i+1, c.Name, c.ID))
		if c.Gender != "" {
			b.WriteString(fmt.Sprintf(" | 性别: %s", c.Gender))
		}
		b.WriteString("\n")
	}
	return b.String()
}

type fmtPersons struct{ *api.Paged[api.PersonDetail] }

func (f fmtPersons) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("搜索结果: 共 %d 条 (显示 %d 条)\n", f.Total, len(f.Data)))
	b.WriteString(strings.Repeat("─", 60) + "\n")
	for i, p := range f.Data {
		b.WriteString(fmt.Sprintf("%d. %s [ID: %d]", i+1, p.Name, p.ID))
		if len(p.Career) > 0 {
			b.WriteString(fmt.Sprintf(" | 职业: %s", strings.Join(p.Career, "/")))
		}
		b.WriteString("\n")
	}
	return b.String()
}

func formatSubjects(r *api.Paged[api.Subject]) fmtSubjects             { return fmtSubjects{r} }
func formatCharacters(r *api.Paged[api.CharacterDetail]) fmtCharacters { return fmtCharacters{r} }
func formatPersons(r *api.Paged[api.PersonDetail]) fmtPersons          { return fmtPersons{r} }
