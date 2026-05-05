package api

import (
	"context"
	"fmt"
	"net/http"
)

// ---------------------------------------------------------------------------
// Calendar 数据结构
// ---------------------------------------------------------------------------

// CalendarItem 每日放送条目
type CalendarItem struct {
	Weekday Weekday        `json:"weekday"`
	Items   []SubjectSmall `json:"items"`
}

// Weekday 星期信息
type Weekday struct {
	EN string `json:"en"` // Mon, Tue, ...
	CN string `json:"cn"` // 星期一, 星期二, ...
	JA string `json:"ja"` // 月耀日, 火耀日, ...
	ID int    `json:"id"` // 1-7
}

// ---------------------------------------------------------------------------
// Calendar 端点实现
// ---------------------------------------------------------------------------

// GetCalendar 获取每日放送信息
func (c *HTTPClient) GetCalendar(ctx context.Context) ([]CalendarItem, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/calendar", nil)
	if err != nil {
		return nil, fmt.Errorf("build calendar request: %w", err)
	}

	var items []CalendarItem
	if err := c.do(req, &items); err != nil {
		return nil, fmt.Errorf("get calendar: %w", err)
	}

	return items, nil
}
