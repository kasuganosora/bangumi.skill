package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kasuganosora/bangumi.skill/cli/log"
)

// ===========================================================================
// Character
// ===========================================================================

// GetCharacterByID 获取角色详情
func (c *HTTPClient) GetCharacterByID(ctx context.Context, id int) (*CharacterDetail, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v0/characters/%d", id), nil)
	if err != nil {
		return nil, err
	}
	log.DebugContext(ctx, "get character", "id", id)
	var result CharacterDetail
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("get character %d: %w", id, err)
	}
	return &result, nil
}

// GetCharacterImage 获取角色图片
func (c *HTTPClient) GetCharacterImage(ctx context.Context, id int, imgType string) (string, error) {
	return c.getImageRedirect(ctx, fmt.Sprintf("/v0/characters/%d/image", id), imgType)
}

// GetCharacterSubjects 获取角色出演条目
func (c *HTTPClient) GetCharacterSubjects(ctx context.Context, id int) ([]V0RelatedSubject, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v0/characters/%d/subjects", id), nil)
	if err != nil {
		return nil, err
	}
	var result []V0RelatedSubject
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("get character %d subjects: %w", id, err)
	}
	return result, nil
}

// GetCharacterPersons 获取角色声优
func (c *HTTPClient) GetCharacterPersons(ctx context.Context, id int) ([]CharacterPerson, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v0/characters/%d/persons", id), nil)
	if err != nil {
		return nil, err
	}
	var result []CharacterPerson
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("get character %d persons: %w", id, err)
	}
	return result, nil
}

// CollectCharacter 收藏角色
func (c *HTTPClient) CollectCharacter(ctx context.Context, id int) error {
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/v0/characters/%d/collect", id), nil)
	if err != nil {
		return err
	}
	log.DebugContext(ctx, "collect character", "id", id)
	return c.do(req, nil)
}

// UncollectCharacter 取消收藏角色
func (c *HTTPClient) UncollectCharacter(ctx context.Context, id int) error {
	req, err := c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/v0/characters/%d/collect", id), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}
