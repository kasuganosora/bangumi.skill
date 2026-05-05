package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kasuganosora/bangumi.skill/cli/log"
)

// ===========================================================================
// Person
// ===========================================================================

// GetPersonByID 获取人物详情
func (c *HTTPClient) GetPersonByID(ctx context.Context, id int) (*PersonDetail, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v0/persons/%d", id), nil)
	if err != nil {
		return nil, err
	}
	log.DebugContext(ctx, "get person", "id", id)
	var result PersonDetail
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("get person %d: %w", id, err)
	}
	return &result, nil
}

// GetPersonImage 获取人物图片
func (c *HTTPClient) GetPersonImage(ctx context.Context, id int, imgType string) (string, error) {
	return c.getImageRedirect(ctx, fmt.Sprintf("/v0/persons/%d/image", id), imgType)
}

// GetPersonSubjects 获取人物参与作品
func (c *HTTPClient) GetPersonSubjects(ctx context.Context, id int) ([]V0RelatedSubject, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v0/persons/%d/subjects", id), nil)
	if err != nil {
		return nil, err
	}
	var result []V0RelatedSubject
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("get person %d subjects: %w", id, err)
	}
	return result, nil
}

// GetPersonCharacters 获取人物配音角色
func (c *HTTPClient) GetPersonCharacters(ctx context.Context, id int) ([]CharacterPerson, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v0/persons/%d/characters", id), nil)
	if err != nil {
		return nil, err
	}
	var result []CharacterPerson
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("get person %d characters: %w", id, err)
	}
	return result, nil
}

// CollectPerson 收藏人物
func (c *HTTPClient) CollectPerson(ctx context.Context, id int) error {
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/v0/persons/%d/collect", id), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// UncollectPerson 取消收藏人物
func (c *HTTPClient) UncollectPerson(ctx context.Context, id int) error {
	req, err := c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/v0/persons/%d/collect", id), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}
