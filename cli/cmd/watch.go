package cmd

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/kasuganosora/bangumi.skill/cli/api"
	"github.com/spf13/cobra"
)

const (
	episodePageSize = 100
	maxEpisodePages = 40
	searchPickLimit = 5
)

func init() {
	collectionCmd.AddCommand(collectionWatchCmd)
	collectionWatchCmd.Flags().Int("id", 0, "条目ID（精确查找）")
	collectionWatchCmd.Flags().Int("ep", 0, "看到第几话（含本话，会标记第 1 话到这一话）")
}

var collectionWatchCmd = &cobra.Command{
	Use:   "watch [作品名称]",
	Short: "标记看到第 N 话",
	Long: `把一部作品标记为看到第 N 话。

会做三件事：
  1. 按作品名匹配条目。多个结果对不上唯一名称时列出候选并停止，不写入
  2. 收藏改为「在看」，进度设为第 N 话。已经是「看过」则保持；进度不会往回调
  3. 把第 1 话到第 N 话的本篇批量标成看过

示例:
  bangumi collection watch "AIR" --ep 3
  bangumi collection watch --id 12 --ep 3`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		epN, _ := cmd.Flags().GetInt("ep")
		if epN <= 0 || epN > 5000 {
			return fmt.Errorf("--ep 需要 1 到 5000 之间的整数，表示看到第几话")
		}
		client, _, err := NewAPIClient()
		if err != nil {
			return err
		}
		subj, err := resolveSubjectArg(cmd, args, client)
		if err != nil {
			return err
		}
		result, err := applyWatch(BackgroundCtx(), client, subj, epN)
		if err != nil {
			return err
		}
		return PrintOutput(result, result)
	},
}

// watchAPI 是「看到第 N 话」需要的接口，便于测试替换。
type watchAPI interface {
	GetUserSubjectCollection(ctx context.Context, username string, subjectID int) (*api.UserSubjectCollection, error)
	UpdateUserSubjectCollection(ctx context.Context, subjectID int, req api.UserSubjectCollectionUpdate) error
	GetEpisodes(ctx context.Context, subjectID int, typ *api.EpType, limit, offset int) (*api.Paged[api.EpisodeDetail], error)
	PatchUserEpisodeCollections(ctx context.Context, subjectID int, episodeIDs []int, typ api.EpisodeCollectionType) error
}

// watchResult 是 watch 命令的结果。
type watchResult struct {
	SubjectID int    `json:"subject_id"`
	Name      string `json:"name"`
	NameCN    string `json:"name_cn,omitempty"`
	Type      int    `json:"type"`
	TypeName  string `json:"type_name"`
	EpStatus  int    `json:"ep_status"`
	KeptDone  bool   `json:"kept_done,omitempty"`
	KeptAhead bool   `json:"kept_ahead,omitempty"`
	Marked    []int  `json:"marked_episodes"`
	Missing   []int  `json:"missing_episodes,omitempty"`
	Note      string `json:"note,omitempty"`
}

func (r watchResult) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "已匹配: %s\n", formatSubjectMatch(api.Subject{ID: r.SubjectID, Name: r.Name, NameCN: r.NameCN}))
	status := r.TypeName
	if r.KeptDone {
		status += "（保持）"
	}
	fmt.Fprintf(&b, "收藏: %s | 进度: 第 %d 话", status, r.EpStatus)
	if r.KeptAhead {
		b.WriteString("（未回退）")
	}
	b.WriteString("\n")
	if r.Note != "" {
		b.WriteString(r.Note)
		b.WriteString("\n")
	}
	if len(r.Marked) > 0 {
		fmt.Fprintf(&b, "已标记看过: %s（%d 集）\n", formatEpisodeSpan(r.Marked), len(r.Marked))
	}
	if len(r.Missing) > 0 {
		fmt.Fprintf(&b, "未找到章节: %s\n", formatEpisodeSpan(r.Missing))
	}
	return b.String()
}

