package cmd

import (
	"testing"
	"time"

	"github.com/kasuganosora/bangumi.skill/cli/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBangumiWeekdayID(t *testing.T) {
	monday := time.Date(2026, 9, 28, 8, 0, 0, 0, time.UTC)
	require.Equal(t, time.Monday, monday.Weekday())
	assert.Equal(t, 1, bangumiWeekdayID(monday))

	sunday := time.Date(2026, 9, 27, 8, 0, 0, 0, time.UTC)
	require.Equal(t, time.Sunday, sunday.Weekday())
	assert.Equal(t, 7, bangumiWeekdayID(sunday))
}

func TestFilterCalendar(t *testing.T) {
	items := []api.CalendarItem{
		{Weekday: api.Weekday{ID: 1, CN: "星期一"}, Items: []api.SubjectSmall{{ID: 1, Name: "A"}}},
		{Weekday: api.Weekday{ID: 3, CN: "星期三"}, Items: []api.SubjectSmall{{ID: 2, Name: "B"}}},
	}
	got := filterCalendar(items, 3)
	require.Len(t, got, 1)
	assert.Equal(t, "B", got[0].Items[0].Name)
	assert.Empty(t, filterCalendar(items, 7))
}
