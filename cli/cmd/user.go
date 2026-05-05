package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/kasuganosora/bangumi.skill/cli/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(userGetCmd)
	userCmd.AddCommand(userMeCmd)
	rootCmd.AddCommand(collectionCmd)
	collectionCmd.AddCommand(collectionListCmd)
	collectionCmd.AddCommand(collectionGetCmd)
	collectionCmd.AddCommand(collectionUpdateCmd)
	collectionCmd.AddCommand(collectionEpisodesCmd)
	collectionCmd.AddCommand(collectionUpdateEpisodeCmd)
	collectionCmd.AddCommand(collectionCharactersCmd)
	collectionCmd.AddCommand(collectionPersonsCmd)
}

// --- user ---

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "查看用户信息",
	Long: `user 命令用于查看 Bangumi 用户信息。

子命令:
  get  - 通过用户名获取用户信息
  me   - 获取当前令牌对应的用户信息`,
}

var userGetCmd = &cobra.Command{
	Use:   "get <用户名>",
	Short: "获取用户信息",
	Long:  "通过用户名获取用户的公开信息。",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		u, err := client.GetUserByName(BackgroundCtx(), args[0])
		if err != nil {
			return err
		}
		return PrintOutput(u, formatUser(u))
	},
}

var userMeCmd = &cobra.Command{
	Use:   "me",
	Short: "获取当前用户信息（需令牌）",
	Long:  "获取当前令牌对应的用户详细信息。",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		u, err := client.GetMe(BackgroundCtx())
		if err != nil {
			return err
		}
		return PrintOutput(u, formatUser(u))
	},
}

func formatUser(u *api.UserDetail) fmt.Stringer {
	return stringerFunc(func() string {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("👤 %s (@%s) [ID:%d]\n", u.Nickname, u.Username, u.ID))
		if u.Sign != "" {
			b.WriteString(fmt.Sprintf("签名: %s\n", u.Sign))
		}
		b.WriteString(fmt.Sprintf("用户组: %s\n", userGroupName(u.UserGroup)))
		return b.String()
	})
}

func userGroupName(g api.UserGroup) string {
	switch g {
	case api.UserGroupAdmin:
		return "管理员"
	case api.UserGroupBangumiAdmin:
		return "Bangumi管理猿"
	case api.UserGroupSkylightAdmin:
		return "天窗管理猿"
	case api.UserGroupMuted:
		return "禁言"
	case api.UserGroupBanned:
		return "禁止访问"
	case api.UserGroupUser:
		return "用户"
	case api.UserGroupWikiUser:
		return "维基人"
	default:
		return "未知"
	}
}

// --- collection ---

var collectionCmd = &cobra.Command{
	Use:   "collection",
	Short: "管理收藏和观看进度 🎯",
	Long: `管理条目收藏和章节观看进度——追番管理的核心功能。

收藏类型:
  1=想看  2=看过  3=在看  4=搁置  5=抛弃

典型追番流程:
  1. collection update --subject "作品名" --type 3         # 标记"在看"
  2. collection update-episode --subject "作品名" --ep 1 --type 2  # 标记第1话看过
  3. collection update --subject "作品名" --type 2 --rate 9       # 看完后评分

子命令:
  list            - 列出用户收藏列表
  get             - 查看某条目的收藏详情
  update          - 修改收藏状态/评分/标签
  episodes        - 查看章节观看进度
  update-episode  - 标记某一话看过/未看`,
}

var collectionListCmd = &cobra.Command{
	Use:   "list [用户名]",
	Short: "列出收藏（默认自己）",
	Long: `列出用户收藏的条目，可按类型和收藏状态筛选。不传默认查看自己的收藏。

筛选示例:
  bangumi collection list                                     # 自己的全部收藏
  bangumi collection list --subject-type 2 --collection-type 3 # 自己在看的动画
  bangumi collection list sai --subject-type 2                 # sai的动画收藏`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		username := ""
		if len(args) == 1 {
			username = args[0]
		} else {
			me, err := client.GetMe(BackgroundCtx())
			if err != nil {
				return fmt.Errorf("请指定用户名或确保已登录: %w", err)
			}
			username = me.Username
		}
		var st *api.SubjectType
		var ct *api.SubjectCollectionType
		if cmd.Flags().Changed("subject-type") {
			v, _ := cmd.Flags().GetInt("subject-type")
			t := api.SubjectType(v)
			st = &t
		}
		if cmd.Flags().Changed("collection-type") {
			v, _ := cmd.Flags().GetInt("collection-type")
			t := api.SubjectCollectionType(v)
			ct = &t
		}
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

		result, err := client.GetUserCollections(BackgroundCtx(), username, st, ct, limit, offset)
		if err != nil {
			return err
		}
		return PrintOutput(result, formatCollections(result))
	},
}

