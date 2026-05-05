package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kasuganosora/bangumi.skill/cli/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(subjectCmd)
	for _, c := range []*cobra.Command{subjectGetCmd, subjectRelationsCmd, subjectCharactersCmd, subjectPersonsCmd} {
		subjectCmd.AddCommand(c)
		c.Flags().Int("id", 0, "条目ID（精确查找）")
	}
}

var subjectCmd = &cobra.Command{
	Use:   "subject",
	Short: "查看条目详情与关联信息",
	Long: `查看动画/书籍/音乐/游戏等条目的详细信息、关联作品、角色和制作人员。

位置参数输入作品名称自动搜索，用 --id 精确指定。

子命令:
  get        - 获取条目基本信息
  relations  - 获取条目关联作品
  characters - 获取条目角色列表
  persons    - 获取条目制作人员`,
}

var subjectGetCmd = &cobra.Command{
	Use:   "get [作品名称]",
	Short: "获取条目详情",
	Long: `获取指定条目的完整信息（名称/简介/评分/标签等）。

示例:
  bangumi subject get "AIR"            # 名称搜索
  bangumi subject get --id 12          # ID精确查找`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		id, err := RequireIDOrName(cmd, args, client, "subject")
		if err != nil {
			return err
		}
		s, err := client.GetSubjectByID(BackgroundCtx(), id)
		if err != nil {
			return err
		}
		return PrintOutput(s, formatSubjectDetail(s))
	},
}

var subjectRelationsCmd = &cobra.Command{
	Use:   "relations [条目ID]",
	Short: "获取条目关联作品",
	Long:  "获取指定条目的关联作品列表。支持 --name 指定作品名。",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		id, err := RequireIDOrName(cmd, args, client, "subject")
		if err != nil {
			return err
		}
		rels, err := client.GetSubjectRelations(BackgroundCtx(), id)
		if err != nil {
			return err
		}
		return PrintOutput(rels, formatRelations(rels))
	},
}

var subjectCharactersCmd = &cobra.Command{
	Use:   "characters [条目ID]",
	Short: "获取条目角色列表",
	Long:  "获取指定条目角色列表。支持 --name 指定作品名。",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		id, err := RequireIDOrName(cmd, args, client, "subject")
		if err != nil {
			return err
		}
		chars, err := client.GetSubjectCharacters(BackgroundCtx(), id)
		if err != nil {
			return err
		}
		return PrintOutput(chars, formatSubjectChars(chars))
	},
}

var subjectPersonsCmd = &cobra.Command{
	Use:   "persons [条目ID]",
	Short: "获取条目制作人员",
	Long:  "获取指定条目制作人员。支持 --name 指定作品名。",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		id, err := RequireIDOrName(cmd, args, client, "subject")
		if err != nil {
			return err
		}
		persons, err := client.GetSubjectPersons(BackgroundCtx(), id)
		if err != nil {
			return err
		}
		return PrintOutput(persons, formatSubjectPersons(persons))
	},
}

// --- 格式化 ---

func formatSubjectDetail(s *api.Subject) fmt.Stringer {
	return stringerFunc(func() string {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("%s\n", s.Name))
		if s.NameCN != "" {
			b.WriteString(fmt.Sprintf("中文名: %s\n", s.NameCN))
		}
		b.WriteString(fmt.Sprintf("ID: %d | 类型: %s", s.ID, subjectTypeName(s.Type)))
		if s.Date != "" {
			b.WriteString(fmt.Sprintf(" | 日期: %s", s.Date))
		}
		b.WriteString(fmt.Sprintf("\n话数: %d | 评分: %.1f (%d人)\n",
			s.Eps, s.Rating.Score, s.Rating.Total))
		b.WriteString(fmt.Sprintf("收藏: 想看%d 在看%d 看过%d 搁置%d 抛弃%d\n",
			s.Collection.Wish, s.Collection.Doing, s.Collection.Collect,
			s.Collection.OnHold, s.Collection.Dropped))
		if len(s.Tags) > 0 {
			var tags []string
			for _, t := range s.Tags {
				tags = append(tags, t.Name)
			}
			b.WriteString(fmt.Sprintf("标签: %s\n", strings.Join(tags, ", ")))
		}
		if s.Summary != "" {
			summary := []rune(s.Summary)
			if len(summary) > 200 {
				summary = append(summary[:200], []rune("...")...)
			}
			fmt.Fprintf(&b, "\n简介: %s\n", string(summary))
		}
		return b.String()
	})
}

func formatRelations(rels []api.V0SubjectRelation) fmt.Stringer {
	return stringerFunc(func() string {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("关联作品 (%d 项):\n", len(rels)))
		for i, r := range rels {
			b.WriteString(fmt.Sprintf("%d. [%s] %s", i+1, r.Relation, r.Name))
			if r.NameCN != "" {
				b.WriteString(fmt.Sprintf(" (%s)", r.NameCN))
			}
			b.WriteString(fmt.Sprintf(" [ID:%d]\n", r.ID))
		}
		return b.String()
	})
}

func formatSubjectChars(chars []api.RelatedCharacter) fmt.Stringer {
	return stringerFunc(func() string {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("角色列表 (%d 项):\n", len(chars)))
		for i, c := range chars {
			b.WriteString(fmt.Sprintf("%d. %s [ID:%d] - %s", i+1, c.Name, c.ID, c.Relation))
			if len(c.Actors) > 0 {
				var actors []string
				for _, a := range c.Actors {
					actors = append(actors, a.Name)
				}
				b.WriteString(fmt.Sprintf(" | CV: %s", strings.Join(actors, ", ")))
			}
			b.WriteString("\n")
		}
		return b.String()
	})
}

func formatSubjectPersons(persons []api.RelatedPerson) fmt.Stringer {
	return stringerFunc(func() string {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("制作人员 (%d 项):\n", len(persons)))
		for i, p := range persons {
			b.WriteString(fmt.Sprintf("%d. %s [ID:%d] - %s", i+1, p.Name, p.ID, p.Relation))
			if len(p.Career) > 0 {
				b.WriteString(fmt.Sprintf(" (%s)", strings.Join(p.Career, "/")))
			}
			b.WriteString("\n")
		}
		return b.String()
	})
}

func subjectTypeName(t api.SubjectType) string {
	switch t {
	case api.SubjectBook:
		return "书籍"
	case api.SubjectAnime:
		return "动画"
	case api.SubjectMusic:
		return "音乐"
	case api.SubjectGame:
		return "游戏"
	case api.SubjectReal:
		return "三次元"
	default:
		return "未知"
	}
}

// 简单辅助
func parseInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}

type stringerFunc func() string

func (f stringerFunc) String() string { return f() }
func (f stringerFunc) MarshalJSON() ([]byte, error) {
	return json.Marshal(f())
}
