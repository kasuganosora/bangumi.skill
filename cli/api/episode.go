package api

import (
	"context"
	"fmt"
	"net/http"
)

// ===========================================================================
// Episode
// ===========================================================================

// GetEpisodes 获取章节列表
func (c *HTTPClient) GetEpisodes(ctx context.Context, subjectID int, typ *EpType, limit, offset int) (*Paged[EpisodeDetail], error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/v0/episodes", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("subject_id", fmt.Sprintf("%d", subjectID))
	if typ != nil {
		q.Set("type", fmt.Sprintf("%d", *typ))
	}
	if limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", limit))
	}
	if offset > 0 {
		q.Set("offset", fmt.Sprintf("%d", offset))
	}
	req.URL.RawQuery = q.Encode()

	var result Paged[EpisodeDetail]
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("get episodes: %w", err)
	}
	return &result, nil
}

// GetEpisodeByID 获取章节详情
func (c *HTTPClient) GetEpisodeByID(ctx context.Context, id int) (*EpisodeDetail, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v0/episodes/%d", id), nil)
	if err != nil {
		return nil, err
	}
	var result EpisodeDetail
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("get episode %d: %w", id, err)
	}
	return &result, nil
}
