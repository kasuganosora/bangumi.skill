package api

import (
	"context"
	"fmt"
	"net/http"
)

// ===========================================================================
// Revisions
// ===========================================================================

func revisionPath(entity string, id int) string {
	if id > 0 {
		return fmt.Sprintf("/v0/revisions/%s/%d", entity, id)
	}
	return fmt.Sprintf("/v0/revisions/%s", entity)
}

// GetSubjectRevisions 获取条目编辑历史
func (c *HTTPClient) GetSubjectRevisions(ctx context.Context, limit, offset int) (*Paged[Revision], error) {
	return c.getRevisions(ctx, "subjects", limit, offset)
}

// GetSubjectRevisionByID 获取条目编辑历史详情
func (c *HTTPClient) GetSubjectRevisionByID(ctx context.Context, id int) (*DetailedRevision, error) {
	return c.getRevisionByID(ctx, "subjects", id)
}

// GetCharacterRevisions 获取角色编辑历史
func (c *HTTPClient) GetCharacterRevisions(ctx context.Context, limit, offset int) (*Paged[Revision], error) {
	return c.getRevisions(ctx, "characters", limit, offset)
}

// GetCharacterRevisionByID 获取角色编辑历史详情
func (c *HTTPClient) GetCharacterRevisionByID(ctx context.Context, id int) (*DetailedRevision, error) {
	return c.getRevisionByID(ctx, "characters", id)
}

// GetPersonRevisions 获取人物编辑历史
func (c *HTTPClient) GetPersonRevisions(ctx context.Context, limit, offset int) (*Paged[Revision], error) {
	return c.getRevisions(ctx, "persons", limit, offset)
}

// GetPersonRevisionByID 获取人物编辑历史详情
func (c *HTTPClient) GetPersonRevisionByID(ctx context.Context, id int) (*DetailedRevision, error) {
	return c.getRevisionByID(ctx, "persons", id)
}

// GetEpisodeRevisions 获取章节编辑历史
func (c *HTTPClient) GetEpisodeRevisions(ctx context.Context, limit, offset int) (*Paged[Revision], error) {
	return c.getRevisions(ctx, "episodes", limit, offset)
}

// GetEpisodeRevisionByID 获取章节编辑历史详情
func (c *HTTPClient) GetEpisodeRevisionByID(ctx context.Context, id int) (*DetailedRevision, error) {
	return c.getRevisionByID(ctx, "episodes", id)
}

func (c *HTTPClient) getRevisions(ctx context.Context, entity string, limit, offset int) (*Paged[Revision], error) {
	req, err := c.newRequest(ctx, http.MethodGet, revisionPath(entity, 0), nil)
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

	var result Paged[Revision]
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("get %s revisions: %w", entity, err)
	}
	return &result, nil
}

func (c *HTTPClient) getRevisionByID(ctx context.Context, entity string, id int) (*DetailedRevision, error) {
	req, err := c.newRequest(ctx, http.MethodGet, revisionPath(entity, id), nil)
	if err != nil {
		return nil, err
	}
	var result DetailedRevision
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("get %s revision %d: %w", entity, id, err)
	}
	return &result, nil
}
