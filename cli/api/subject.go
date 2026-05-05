package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/kasuganosora/bangumi.skill/cli/log"
)

// ===========================================================================
// Legacy Search
// ===========================================================================

// SearchSubjectByKeywords 旧版关键词搜索
func (c *HTTPClient) SearchSubjectByKeywords(ctx context.Context, keywords string, typ SubjectType, responseGroup string, start, max int) (*Paged[SubjectSmall], error) {
	path := fmt.Sprintf("/search/subject/%s", keywords)
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	if typ != 0 {
		q.Set("type", fmt.Sprintf("%d", typ))
	}
	if responseGroup != "" {
		q.Set("responseGroup", responseGroup)
	}
	if start > 0 {
		q.Set("start", fmt.Sprintf("%d", start))
	}
	if max > 0 {
		q.Set("max_results", fmt.Sprintf("%d", max))
	}
	req.URL.RawQuery = q.Encode()

	var result struct {
		Results int            `json:"results"`
		List    []SubjectSmall `json:"list"`
	}
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("search subject by keywords: %w", err)
	}
	return &Paged[SubjectSmall]{Data: result.List, Total: result.Results}, nil
}

// ===========================================================================
// v0 Search
// ===========================================================================

// SearchSubjects 条目搜索
func (c *HTTPClient) SearchSubjects(ctx context.Context, r SearchSubjectRequest, limit, offset int) (*Paged[Subject], error) {
	return postSearch[Paged[Subject]](c, ctx, "/v0/search/subjects", r, limit, offset)
}

// SearchCharacters 角色搜索
func (c *HTTPClient) SearchCharacters(ctx context.Context, r SearchCharacterRequest, limit, offset int) (*Paged[CharacterDetail], error) {
	return postSearch[Paged[CharacterDetail]](c, ctx, "/v0/search/characters", r, limit, offset)
}

// SearchPersons 人物搜索
func (c *HTTPClient) SearchPersons(ctx context.Context, r SearchPersonRequest, limit, offset int) (*Paged[PersonDetail], error) {
	return postSearch[Paged[PersonDetail]](c, ctx, "/v0/search/persons", r, limit, offset)
}

func postSearch[T any](c *HTTPClient, ctx context.Context, path string, body interface{}, limit, offset int) (*T, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal search request: %w", err)
	}
	req, err := c.newRequest(ctx, http.MethodPost, path, strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	if limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", limit))
	}
	if offset > 0 {
		q.Set("offset", fmt.Sprintf("%d", offset))
	}
	req.URL.RawQuery = q.Encode()

	var result T
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}
	return &result, nil
}

// ===========================================================================
// Subject CRUD
// ===========================================================================

// GetSubjects 浏览条目
func (c *HTTPClient) GetSubjects(ctx context.Context, typ SubjectType, cat SubjectCategory, sort string, year, month, limit, offset int) (*Paged[Subject], error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/v0/subjects", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("type", fmt.Sprintf("%d", typ))
	if cat != 0 {
		q.Set("cat", fmt.Sprintf("%d", cat))
	}
	if sort != "" {
		q.Set("sort", sort)
	}
	if year > 0 {
		q.Set("year", fmt.Sprintf("%d", year))
	}
	if month > 0 {
		q.Set("month", fmt.Sprintf("%d", month))
	}
	if limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", limit))
	}
	if offset > 0 {
		q.Set("offset", fmt.Sprintf("%d", offset))
	}
	req.URL.RawQuery = q.Encode()

	log.DebugContext(ctx, "get subjects", "type", typ, "sort", sort)
	var result Paged[Subject]
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("get subjects: %w", err)
	}
	return &result, nil
}

// GetSubjectByID 获取条目详情
func (c *HTTPClient) GetSubjectByID(ctx context.Context, id int) (*Subject, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v0/subjects/%d", id), nil)
	if err != nil {
		return nil, err
	}
	var result Subject
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("get subject %d: %w", id, err)
	}
	return &result, nil
}

// GetSubjectImage 获取条目图片
func (c *HTTPClient) GetSubjectImage(ctx context.Context, id int, imgType string) (string, error) {
	return c.getImageRedirect(ctx, fmt.Sprintf("/v0/subjects/%d/image", id), imgType)
}

// GetSubjectPersons 获取条目关联人物
func (c *HTTPClient) GetSubjectPersons(ctx context.Context, id int) ([]RelatedPerson, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v0/subjects/%d/persons", id), nil)
	if err != nil {
		return nil, err
	}
	var result []RelatedPerson
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("get subject %d persons: %w", id, err)
	}
	return result, nil
}

// GetSubjectCharacters 获取条目关联角色
func (c *HTTPClient) GetSubjectCharacters(ctx context.Context, id int) ([]RelatedCharacter, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v0/subjects/%d/characters", id), nil)
	if err != nil {
		return nil, err
	}
	var result []RelatedCharacter
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("get subject %d characters: %w", id, err)
	}
	return result, nil
}

// GetSubjectRelations 获取条目关联
func (c *HTTPClient) GetSubjectRelations(ctx context.Context, id int) ([]V0SubjectRelation, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v0/subjects/%d/subjects", id), nil)
	if err != nil {
		return nil, err
	}
	var result []V0SubjectRelation
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("get subject %d relations: %w", id, err)
	}
	return result, nil
}
