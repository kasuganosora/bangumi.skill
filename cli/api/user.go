package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ===========================================================================
// User
// ===========================================================================

// GetUserByName 获取用户信息
func (c *HTTPClient) GetUserByName(ctx context.Context, username string) (*UserDetail, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v0/users/%s", username), nil)
	if err != nil {
		return nil, err
	}
	var result UserDetail
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("get user %s: %w", username, err)
	}
	return &result, nil
}

// GetUserAvatar 获取用户头像
func (c *HTTPClient) GetUserAvatar(ctx context.Context, username string, imgType string) (string, error) {
	return c.getImageRedirect(ctx, fmt.Sprintf("/v0/users/%s/avatar", username), imgType)
}

// GetMe 获取当前用户信息（需 Token）
func (c *HTTPClient) GetMe(ctx context.Context) (*UserDetail, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/v0/me", nil)
	if err != nil {
		return nil, err
	}
	var result UserDetail
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("get me: %w", err)
	}
	return &result, nil
}

// ===========================================================================
// User Collections
// ===========================================================================

// GetUserCollections 获取用户收藏列表
func (c *HTTPClient) GetUserCollections(ctx context.Context, username string, subjectType *SubjectType, collectionType *SubjectCollectionType, limit, offset int) (*Paged[UserSubjectCollection], error) {
	path := fmt.Sprintf("/v0/users/%s/collections", username)
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	if subjectType != nil {
		q.Set("subject_type", fmt.Sprintf("%d", *subjectType))
	}
	if collectionType != nil {
		q.Set("type", fmt.Sprintf("%d", *collectionType))
	}
	if limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", limit))
	}
	if offset > 0 {
		q.Set("offset", fmt.Sprintf("%d", offset))
	}
	req.URL.RawQuery = q.Encode()

	var result Paged[UserSubjectCollection]
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("get user %s collections: %w", username, err)
	}
	return &result, nil
}

// GetUserSubjectCollection 获取用户对指定条目的收藏
func (c *HTTPClient) GetUserSubjectCollection(ctx context.Context, username string, subjectID int) (*UserSubjectCollection, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v0/users/%s/collections/%d", username, subjectID), nil)
	if err != nil {
		return nil, err
	}
	var result UserSubjectCollection
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("get user %s collection %d: %w", username, subjectID, err)
	}
	return &result, nil
}

// UpdateUserSubjectCollection 更新当前用户对条目的收藏状态
func (c *HTTPClient) UpdateUserSubjectCollection(ctx context.Context, subjectID int, r UserSubjectCollectionUpdate) error {
	data, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("marshal collection update: %w", err)
	}
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/v0/users/-/collections/%d", subjectID), strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.do(req, nil)
}

// GetUserSubjectEpisodeCollection 获取用户章节收藏
func (c *HTTPClient) GetUserSubjectEpisodeCollection(ctx context.Context, subjectID int, limit, offset int) (*Paged[UserEpisodeCollection], error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v0/users/-/collections/%d/episodes", subjectID), nil)
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

	var result Paged[UserEpisodeCollection]
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("get episode collections: %w", err)
	}
	return &result, nil
}

// UpdateUserEpisodeCollection 更新当前用户章节收藏状态
func (c *HTTPClient) UpdateUserEpisodeCollection(ctx context.Context, episodeID int, typ EpisodeCollectionType) error {
	body := map[string]int{"type": int(typ)}
	data, _ := json.Marshal(body)
	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/v0/users/-/collections/-/episodes/%d", episodeID), strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.do(req, nil)
}

// GetUserCharacterCollections 获取用户角色收藏
func (c *HTTPClient) GetUserCharacterCollections(ctx context.Context, username string) ([]UserCharacterCollection, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v0/users/%s/collections/-/characters", username), nil)
	if err != nil {
		return nil, err
	}
	var result Paged[UserCharacterCollection]
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("get character collections: %w", err)
	}
	return result.Data, nil
}

// UpdateUserCharacterCollection 更新当前用户角色收藏
func (c *HTTPClient) UpdateUserCharacterCollection(ctx context.Context, characterID int) error {
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/v0/users/-/collections/-/characters/%d", characterID), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// GetUserPersonCollections 获取用户人物收藏
func (c *HTTPClient) GetUserPersonCollections(ctx context.Context, username string) ([]UserPersonCollection, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v0/users/%s/collections/-/persons", username), nil)
	if err != nil {
		return nil, err
	}
	var result Paged[UserPersonCollection]
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("get person collections: %w", err)
	}
	return result.Data, nil
}

// UpdateUserPersonCollection 更新当前用户人物收藏
func (c *HTTPClient) UpdateUserPersonCollection(ctx context.Context, personID int) error {
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/v0/users/-/collections/-/persons/%d", personID), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}