var collectionGetCmd = &cobra.Command{
	Use:   "get <用户名> <条目ID>",
	Short: "查看条目收藏详情",
	Long:  "查看某用户对某条目的详细收藏状态：评分、评论、标签、观看到了第几话等。",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		sid, _ := parseInt(args[1])
		c, err := client.GetUserSubjectCollection(BackgroundCtx(), args[0], sid)
		if err != nil {
			return err
		}
		return PrintOutput(c, formatCollection(c))
	},
}

var collectionUpdateCmd = &cobra.Command{
	Use:   "update [条目ID]",
	Short: "更新条目收藏状态（需令牌）",
	Long: `更新当前用户对指定条目的收藏状态。

位置参数输入作品名称自动搜索，用 --id 精确指定。

参数:
  --type <1|2|3|4|5>   收藏状态: 1=想看 2=看过 3=在看 4=搁置 5=抛弃
  --rate <1-10>         评分
  --comment "文字"      简短评论
  --comment-file <路径> 从文件读取评论（适合长文本）
  --tags "标签1,标签2"   标签（逗号分隔）
  --ep-status <集数>    观看进度
  --private              设为私有

示例:
  bangumi collection update "AIR" --type 3          # 标记"在看"
  bangumi collection update "AIR" --type 2 --rate 9  # 看过 + 9分
  bangumi collection update --id 12 --type 2        # 通过ID`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		var sid int
		if cmd.Flags().Changed("id") {
			sid, _ = cmd.Flags().GetInt("id")
		} else if len(args) == 1 {
			var err error
			sid, err = ResolveSubjectID(client, args[0])
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("请指定作品名称或通过 --id 指定条目ID")
		}

		var req api.UserSubjectCollectionUpdate
		if cmd.Flags().Changed("type") {
			v, _ := cmd.Flags().GetInt("type")
			t := api.SubjectCollectionType(v)
			req.Type = &t
		}
		if cmd.Flags().Changed("rate") {
			v, _ := cmd.Flags().GetInt("rate")
			req.Rate = &v
		}
		commentFile, _ := cmd.Flags().GetString("comment-file")
		commentVal, _ := cmd.Flags().GetString("comment")
		comment, err := FileOrValue(commentFile, commentVal, "comment")
		if err != nil {
			return err
		}
		if comment != "" {
			req.Comment = &comment
		}
		if cmd.Flags().Changed("tags") {
			tags, _ := cmd.Flags().GetString("tags")
			req.Tags = strings.Split(tags, ",")
		}
		if cmd.Flags().Changed("ep-status") {
			v, _ := cmd.Flags().GetInt("ep-status")
			req.EpStatus = &v
		}
		if cmd.Flags().Changed("private") {
			v, _ := cmd.Flags().GetBool("private")
			req.Private = &v
		}

		if err := client.UpdateUserSubjectCollection(BackgroundCtx(), sid, req); err != nil {
			return err
		}
		fmt.Println("✅ 收藏状态已更新")
		return nil
	},
}

var collectionEpisodesCmd = &cobra.Command{
	Use:   "episodes [条目ID]",
	Short: "查看章节观看进度（需令牌）",
	Long: `获取当前用户对指定条目的章节观看进度（每一话是否已看）。

位置参数输入作品名称自动搜索，用 --id 精确指定。

示例:
  bangumi collection episodes "AIR"           # 名称搜索
  bangumi collection episodes --id 12          # ID精确查找`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		var sid int
		if cmd.Flags().Changed("id") {
			sid, _ = cmd.Flags().GetInt("id")
		} else if len(args) == 1 {
			var err error
			sid, err = ResolveSubjectID(client, args[0])
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("请指定作品名称或通过 --id 指定条目ID")
		}
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		eps, err := client.GetUserSubjectEpisodeCollection(BackgroundCtx(), sid, limit, offset)
		if err != nil {
			return err
		}
		return PrintOutput(eps, eps)
	},
}

