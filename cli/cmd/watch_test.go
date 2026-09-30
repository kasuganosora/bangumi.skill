package cmd

import (
	"context"
	"fmt"
	"testing"

	"github.com/kasuganosora/bangumi.skill/cli/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChooseSubject(t *testing.T) {
	air := api.Subject{ID: 12, Name: "AIR", NameCN: "AIR"}
	gear := api.Subject{ID: 99, Name: "AIR GEAR", NameCN: "飞轮少年"}
	tests := []struct {
		name    string
		query   string
		items   []api.Subject
		wantID  int
		wantErr string
	}{
		{name: "exact among many", query: "AIR", items: []api.Subject{gear, air}, wantID: 12},
		{name: "case fold", query: "air", items: []api.Subject{air, gear}, wantID: 12},
		{name: "single fuzzy", query: "飞轮", items: []api.Subject{gear}, wantID: 99},
		{name: "ambiguous fuzzy", query: "air", items: []api.Subject{gear, {ID: 3, Name: "Aria"}}, wantErr: "匹配到多个条目"},
		{name: "two exact", query: "AIR", items: []api.Subject{air, {ID: 13, Name: "AIR", NameCN: "AIR"}}, wantErr: "匹配到多个条目"},
		{name: "empty", query: "没有", wantErr: "未找到作品"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := chooseSubject(tt.query, tt.items)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantID, got.ID)
		})
	}
}

func TestEpisodeNumber(t *testing.T) {
	n, ok := episodeNumber(api.EpisodeDetail{Ep: 3, Sort: 9})
	assert.True(t, ok)
	assert.Equal(t, 3, n)

	n, ok = episodeNumber(api.EpisodeDetail{Sort: 4})
	assert.True(t, ok)
	assert.Equal(t, 4, n)

	_, ok = episodeNumber(api.EpisodeDetail{Ep: 1.5, Sort: 2})
	assert.False(t, ok)

	_, ok = episodeNumber(api.EpisodeDetail{})
	assert.False(t, ok)
}

func TestListEpisodes_Paginates(t *testing.T) {
	eps := make([]api.EpisodeDetail, 150)
	for i := range eps {
		eps[i] = api.EpisodeDetail{ID: i + 1, Ep: float64(i + 1)}
	}
	fake := &fakeWatch{episodes: eps}
	got, err := listEpisodes(context.Background(), fake, 12, nil)
	require.NoError(t, err)
	assert.Len(t, got, 150)
	assert.GreaterOrEqual(t, fake.episodeCalls, 2)
}

func TestApplyWatch_MarksRangeAndSetsDoing(t *testing.T) {
	fake := &fakeWatch{
		getErr:   fmt.Errorf("missing: %w", &api.APIError{StatusCode: 404, Message: "not found"}),
		episodes: mainEpisodes(1, 12),
	}
	subj := api.Subject{ID: 12, Name: "AIR", NameCN: "AIR"}
	result, err := applyWatch(context.Background(), fake, subj, 3)
	require.NoError(t, err)
	assert.Equal(t, 3, result.Type)
	assert.Equal(t, 3, result.EpStatus)
	assert.Equal(t, []int{1, 2, 3}, result.Marked)
	assert.Empty(t, result.Missing)
	require.NotNil(t, fake.updated.Type)
	assert.Equal(t, api.CollectionDoing, *fake.updated.Type)
	assert.Equal(t, []int{101, 102, 103}, fake.patched)
	assert.Contains(t, result.String(), "已匹配: AIR [ID:12]")
	assert.Contains(t, result.String(), "在看")
}

func TestApplyWatch_KeepsDoneAndDoesNotRewind(t *testing.T) {
	fake := &fakeWatch{
		current: &api.UserSubjectCollection{
			Type:     api.CollectionDone,
			EpStatus: 20,
		},
		episodes: mainEpisodes(1, 20),
	}
	result, err := applyWatch(context.Background(), fake, api.Subject{ID: 1, Name: "X"}, 5)
	require.NoError(t, err)
	assert.True(t, result.KeptDone)
	assert.True(t, result.KeptAhead)
	assert.Equal(t, 20, result.EpStatus)
	assert.Equal(t, api.CollectionDone, *fake.updated.Type)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result.Marked)
}

