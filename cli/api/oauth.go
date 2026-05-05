package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// ---------------------------------------------------------------------------
// OAuth 配置
// ---------------------------------------------------------------------------

// OAuthConfig 存储 OAuth App 注册信息。
// client_id 和 client_secret 在 https://bgm.tv/dev/app 注册获得。
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// ---------------------------------------------------------------------------
// Token 数据结构
// ---------------------------------------------------------------------------

// Token OAuth 授权返回的 Token
type Token struct {
	AccessToken  string      `json:"access_token"`
	ExpiresIn    int         `json:"expires_in"`
	TokenType    string      `json:"token_type"`
	Scope        interface{} `json:"scope"`
	RefreshToken string      `json:"refresh_token"`
	UserID       int         `json:"user_id"`
}

// TokenStatus Access Token 状态
type TokenStatus struct {
	AccessToken string      `json:"access_token"`
	ClientID    string      `json:"client_id"`
	Expires     int64       `json:"expires"`
	Scope       interface{} `json:"scope"`
	UserID      int         `json:"user_id"`
}

// ---------------------------------------------------------------------------
// Step 1: 生成授权 URL
// ---------------------------------------------------------------------------

// AuthorizeURL 生成引导用户授权的 URL。
// state 为随机字符串，用于防止 CSRF 攻击。
func (c *HTTPClient) AuthorizeURL(cfg OAuthConfig, state string) string {
	u, _ := url.Parse(c.oauthBaseURL.String() + "/oauth/authorize")
	q := u.Query()
	q.Set("client_id", cfg.ClientID)
	q.Set("response_type", "code")
	if cfg.RedirectURI != "" {
		q.Set("redirect_uri", cfg.RedirectURI)
	}
	if state != "" {
		q.Set("state", state)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// GenerateState 生成一个随机 state 字符串
func GenerateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate state: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// ---------------------------------------------------------------------------
// Step 3: 使用 code 换取 Access Token
// ---------------------------------------------------------------------------

// ExchangeCode 使用授权回调返回的 code 换取 Access Token。
// code 有效期为 60 秒。
func (c *HTTPClient) ExchangeCode(ctx context.Context, cfg OAuthConfig, code, state string) (*Token, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", cfg.ClientID)
	form.Set("client_secret", cfg.ClientSecret)
	form.Set("code", code)
	form.Set("redirect_uri", cfg.RedirectURI)
	if state != "" {
		form.Set("state", state)
	}

	req, err := c.newOAuthRequest(ctx, http.MethodPost, "/oauth/access_token",
		strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("exchange code: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var token Token
	if err := c.do(req, &token); err != nil {
		return nil, fmt.Errorf("exchange code: %w", err)
	}

	return &token, nil
}

// ---------------------------------------------------------------------------
// 刷新 Access Token
// ---------------------------------------------------------------------------

// RefreshAccessToken 使用 Refresh Token 刷新 Access Token。
func (c *HTTPClient) RefreshAccessToken(ctx context.Context, cfg OAuthConfig, refreshToken string) (*Token, error) {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("client_id", cfg.ClientID)
	form.Set("client_secret", cfg.ClientSecret)
	form.Set("refresh_token", refreshToken)
	form.Set("redirect_uri", cfg.RedirectURI)

	req, err := c.newOAuthRequest(ctx, http.MethodPost, "/oauth/access_token",
		strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("refresh token: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var token Token
	if err := c.do(req, &token); err != nil {
		return nil, fmt.Errorf("refresh token: %w", err)
	}

	return &token, nil
}

// ---------------------------------------------------------------------------
// 查询 Token 状态
// ---------------------------------------------------------------------------

// GetTokenStatus 查询 Access Token 的状态信息。
func (c *HTTPClient) GetTokenStatus(ctx context.Context, accessToken string) (*TokenStatus, error) {
	form := url.Values{}
	form.Set("access_token", accessToken)

	req, err := c.newOAuthRequest(ctx, http.MethodPost, "/oauth/token_status",
		strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("token status: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var status TokenStatus
	if err := c.do(req, &status); err != nil {
		return nil, fmt.Errorf("token status: %w", err)
	}

	return &status, nil
}

// ---------------------------------------------------------------------------
// JSON 序列化辅助
// ---------------------------------------------------------------------------

// UnmarshalJSON 处理 Token.scope 可能为 null 或字符串的情况
func (t *Token) UnmarshalJSON(data []byte) error {
	type alias Token
	aux := &struct {
		*alias
	}{
		alias: (*alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}
