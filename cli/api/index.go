package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ===========================================================================
// Index
// ===========================================================================

// NewIndex 创建目录
func (c *HTTPClient) NewIndex(ctx context.Context, r NewIndexRequest) (*Index, error) {
	data, _ := json.Marshal(r)
	req, err := c.newRequest(ctx, http.MethodPost, "/v0/indices", strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var result Index
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("new index: %w", err)
	}
	return &result, nil
}

// GetIndexByID 获取目录详情
func (c *HTTPClient) GetIndexByID(ctx context.Context, id int) (*Index, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v0/indices/%d", id), nil)
	if err != nil {
		return nil, err
	}
	var result Index
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("get index %d: %w", id, err)
	}
	return &result, nil
}

// EditIndexByID 编辑目录
func (c *HTTPClient) EditIndexByID(ctx context.Context, id int, r UpdateIndexRequest) (*Index, error) {
	data, _ := json.Marshal(r)
	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/v0/indices/%d", id), strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var result Index
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("edit index %d: %w", id, err)
	}
	return &result, nil
}

// GetIndexSubjects 获取目录条目列表
func (c *HTTPClient) GetIndexSubjects(ctx context.Context, id int, limit, offset int) (*Paged[IndexSubject], error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v0/indices/%d/subjects", id), nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	if limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", limit))
	}
	if offset > 0 {
		q.Set("offset", fmt.Sprintf("%d", offset))
	}
	req.URL.RawQuery = q.Encode()

	var result Paged[IndexSubject]
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("get index %d subjects: %w", id, err)
	}
	return &result, nil
}

// AddIndexSubject 向目录添加条目
func (c *HTTPClient) AddIndexSubject(ctx context.Context, indexID int, r AddIndexSubjectRequest) error {
	data, _ := json.Marshal(r)
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/v0/indices/%d/subjects", indexID), strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.do(req, nil)
}

// EditIndexSubject 编辑目录中的条目
func (c *HTTPClient) EditIndexSubject(ctx context.Context, indexID, subjectID int, r EditIndexSubjectRequest) error {
	data, _ := json.Marshal(r)
	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/v0/indices/%d/subjects/%d", indexID, subjectID), strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.do(req, nil)
}

// DeleteIndexSubject 从目录删除条目
func (c *HTTPClient) DeleteIndexSubject(ctx context.Context, indexID, subjectID int) error {
	req, err := c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/v0/indices/%d/subjects/%d", indexID, subjectID), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// CollectIndex 收藏目录
func (c *HTTPClient) CollectIndex(ctx context.Context, id int) error {
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/v0/indices/%d/collect", id), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// UncollectIndex 取消收藏目录
func (c *HTTPClient) UncollectIndex(ctx context.Context, id int) error {
	req, err := c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/v0/indices/%d/collect", id), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}