func TestApplyWatch_ReportsMissingEpisodes(t *testing.T) {
	fake := &fakeWatch{
		getErr:   &api.APIError{StatusCode: 404},
		episodes: []api.EpisodeDetail{{ID: 1, Ep: 1}, {ID: 3, Ep: 3}},
	}
	result, err := applyWatch(context.Background(), fake, api.Subject{ID: 1, Name: "X"}, 3)
	require.NoError(t, err)
	assert.Equal(t, []int{1, 3}, result.Marked)
	assert.Equal(t, []int{2}, result.Missing)
	assert.Equal(t, []int{1, 3}, fake.patched)
	assert.Contains(t, result.String(), "未找到章节: 第 2 话")
}

func TestApplyWatch_RefusesWhenNumberingMissing(t *testing.T) {
	fake := &fakeWatch{episodes: mainEpisodes(8, 10)}
	_, err := applyWatch(context.Background(), fake, api.Subject{ID: 1, Name: "X"}, 3)
	require.Error(t, err)
	assert.Nil(t, fake.updated.Type)
	assert.Empty(t, fake.patched)
}

func TestApplyWatch_NoEpisodesStillUpdatesCollection(t *testing.T) {
	fake := &fakeWatch{getErr: &api.APIError{StatusCode: 404}}
	result, err := applyWatch(context.Background(), fake, api.Subject{ID: 1, Name: "X"}, 2)
	require.NoError(t, err)
	assert.Equal(t, 2, result.EpStatus)
	assert.Contains(t, result.Note, "没有章节数据")
	assert.Empty(t, fake.patched)
}

func TestParsePositiveID(t *testing.T) {
	n, err := parsePositiveID("12", "条目ID")
	require.NoError(t, err)
	assert.Equal(t, 12, n)
	_, err = parsePositiveID("0", "条目ID")
	require.Error(t, err)
	_, err = parsePositiveID("abc", "条目ID")
	require.Error(t, err)
}

type fakeWatch struct {
	current      *api.UserSubjectCollection
	getErr       error
	episodes     []api.EpisodeDetail
	episodeCalls int
	updated      api.UserSubjectCollectionUpdate
	patched      []int
}

func (f *fakeWatch) GetUserSubjectCollection(context.Context, string, int) (*api.UserSubjectCollection, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	if f.current == nil {
		return nil, &api.APIError{StatusCode: 404, Message: "not found"}
	}
	return f.current, nil
}

func (f *fakeWatch) UpdateUserSubjectCollection(_ context.Context, _ int, req api.UserSubjectCollectionUpdate) error {
	f.updated = req
	return nil
}

func (f *fakeWatch) GetEpisodes(_ context.Context, _ int, _ *api.EpType, limit, offset int) (*api.Paged[api.EpisodeDetail], error) {
	f.episodeCalls++
	if offset >= len(f.episodes) {
		return &api.Paged[api.EpisodeDetail]{Total: len(f.episodes)}, nil
	}
	end := offset + limit
	if limit <= 0 || end > len(f.episodes) {
		end = len(f.episodes)
	}
	return &api.Paged[api.EpisodeDetail]{
		Data:  append([]api.EpisodeDetail(nil), f.episodes[offset:end]...),
		Total: len(f.episodes),
	}, nil
}

func (f *fakeWatch) PatchUserEpisodeCollections(_ context.Context, _ int, episodeIDs []int, _ api.EpisodeCollectionType) error {
	f.patched = append(f.patched, episodeIDs...)
	return nil
}

func mainEpisodes(from, to int) []api.EpisodeDetail {
	var eps []api.EpisodeDetail
	for n := from; n <= to; n++ {
		eps = append(eps, api.EpisodeDetail{ID: 100 + n, Ep: float64(n)})
	}
	return eps
}