func applyWatch(ctx context.Context, client watchAPI, subj api.Subject, epN int) (watchResult, error) {
	if epN <= 0 || epN > 5000 {
		return watchResult{}, fmt.Errorf("--ep 需要 1 到 5000 之间的整数，表示看到第几话")
	}
	onlyMain := api.EpMainStory
	eps, err := listEpisodes(ctx, client, subj.ID, &onlyMain)
	if err != nil {
		return watchResult{}, err
	}
	marked, missing := selectEpisodesUpTo(eps, epN)
	if len(eps) > 0 && len(marked) == 0 {
		return watchResult{}, fmt.Errorf("作品 ID=%d 有 %d 话本篇，但没有第 1-%d 话，未写入", subj.ID, len(eps), epN)
	}

	current, err := client.GetUserSubjectCollection(ctx, "-", subj.ID)
	if err != nil {
		if !isNotFound(err) {
			return watchResult{}, fmt.Errorf("读取当前收藏失败: %w", err)
		}
		current = nil
	}

	typ := api.CollectionDoing
	epStatus := epN
	keptDone := false
	keptAhead := false
	if current != nil {
		if current.Type == api.CollectionDone {
			typ = api.CollectionDone
			keptDone = true
		}
		if current.EpStatus > epStatus {
			epStatus = current.EpStatus
			keptAhead = true
		}
	}

	ids := make([]int, len(marked))
	nums := make([]int, len(marked))
	for i, item := range marked {
		ids[i] = item.id
		nums[i] = item.num
	}
	// 批量标记会让服务端重算完成度，所以先标章节，再用下面的请求把进度和收藏类型定下来。
	if len(ids) > 0 {
		if err := client.PatchUserEpisodeCollections(ctx, subj.ID, ids, api.EpCollectionDone); err != nil {
			return watchResult{}, fmt.Errorf("标记章节失败: %w", err)
		}
	}
	if err := client.UpdateUserSubjectCollection(ctx, subj.ID, api.UserSubjectCollectionUpdate{
		Type:     &typ,
		EpStatus: &epStatus,
	}); err != nil {
		if len(ids) > 0 {
			return watchResult{}, fmt.Errorf("章节已标记，但更新收藏进度失败: %w", err)
		}
		return watchResult{}, err
	}

	result := watchResult{
		SubjectID: subj.ID,
		Name:      subj.Name,
		NameCN:    subj.NameCN,
		Type:      int(typ),
		TypeName:  collectionTypeName(typ),
		EpStatus:  epStatus,
		KeptDone:  keptDone,
		KeptAhead: keptAhead,
		Marked:    nums,
		Missing:   missing,
	}
	if len(eps) == 0 {
		result.Note = "条目还没有章节数据，只更新了收藏进度"
	}
	return result, nil
}

type numberedEpisode struct {
	num int
	id  int
}

func selectEpisodesUpTo(eps []api.EpisodeDetail, epN int) (marked []numberedEpisode, missing []int) {
	found := map[int]int{}
	for _, ep := range eps {
		n, ok := episodeNumber(ep)
		if !ok || n < 1 || n > epN {
			continue
		}
		if _, exists := found[n]; !exists {
			found[n] = ep.ID
		}
	}
	for n := 1; n <= epN; n++ {
		id, ok := found[n]
		if !ok {
			missing = append(missing, n)
			continue
		}
		marked = append(marked, numberedEpisode{num: n, id: id})
	}
	if len(missing) == epN {
		missing = nil
	}
	return marked, missing
}

func episodeNumber(ep api.EpisodeDetail) (int, bool) {
	if ep.Ep != 0 {
		return wholeEpisodeNumber(ep.Ep)
	}
	return wholeEpisodeNumber(ep.Sort)
}

func wholeEpisodeNumber(v float64) (int, bool) {
	if v <= 0 {
		return 0, false
	}
	n := int(math.Round(v))
	if math.Abs(v-float64(n)) > 0.001 {
		return 0, false
	}
	return n, true
}

type episodeLister interface {
	GetEpisodes(ctx context.Context, subjectID int, typ *api.EpType, limit, offset int) (*api.Paged[api.EpisodeDetail], error)
}

