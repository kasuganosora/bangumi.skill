package cmd

import (
	"fmt"
	"strings"

	"github.com/kasuganosora/bangumi.skill/cli/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(characterCmd)
	for _, c := range []*cobra.Command{characterGetCmd, characterSubjectsCmd, characterPersonsCmd, characterCollectCmd, characterUncollectCmd} {
		characterCmd.AddCommand(c)
		c.Flags().Int("id", 0, "角色ID（精确查找，跳过名称搜索）")
	}
}

var characterCmd = &cobra.Command{
	Use:   "character",
	Short: "查看角色信息",
	Long: `查看虚拟角色详情、出演条目和声优信息。

所有子命令通过位置参数输入角色名称自动搜索，也可用 --id 精确指定。

子命令:
  get        - 获取角色详情
  subjects   - 获取角色出演条目
  persons    - 获取角色声优
  collect    - 收藏角色
  uncollect  - 取消收藏角色`,
}

var characterGetCmd = &cobra.Command{
	Use:   "get [角色名称]",
	Short: "获取角色详情",
	Long: `获取指定角色的详细信息（性别/生日/简介等）。

示例:
  bangumi character get "神尾观铃"        # 名称搜索
  bangumi character get --id 303           # ID精确查找`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		id, err := RequireIDOrName(cmd, args, client, "character")
		if err != nil {
			return err
		}
		c, err := client.GetCharacterByID(BackgroundCtx(), id)
		if err != nil {
			return err
		}
		return LoadFormat().Print(formatCharacter(c))
	},
}

var characterSubjectsCmd = &cobra.Command{
	Use:   "subjects [角色名称]",
	Short: "获取角色出演条目",
	Long: `获取指定角色出演的所有条目列表。

示例:
  bangumi character subjects "神尾观铃"    # 名称搜索
  bangumi character subjects --id 303       # ID精确查找`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		id, err := RequireIDOrName(cmd, args, client, "character")
		if err != nil {
			return err
		}
		subs, err := client.GetCharacterSubjects(BackgroundCtx(), id)
		if err != nil {
			return err
		}
		return LoadFormat().Print(formatCharSubjects(subs))
	},
}

var characterPersonsCmd = &cobra.Command{
	Use:   "persons [角色名称]",
	Short: "获取角色声优",
	Long: `获取指定角色的声优（CV）列表。

示例:
  bangumi character persons "神尾观铃"     # 名称搜索
  bangumi character persons --id 303        # ID精确查找`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		id, err := RequireIDOrName(cmd, args, client, "character")
		if err != nil {
			return err
		}
		persons, err := client.GetCharacterPersons(BackgroundCtx(), id)
		if err != nil {
			return err
		}
		return LoadFormat().Print(formatCharPersons(persons))
	},
}

func formatCharacter(c *api.CharacterDetail) fmt.Stringer {
	return stringerFunc(func() string {
		var b strings.Builder
		fmt.Fprintf(&b, "🎭 %s [ID:%d]\n", c.Name, c.ID)
		if c.Gender != "" {
			fmt.Fprintf(&b, "性别: %s", c.Gender)
		}
		if c.BloodType != nil {
			fmt.Fprintf(&b, " | 血型: %s", bloodName(*c.BloodType))
		}
		b.WriteString("\n")
		if c.BirthYear != nil {
			fmt.Fprintf(&b, "生日: %d-%d-%d\n", *c.BirthYear, *c.BirthMon, *c.BirthDay)
		}
		fmt.Fprintf(&b, "收藏数: %d | 评论数: %d\n", c.Stat.Collects, c.Stat.Comments)
		if c.Summary != "" {
			fmt.Fprintf(&b, "简介: %s\n", c.Summary)
		}
		return b.String()
	})
}

func formatCharSubjects(subs []api.V0RelatedSubject) fmt.Stringer {
	return stringerFunc(func() string {
		var b strings.Builder
		fmt.Fprintf(&b, "出演条目 (%d 项):\n", len(subs))
		for i, s := range subs {
			fmt.Fprintf(&b, "%d. %s", i+1, s.Name)
			if s.NameCN != "" {
				fmt.Fprintf(&b, " (%s)", s.NameCN)
			}
			b.WriteString(fmt.Sprintf(" [ID:%d]", s.ID))
			if s.Staff != "" {
				b.WriteString(fmt.Sprintf(" - %s", s.Staff))
			}
			b.WriteString("\n")
		}
		return b.String()
	})
}

func formatCharPersons(persons []api.CharacterPerson) fmt.Stringer {
	return stringerFunc(func() string {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("声优列表 (%d 项):\n", len(persons)))
		for i, p := range persons {
			b.WriteString(fmt.Sprintf("%d. %s [ID:%d]", i+1, p.Name, p.ID))
			b.WriteString(fmt.Sprintf(" | 角色: %s", p.SubjectName))
			if p.SubjectNameCN != "" {
				b.WriteString(fmt.Sprintf(" (%s)", p.SubjectNameCN))
			}
			b.WriteString("\n")
		}
		return b.String()
	})
}

// --- collect / uncollect ---

var characterCollectCmd = &cobra.Command{
	Use:   "collect [角色名称]",
	Short: "收藏角色（需令牌）",
	Long:  "收藏指定角色到个人收藏夹。名称或 --id 二选一。",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		id, err := RequireIDOrName(cmd, args, client, "character")
		if err != nil {
			return err
		}
		if err := client.CollectCharacter(BackgroundCtx(), id); err != nil {
			return err
		}
		fmt.Println("✅ 角色已收藏")
		return nil
	},
}

var characterUncollectCmd = &cobra.Command{
	Use:   "uncollect [角色ID]",
	Short: "取消收藏角色（需令牌）",
	Long:  "取消收藏指定角色。支持 --name 指定角色名。",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		id, err := RequireIDOrName(cmd, args, client, "character")
		if err != nil {
			return err
		}
		if err := client.UncollectCharacter(BackgroundCtx(), id); err != nil {
			return err
		}
		fmt.Println("✅ 已取消收藏")
		return nil
	},
}

func bloodName(b api.BloodType) string {
	switch b {
	case api.BloodA:
		return "A"
	case api.BloodB:
		return "B"
	case api.BloodAB:
		return "AB"
	case api.BloodO:
		return "O"
	default:
		return "?"
	}
}
