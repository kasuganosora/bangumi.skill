package cmd

import (
	"fmt"
	"strings"

	"github.com/kasuganosora/bangumi.skill/cli/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(personCmd)
	for _, c := range []*cobra.Command{personGetCmd, personSubjectsCmd, personCharactersCmd, personCollectCmd, personUncollectCmd} {
		personCmd.AddCommand(c)
		c.Flags().Int("id", 0, "人物ID（精确查找）")
	}
}

var personCmd = &cobra.Command{
	Use:   "person",
	Short: "查看人物信息",
	Long: `查看现实人物（声优/导演/艺术家等）详情和关联作品。

所有子命令通过位置参数输入人物名称自动搜索，也可用 --id 精确指定。

子命令:
  get        - 获取人物详情
  subjects   - 获取人物参与作品
  characters - 获取人物配音角色
  collect    - 收藏人物
  uncollect  - 取消收藏人物`,
}

var personGetCmd = &cobra.Command{
	Use:   "get [人物名称]",
	Short: "获取人物详情",
	Long: `获取指定人物的详细信息。

示例:
  bangumi person get "神尾观铃"           # 名称搜索
  bangumi person get --id 1               # ID精确查找`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		id, err := RequireIDOrName(cmd, args, client, "person")
		if err != nil {
			return err
		}
		ctx := BackgroundCtx()

		// 并行获取人物详情 + 参与作品 + 配音角色
		type result struct {
			p     *api.PersonDetail
			sub   []api.V0RelatedSubject
			chars []api.CharacterPerson
		}
		ch := make(chan result, 1)
		go func() {
			var r result
			r.p, _ = client.GetPersonByID(ctx, id)
			r.sub, _ = client.GetPersonSubjects(ctx, id)
			r.chars, _ = client.GetPersonCharacters(ctx, id)
			ch <- r
		}()
		r := <-ch

		if r.p == nil {
			return fmt.Errorf("未找到人物 (ID: %d)", id)
		}

		full := &api.PersonFull{
			Detail:     *r.p,
			Subjects:   r.sub,
			Characters: r.chars,
		}
		return PrintOutput(full, formatPersonFull(full))
	},
}

var personSubjectsCmd = &cobra.Command{
	Use:   "subjects [人物名称]",
	Short: "获取人物参与作品",
	Long: `获取指定人物参与的所有作品列表。

示例:
  bangumi person subjects "神尾观铃"      # 名称搜索`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		id, err := RequireIDOrName(cmd, args, client, "person")
		if err != nil {
			return err
		}
		subs, err := client.GetPersonSubjects(BackgroundCtx(), id)
		if err != nil {
			return err
		}
		return PrintOutput(subs, formatCharSubjects(subs))
	},
}

var personCharactersCmd = &cobra.Command{
	Use:   "characters [人物名称]",
	Short: "获取人物配音角色",
	Long: `获取指定人物配音的所有角色列表。

示例:
  bangumi person characters "神尾观铃"    # 名称搜索`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		id, err := RequireIDOrName(cmd, args, client, "person")
		if err != nil {
			return err
		}
		chars, err := client.GetPersonCharacters(BackgroundCtx(), id)
		if err != nil {
			return err
		}
		return PrintOutput(chars, formatCharPersons(chars))
	},
}

func formatPersonFull(f *api.PersonFull) fmt.Stringer {
	return stringerFunc(func() string {
		var b strings.Builder
		p := f.Detail

		// ── 基本信息 ──
		fmt.Fprintf(&b, "👤 %s [ID:%d]\n", p.Name, p.ID)
		if len(p.Career) > 0 {
			fmt.Fprintf(&b, "职业: %s\n", strings.Join(p.Career, "/"))
		}
		if p.Gender != "" {
			fmt.Fprintf(&b, "性别: %s", p.Gender)
		}
		if p.BloodType != nil {
			fmt.Fprintf(&b, " | 血型: %s", bloodName(*p.BloodType))
		}
		b.WriteString("\n")
		if p.BirthYear != nil {
			fmt.Fprintf(&b, "生日: %d-%d-%d\n", *p.BirthYear, *p.BirthMon, *p.BirthDay)
		}
		fmt.Fprintf(&b, "收藏数: %d | 评论数: %d\n", p.Stat.Collects, p.Stat.Comments)
		if p.Summary != "" {
			fmt.Fprintf(&b, "简介: %s\n", p.Summary)
		}

		// ── 参与作品 ──
		if len(f.Subjects) > 0 {
			b.WriteString("\n")
			fmt.Fprintf(&b, "参与作品 (%d 项):\n", len(f.Subjects))
			for i, s := range f.Subjects {
				fmt.Fprintf(&b, "  %d. %s", i+1, s.Name)
				if s.NameCN != "" {
					fmt.Fprintf(&b, " (%s)", s.NameCN)
				}
				fmt.Fprintf(&b, " [ID:%d]", s.ID)
				if s.Staff != "" {
					fmt.Fprintf(&b, " - %s", s.Staff)
				}
				b.WriteString("\n")
			}
		}

		// ── 配音角色 ──
		if len(f.Characters) > 0 {
			b.WriteString("\n")
			fmt.Fprintf(&b, "配音角色 (%d 人):\n", len(f.Characters))
			for i, c := range f.Characters {
				fmt.Fprintf(&b, "  %d. %s [ID:%d] | 出自: %s", i+1, c.Name, c.ID, c.SubjectName)
				if c.SubjectNameCN != "" {
					fmt.Fprintf(&b, " (%s)", c.SubjectNameCN)
				}
				b.WriteString("\n")
			}
		}

		return b.String()
	})
}

// --- collect / uncollect ---

var personCollectCmd = &cobra.Command{
	Use:   "collect [人物ID]",
	Short: "收藏人物（需令牌）",
	Long:  "收藏指定人物。支持 --name 指定人物名。",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		id, err := RequireIDOrName(cmd, args, client, "person")
		if err != nil {
			return err
		}
		if err := client.CollectPerson(BackgroundCtx(), id); err != nil {
			return err
		}
		fmt.Println("✅ 人物已收藏")
		return nil
	},
}

var personUncollectCmd = &cobra.Command{
	Use:   "uncollect [人物ID]",
	Short: "取消收藏人物（需令牌）",
	Long:  "取消收藏指定人物。支持 --name 指定人物名。",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		id, err := RequireIDOrName(cmd, args, client, "person")
		if err != nil {
			return err
		}
		if err := client.UncollectPerson(BackgroundCtx(), id); err != nil {
			return err
		}
		fmt.Println("✅ 已取消收藏")
		return nil
	},
}