func listEpisodes(ctx context.Context, client episodeLister, subjectID int, typ *api.EpType) ([]api.EpisodeDetail, error) {
	var all []api.EpisodeDetail
	offset := 0
	for page := 0; page < maxEpisodePages; page++ {
		result, err := client.GetEpisodes(ctx, subjectID, typ, episodePageSize, offset)
		if err != nil {
			return nil, fmt.Errorf("获取章节列表失败: %w", err)
		}
		if len(result.Data) == 0 {
			return all, nil
		}
		all = append(all, result.Data...)
		offset += len(result.Data)
		if result.Total > 0 && len(all) >= result.Total {
			return all, nil
		}
		if len(result.Data) < episodePageSize {
			return all, nil
		}
	}
	return nil, fmt.Errorf("章节过多，已停止在第 %d 页", maxEpisodePages)
}

func chooseSubject(query string, results []api.Subject) (api.Subject, error) {
	if len(results) == 0 {
		return api.Subject{}, fmt.Errorf("未找到作品 '%s'", query)
	}
	q := strings.TrimSpace(query)
	var exact []api.Subject
	for _, s := range results {
		if subjectNameEquals(s.Name, q) || subjectNameEquals(s.NameCN, q) {
			exact = append(exact, s)
		}
	}
	switch {
	case len(exact) == 1:
		return exact[0], nil
	case len(exact) > 1:
		return api.Subject{}, ambiguousSubjectError(query, exact)
	case len(results) == 1:
		return results[0], nil
	default:
		return api.Subject{}, ambiguousSubjectError(query, results)
	}
}

func subjectNameEquals(name, query string) bool {
	name = strings.TrimSpace(name)
	if name == "" || query == "" {
		return false
	}
	return strings.EqualFold(name, query)
}

func ambiguousSubjectError(query string, items []api.Subject) error {
	var b strings.Builder
	fmt.Fprintf(&b, "作品名 %q 匹配到多个条目，请用 --id 指定：\n", query)
	for i, s := range items {
		fmt.Fprintf(&b, "%d. %s\n", i+1, formatSubjectMatch(s))
	}
	return fmt.Errorf("%s", strings.TrimRight(b.String(), "\n"))
}

func formatSubjectMatch(s api.Subject) string {
	if s.NameCN != "" && !strings.EqualFold(s.NameCN, s.Name) {
		return fmt.Sprintf("%s (%s) [ID:%d]", s.Name, s.NameCN, s.ID)
	}
	return fmt.Sprintf("%s [ID:%d]", s.Name, s.ID)
}

func formatEpisodeSpan(nums []int) string {
	if len(nums) == 0 {
		return ""
	}
	if len(nums) == 1 {
		return fmt.Sprintf("第 %d 话", nums[0])
	}
	contiguous := true
	for i := 1; i < len(nums); i++ {
		if nums[i] != nums[i-1]+1 {
			contiguous = false
			break
		}
	}
	if contiguous {
		return fmt.Sprintf("第 %d-%d 话", nums[0], nums[len(nums)-1])
	}
	parts := make([]string, len(nums))
	for i, n := range nums {
		parts[i] = fmt.Sprintf("%d", n)
	}
	return "第 " + strings.Join(parts, "、") + " 话"
}

func resolveSubjectForWrite(client *api.HTTPClient, name string) (api.Subject, error) {
	result, err := client.SearchSubjects(BackgroundCtx(), api.SearchSubjectRequest{
		Keyword: name,
		Sort:    "match",
	}, searchPickLimit, 0)
	if err != nil {
		return api.Subject{}, fmt.Errorf("搜索 '%s' 失败: %w", name, err)
	}
	return chooseSubject(name, result.Data)
}

func subjectByID(client *api.HTTPClient, id int) (api.Subject, error) {
	if id <= 0 {
		return api.Subject{}, fmt.Errorf("--id 需要正整数")
	}
	s, err := client.GetSubjectByID(BackgroundCtx(), id)
	if err != nil {
		return api.Subject{}, err
	}
	if s == nil {
		return api.Subject{}, fmt.Errorf("未找到条目 ID %d", id)
	}
	return *s, nil
}

func resolveSubjectArg(cmd *cobra.Command, args []string, client *api.HTTPClient) (api.Subject, error) {
	if cmd.Flags().Changed("id") {
		id, _ := cmd.Flags().GetInt("id")
		return subjectByID(client, id)
	}
	if len(args) == 1 {
		return resolveSubjectForWrite(client, args[0])
	}
	return api.Subject{}, fmt.Errorf("请指定作品名称或通过 --id 指定条目ID")
}