var collectionUpdateEpisodeCmd = &cobra.Command{
	Use:   "update-episode [章节ID]",
	Short: "更新章节观看状态（需令牌）",
	Long: `更新当前用户对指定章节的观看状态。

位置参数输入作品名称自动搜索，配合 --ep 指定集数。
也可通过 --id 精确指定条目ID，或直接传章节ID。

--type: 0=未看 1=想看 2=看过 3=抛弃

示例:
  bangumi collection update-episode "AIR" --ep 1 --type 2    # 标记AIR第1话看过
  bangumi collection update-episode --id 265 --ep 3 --type 2 # 通过条目ID
  bangumi collection update-episode 1027 --type 2            # 直接指定章节ID`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		t, _ := cmd.Flags().GetInt("type")

		var eid int
		subjID, _ := cmd.Flags().GetInt("id")
		epNum, _ := cmd.Flags().GetInt("ep")

		// 通过 --ep + 名称/ID 自动查找章节
		if epNum > 0 && (len(args) == 1 || subjID > 0) {
			var sid int
			if subjID > 0 {
				sid = subjID
			} else {
				sid, err = ResolveSubjectID(client, args[0])
				if err != nil {
					return err
				}
			}
			eid, err = resolveEpisodeID(client, "", sid, epNum)
			if err != nil {
				return err
			}
		} else if len(args) == 1 {
			// 直接章节ID
			eid, _ = parseInt(args[0])
		} else {
			return fmt.Errorf("请指定章节ID，或用 '作品名 --ep 集数' 自动查找")
		}

		if err := client.UpdateUserEpisodeCollection(BackgroundCtx(), eid, api.EpisodeCollectionType(t)); err != nil {
			return err
		}
		var emoji string
		switch t {
		case 2:
			emoji = "✅ 已标记为看过"
		case 1:
			emoji = "📌 已标记为想看"
		case 0:
			emoji = "↩️ 已取消标记"
		case 3:
			emoji = "🗑️ 已标记为抛弃"
		default:
			emoji = "✅ 已更新"
		}
		fmt.Println(emoji)
		return nil
	},
}

// ── collection characters ──

var collectionCharactersCmd = &cobra.Command{
	Use:   "characters [用户名]",
	Short: "查看收藏的角色（默认自己）",
	Long: `查看用户收藏的角色列表。不传参数默认查看自己。

示例:
  bangumi collection characters              # 自己的收藏
  bangumi collection characters sai          # sai的收藏`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		username := ""
		if len(args) == 1 {
			username = args[0]
		} else {
			me, err := client.GetMe(BackgroundCtx())
			if err != nil {
				return fmt.Errorf("请指定用户名或确保已登录: %w", err)
			}
			username = me.Username
		}
		chars, err := client.GetUserCharacterCollections(BackgroundCtx(), username)
		if err != nil {
			return err
		}
		return PrintOutput(chars, formatCharCollection(chars))
	},
}

// ── collection persons ──

var collectionPersonsCmd = &cobra.Command{
	Use:   "persons [用户名]",
	Short: "查看收藏的人物（默认自己）",
	Long: `查看用户收藏的人物列表（声优/导演等）。不传参数默认查看自己。

示例:
  bangumi collection persons                  # 自己的收藏
  bangumi collection persons sai              # sai的收藏`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		username := ""
		if len(args) == 1 {
			username = args[0]
		} else {
			me, err := client.GetMe(BackgroundCtx())
			if err != nil {
				return fmt.Errorf("请指定用户名或确保已登录: %w", err)
			}
			username = me.Username
		}
		persons, err := client.GetUserPersonCollections(BackgroundCtx(), username)
		if err != nil {
			return err
		}
		return PrintOutput(persons, formatPersonCollection(persons))
	},
}

// resolveEpisodeID 根据作品名/ID + 集数查找章节ID
func resolveEpisodeID(client *api.HTTPClient, subjectName string, subjectID int, epNum int) (int, error) {
	ctx := BackgroundCtx()

	// Step 1: 获取条目ID
	if subjectID == 0 {
		result, err := client.SearchSubjects(ctx, api.SearchSubjectRequest{
			Keyword: subjectName,
			Sort:    "match",
		}, 5, 0)
		if err != nil {
			return 0, fmt.Errorf("搜索作品 '%s' 失败: %w", subjectName, err)
		}
		if len(result.Data) == 0 {
			return 0, fmt.Errorf("未找到作品 '%s'", subjectName)
		}
		subjectID = result.Data[0].ID
	}

	// Step 2: 获取章节列表
	onlyMain := api.EpMainStory
	eps, err := client.GetEpisodes(ctx, subjectID, &onlyMain, 200, 0)
	if err != nil {
		return 0, fmt.Errorf("获取章节列表失败: %w", err)
	}

	// Step 3: 匹配集数
	for _, ep := range eps.Data {
		if int(ep.Ep) == epNum || int(ep.Sort) == epNum {
			return ep.ID, nil
		}
	}
	// 放宽匹配：检查 sort
	for _, ep := range eps.Data {
		if int(ep.Sort) == epNum {
			return ep.ID, nil
		}
	}
	return 0, fmt.Errorf("未找到作品 ID=%d 的第 %d 话 (共 %d 话)", subjectID, epNum, len(eps.Data))
}

