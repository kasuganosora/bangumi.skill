package cmd

import (
	"fmt"
	"strings"

	"github.com/kasuganosora/bangumi.skill/cli/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(indexCmd)
	indexCmd.AddCommand(indexCreateCmd)
	indexCmd.AddCommand(indexEditCmd)
	indexCmd.AddCommand(indexSubjectsCmd)
	indexCmd.AddCommand(indexAddSubjectCmd)
	indexCmd.AddCommand(indexEditSubjectCmd)
	indexCmd.AddCommand(indexDeleteSubjectCmd)
	indexCmd.AddCommand(indexCollectCmd)
	indexCmd.AddCommand(indexUncollectCmd)
}

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "管理目录/收藏夹（需令牌）",
	Long: `管理条目目录（类似收藏夹/播放列表功能），可收录多个条目并添加备注。

使用场景:
  - 创建个人推荐列表
  - 整理某个系列的关联作品
  - 收藏别人创建的目录

子命令:
  create         - 创建新目录
  edit           - 编辑目录信息
  subjects       - 查看目录中包含的条目
  add-subject    - 向目录添加条目
  edit-subject   - 编辑目录中条目的备注
  delete-subject - 从目录移除条目
  collect        - 收藏目录
  uncollect      - 取消收藏目录`,
}

// ── create ──

var indexCreateCmd = &cobra.Command{
	Use:   "create <标题>",
	Short: "创建新目录",
	Long: `创建一个新的条目目录，可用于整理推荐列表或关联作品。

示例:
  bangumi index create "2024年度最佳" --desc "个人年度推荐"       # 创建推荐列表
  bangumi index create "物语系列" --desc-file ./description.md     # 从文件读取描述`,

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		descFile, _ := cmd.Flags().GetString("desc-file")
		descVal, _ := cmd.Flags().GetString("desc")
		desc, err := FileOrValue(descFile, descVal, "desc")
		if err != nil {
			return err
		}
		nsfw, _ := cmd.Flags().GetBool("nsfw")

		idx, err := client.NewIndex(BackgroundCtx(), api.NewIndexRequest{
			Title: args[0],
			Desc:  desc,
			NSFW:  nsfw,
		})
		if err != nil {
			return err
		}
		return PrintOutput(idx, formatIndex(idx))
	},
}

// ── edit ──

var indexEditCmd = &cobra.Command{
	Use:   "edit <目录ID>",
	Short: "编辑目录信息",
	Long: `修改目录的标题、描述或NSFW标记。只需要传要修改的字段。

示例:
  bangumi index edit 123 --title "新标题"                        # 仅修改标题
  bangumi index edit 123 --desc "更新后的描述" --nsfw            # 修改描述和标记`,

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := parseInt(args[0])
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		var req api.UpdateIndexRequest
		if cmd.Flags().Changed("title") {
			v, _ := cmd.Flags().GetString("title")
			req.Title = &v
		}
		descFile, _ := cmd.Flags().GetString("desc-file")
		descVal, _ := cmd.Flags().GetString("desc")
		if descFile != "" || descVal != "" {
			desc, err := FileOrValue(descFile, descVal, "desc")
			if err != nil {
				return err
			}
			req.Desc = &desc
		}
		if cmd.Flags().Changed("nsfw") {
			v, _ := cmd.Flags().GetBool("nsfw")
			req.NSFW = &v
		}

		idx, err := client.EditIndexByID(BackgroundCtx(), id, req)
		if err != nil {
			return err
		}
		return PrintOutput(idx, formatIndex(idx))
	},
}

// ── subjects ──

var indexSubjectsCmd = &cobra.Command{
	Use:   "subjects <目录ID>",
	Short: "查看目录条目列表",
	Long:  "查看指定目录中收录的条目列表。",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := parseInt(args[0])
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		subs, err := client.GetIndexSubjects(BackgroundCtx(), id, limit, offset)
		if err != nil {
			return err
		}
		return PrintOutput(subs, subs)
	},
}

// ── add-subject ──