func init() {
	collectionListCmd.Flags().Int("subject-type", -1, "条目类型")
	collectionListCmd.Flags().Int("collection-type", -1, "收藏类型 (1=想看 2=看过 3=在看 4=搁置 5=抛弃)")
	collectionListCmd.Flags().Int("limit", 30, "每页数量")
	collectionListCmd.Flags().Int("offset", 0, "偏移")

	collectionUpdateCmd.Flags().Int("id", 0, "条目ID（精确查找）")
	collectionUpdateCmd.Flags().Int("type", -1, "收藏类型")
	collectionUpdateCmd.Flags().Int("rate", -1, "评分 (1-10)")
	collectionUpdateCmd.Flags().String("comment", "", "简短评论")
	collectionUpdateCmd.Flags().String("comment-file", "", "从文件读取评论")
	collectionUpdateCmd.Flags().String("tags", "", "标签（逗号分隔）")
	collectionUpdateCmd.Flags().Int("ep-status", -1, "观看进度（集数）")
	collectionUpdateCmd.Flags().Bool("private", false, "设为私有")

	collectionEpisodesCmd.Flags().Int("id", 0, "条目ID（精确查找）")
	collectionEpisodesCmd.Flags().Int("limit", 100, "每页数量")
	collectionEpisodesCmd.Flags().Int("offset", 0, "偏移")

	collectionUpdateEpisodeCmd.Flags().Int("type", 2, "状态: 0=未看 1=想看 2=看过 3=抛弃")
	collectionUpdateEpisodeCmd.Flags().Int("id", 0, "条目ID（配合 --ep 使用）")
	collectionUpdateEpisodeCmd.Flags().Int("ep", 0, "集数（第几话）")
}

func formatCollections(r *api.Paged[api.UserSubjectCollection]) fmt.Stringer {
	return stringerFunc(func() string {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("收藏列表: 共 %d 项 (显示 %d 项)\n", r.Total, len(r.Data)))
		b.WriteString(strings.Repeat("─", 60) + "\n")
		for i, c := range r.Data {
			name := c.Subject.Name
			if c.Subject.NameCN != "" {
				name += fmt.Sprintf(" (%s)", c.Subject.NameCN)
			}
			b.WriteString(fmt.Sprintf("%d. %s [ID:%d]", i+1, name, c.SubjectID))
			if c.Rate > 0 {
				b.WriteString(fmt.Sprintf(" ⭐%d", c.Rate))
			}
			b.WriteString(fmt.Sprintf(" | %s", collectionTypeName(c.Type)))
			b.WriteString(fmt.Sprintf(" | 进度: %d话\n", c.EpStatus))
		}
		return b.String()
	})
}

func formatCollection(c *api.UserSubjectCollection) fmt.Stringer {
	return stringerFunc(func() string {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("条目: %s", c.Subject.Name))
		if c.Subject.NameCN != "" {
			b.WriteString(fmt.Sprintf(" (%s)", c.Subject.NameCN))
		}
		b.WriteString(fmt.Sprintf("\n收藏状态: %s\n", collectionTypeName(c.Type)))
		if c.Rate > 0 {
			b.WriteString(fmt.Sprintf("评分: %d/10\n", c.Rate))
		}
		b.WriteString(fmt.Sprintf("观看进度: %d 话\n", c.EpStatus))
		if len(c.Tags) > 0 {
			b.WriteString(fmt.Sprintf("标签: %s\n", strings.Join(c.Tags, ", ")))
		}
		if c.Comment != "" {
			b.WriteString(fmt.Sprintf("评论: %s\n", c.Comment))
		}
		return b.String()
	})
}

func collectionTypeName(t api.SubjectCollectionType) string {
	switch t {
	case api.CollectionWish:
		return "想看"
	case api.CollectionDone:
		return "看过"
	case api.CollectionDoing:
		return "在看"
	case api.CollectionOnHold:
		return "搁置"
	case api.CollectionDropped:
		return "抛弃"
	default:
		return "未知"
	}
}

func formatCharCollection(items []api.UserCharacterCollection) fmt.Stringer {
	return stringerFunc(func() string {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("收藏的角色 (%d 项):\n", len(items)))
		b.WriteString(strings.Repeat("─", 40) + "\n")
		for i, c := range items {
			b.WriteString(fmt.Sprintf("%d. %s [ID:%d]  %s\n", i+1, c.Name, c.ID, fmtDate(c.CreatedAt)))
		}
		return b.String()
	})
}

func formatPersonCollection(items []api.UserPersonCollection) fmt.Stringer {
	return stringerFunc(func() string {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("收藏的人物 (%d 项):\n", len(items)))
		b.WriteString(strings.Repeat("─", 40) + "\n")
		for i, p := range items {
			b.WriteString(fmt.Sprintf("%d. %s [ID:%d]  %s\n", i+1, p.Name, p.ID, fmtDate(p.CreatedAt)))
		}
		return b.String()
	})
}

func fmtDate(iso string) string {
	t, err := time.Parse(time.RFC3339, iso)
	if err != nil {
		return iso
	}
	return t.Format("2006-01-02 15:04:05")
}