var indexAddSubjectCmd = &cobra.Command{
	Use:   "add-subject <目录ID> <条目ID>",
	Short: "向目录添加条目",
	Long: `向目录中添加一个条目，可选附加备注。

使用场景:
  - 把看过的动画加入推荐列表
  - 把系列作品整理到同一个目录

示例:
  bangumi index add-subject 123 265                           # 添加条目265到目录123
  bangumi index add-subject 123 265 --comment "心目中第一神作"   # 附带备注`,

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		iid, _ := parseInt(args[0])
		sid, _ := parseInt(args[1])
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		comment, _ := cmd.Flags().GetString("comment")
		if err := client.AddIndexSubject(BackgroundCtx(), iid, api.AddIndexSubjectRequest{
			SubjectID: sid,
			Comment:   comment,
		}); err != nil {
			return err
		}
		fmt.Println("✅ 条目已添加到目录")
		return nil
	},
}

// ── edit-subject ──

var indexEditSubjectCmd = &cobra.Command{
	Use:   "edit-subject <目录ID> <条目ID>",
	Short: "编辑目录中条目描述（需令牌）",
	Long:  "修改目录中指定条目的备注描述。",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		iid, _ := parseInt(args[0])
		sid, _ := parseInt(args[1])
		commentFile, _ := cmd.Flags().GetString("comment-file")
		commentVal, _ := cmd.Flags().GetString("comment")
		comment, err := FileOrValue(commentFile, commentVal, "comment")
		if err != nil {
			return err
		}
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		if err := client.EditIndexSubject(BackgroundCtx(), iid, sid, api.EditIndexSubjectRequest{
			Comment: comment,
		}); err != nil {
			return err
		}
		fmt.Println("✅ 条目描述已更新")
		return nil
	},
}

// ── delete-subject ──

var indexDeleteSubjectCmd = &cobra.Command{
	Use:   "delete-subject <目录ID> <条目ID>",
	Short: "从目录删除条目（需令牌）",
	Long:  "从指定目录中移除条目。",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		iid, _ := parseInt(args[0])
		sid, _ := parseInt(args[1])
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		if err := client.DeleteIndexSubject(BackgroundCtx(), iid, sid); err != nil {
			return err
		}
		fmt.Println("✅ 条目已从目录移除")
		return nil
	},
}

// ── collect/uncollect ──

var indexCollectCmd = &cobra.Command{
	Use:   "collect <目录ID>",
	Short: "收藏目录（需令牌）",
	Long:  "收藏指定目录。",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := parseInt(args[0])
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		if err := client.CollectIndex(BackgroundCtx(), id); err != nil {
			return err
		}
		fmt.Println("✅ 目录已收藏")
		return nil
	},
}

var indexUncollectCmd = &cobra.Command{
	Use:   "uncollect <目录ID>",
	Short: "取消收藏目录（需令牌）",
	Long:  "取消收藏指定目录。",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := parseInt(args[0])
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		if err := client.UncollectIndex(BackgroundCtx(), id); err != nil {
			return err
		}
		fmt.Println("✅ 已取消收藏")
		return nil
	},
}

func init() {
	indexCreateCmd.Flags().String("desc", "", "目录描述")
	indexCreateCmd.Flags().String("desc-file", "", "从文件读取描述")
	indexCreateCmd.Flags().Bool("nsfw", false, "标记为NSFW")

	indexEditCmd.Flags().String("title", "", "新标题")
	indexEditCmd.Flags().String("desc", "", "新描述")
	indexEditCmd.Flags().String("desc-file", "", "从文件读取描述")
	indexEditCmd.Flags().Bool("nsfw", false, "NSFW标记")

	indexSubjectsCmd.Flags().Int("limit", 30, "每页数量")
	indexSubjectsCmd.Flags().Int("offset", 0, "偏移")

	indexAddSubjectCmd.Flags().String("comment", "", "备注描述")
	indexEditSubjectCmd.Flags().String("comment", "", "备注描述")
	indexEditSubjectCmd.Flags().String("comment-file", "", "从文件读取备注")
}

func formatIndex(idx *api.Index) fmt.Stringer {
	return stringerFunc(func() string {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("📁 %s [ID:%d]\n", idx.Title, idx.ID))
		if idx.Desc != "" {
			b.WriteString(fmt.Sprintf("描述: %s\n", idx.Desc))
		}
		b.WriteString(fmt.Sprintf("条目数: %d | 收藏: %d | 评论: %d\n",
			idx.Total, idx.Stat.Collects, idx.Stat.Comments))
		b.WriteString(fmt.Sprintf("创建者: %s | NSFW: %v\n", idx.Creator.Nickname, idx.NSFW))
		return b.String()
	})
}
