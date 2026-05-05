package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// ---------------------------------------------------------------------------
// Client 接口 —— 所有 API 方法的契约，便于 Mock 测试
// ---------------------------------------------------------------------------

// Client 定义 Bangumi API 客户端接口。
// 上层代码依赖此接口而非具体实现，实现高可测试性。
type Client interface {
	// ---- Legacy ----
	GetCalendar(ctx context.Context) ([]CalendarItem, error)
	SearchSubjectByKeywords(ctx context.Context, keywords string, typ SubjectType, responseGroup string, start, max int) (*Paged[SubjectSmall], error)

	// ---- Search ----
	SearchSubjects(ctx context.Context, req SearchSubjectRequest, limit, offset int) (*Paged[Subject], error)
	SearchCharacters(ctx context.Context, req SearchCharacterRequest, limit, offset int) (*Paged[CharacterDetail], error)
	SearchPersons(ctx context.Context, req SearchPersonRequest, limit, offset int) (*Paged[PersonDetail], error)

	// ---- Subject ----
	GetSubjects(ctx context.Context, typ SubjectType, cat SubjectCategory, sort string, year, month, limit, offset int) (*Paged[Subject], error)
	GetSubjectByID(ctx context.Context, id int) (*Subject, error)
	GetSubjectImage(ctx context.Context, id int, imgType string) (string, error)
	GetSubjectPersons(ctx context.Context, id int) ([]RelatedPerson, error)
	GetSubjectCharacters(ctx context.Context, id int) ([]RelatedCharacter, error)
	GetSubjectRelations(ctx context.Context, id int) ([]V0SubjectRelation, error)

	// ---- Episode ----
	GetEpisodes(ctx context.Context, subjectID int, typ *EpType, limit, offset int) (*Paged[EpisodeDetail], error)
	GetEpisodeByID(ctx context.Context, id int) (*EpisodeDetail, error)

	// ---- Character ----
	GetCharacterByID(ctx context.Context, id int) (*CharacterDetail, error)
	GetCharacterImage(ctx context.Context, id int, imgType string) (string, error)
	GetCharacterSubjects(ctx context.Context, id int) ([]V0RelatedSubject, error)
	GetCharacterPersons(ctx context.Context, id int) ([]CharacterPerson, error)
	CollectCharacter(ctx context.Context, id int) error
	UncollectCharacter(ctx context.Context, id int) error

	// ---- Person ----
	GetPersonByID(ctx context.Context, id int) (*PersonDetail, error)
	GetPersonImage(ctx context.Context, id int, imgType string) (string, error)
	GetPersonSubjects(ctx context.Context, id int) ([]V0RelatedSubject, error)
	GetPersonCharacters(ctx context.Context, id int) ([]CharacterPerson, error)
	CollectPerson(ctx context.Context, id int) error
	UncollectPerson(ctx context.Context, id int) error

	// ---- User ----
	GetUserByName(ctx context.Context, username string) (*UserDetail, error)
	GetUserAvatar(ctx context.Context, username string, imgType string) (string, error)
	GetMe(ctx context.Context) (*UserDetail, error)

	// ---- User Collections ----
	GetUserCollections(ctx context.Context, username string, subjectType *SubjectType, collectionType *SubjectCollectionType, limit, offset int) (*Paged[UserSubjectCollection], error)
	GetUserSubjectCollection(ctx context.Context, username string, subjectID int) (*UserSubjectCollection, error)
	UpdateUserSubjectCollection(ctx context.Context, subjectID int, req UserSubjectCollectionUpdate) error
	GetUserSubjectEpisodeCollection(ctx context.Context, subjectID int, limit, offset int) (*Paged[UserEpisodeCollection], error)
	UpdateUserEpisodeCollection(ctx context.Context, episodeID int, typ EpisodeCollectionType) error
	GetUserCharacterCollections(ctx context.Context, username string) ([]UserCharacterCollection, error)
	UpdateUserCharacterCollection(ctx context.Context, characterID int) error
	GetUserPersonCollections(ctx context.Context, username string) ([]UserPersonCollection, error)
	UpdateUserPersonCollection(ctx context.Context, personID int) error

	// ---- Revisions ----
	GetSubjectRevisions(ctx context.Context, limit, offset int) (*Paged[Revision], error)
	GetSubjectRevisionByID(ctx context.Context, id int) (*DetailedRevision, error)
	GetCharacterRevisions(ctx context.Context, limit, offset int) (*Paged[Revision], error)
	GetCharacterRevisionByID(ctx context.Context, id int) (*DetailedRevision, error)
	GetPersonRevisions(ctx context.Context, limit, offset int) (*Paged[Revision], error)
	GetPersonRevisionByID(ctx context.Context, id int) (*DetailedRevision, error)
	GetEpisodeRevisions(ctx context.Context, limit, offset int) (*Paged[Revision], error)
	GetEpisodeRevisionByID(ctx context.Context, id int) (*DetailedRevision, error)

	// ---- Index ----
	NewIndex(ctx context.Context, req NewIndexRequest) (*Index, error)
	GetIndexByID(ctx context.Context, id int) (*Index, error)
	EditIndexByID(ctx context.Context, id int, req UpdateIndexRequest) (*Index, error)
	GetIndexSubjects(ctx context.Context, id int, limit, offset int) (*Paged[IndexSubject], error)
	AddIndexSubject(ctx context.Context, indexID int, req AddIndexSubjectRequest) error
	EditIndexSubject(ctx context.Context, indexID, subjectID int, req EditIndexSubjectRequest) error
	DeleteIndexSubject(ctx context.Context, indexID, subjectID int) error
	CollectIndex(ctx context.Context, id int) error
	UncollectIndex(ctx context.Context, id int) error
}

// ---------------------------------------------------------------------------
// HTTP 客户端实现
// ---------------------------------------------------------------------------

const (
	defaultBaseURL = "https://api.bgm.tv"
	oauthBaseURL   = "https://bgm.tv"
	defaultTimeout = 30 * time.Second

	// UA 格式遵循 Bangumi API 规范：
	// {开发者ID}/{应用名}[/{版本}] ({项目主页})
	uaDeveloperID = "kasuganosora"
	uaAppName     = "bangumi-cli"
	uaProjectURL  = "https://github.com/kasuganosora/bangumi.skill"
)

// version 通过 ldflags 在构建时注入
var version = "dev"

// BuildUserAgent 构建符合 Bangumi API 规范的 User-Agent。
//
// 规范要求：非浏览器 API 使用者须指定开发者 ID + 应用名称；
// 可分发的应用须附带版本号；开源项目须附上项目主页。
func BuildUserAgent(v string) string {
	if v == "" {
		v = version
	}
	return fmt.Sprintf("%s/%s/%s (%s)", uaDeveloperID, uaAppName, v, uaProjectURL)
}

// HTTPClient 是 Client 接口的 HTTP 实现。
type HTTPClient struct {
	baseURL        *url.URL
	oauthBaseURL   *url.URL
	httpClient     *http.Client
	userAgent      string
	accessToken    string
	onUnauthorized func() // 401 回调：清除令牌并提示重新登录
}

// ClientOption 函数式选项模式
type ClientOption func(*HTTPClient)

// WithBaseURL 设置自定义 Base URL（便于测试时指向 mock server）
func WithBaseURL(rawURL string) ClientOption {
	return func(c *HTTPClient) {
		if u, err := url.Parse(rawURL); err == nil {
			c.baseURL = u
		}
	}
}

// WithHTTPClient 注入自定义 *http.Client（便于测试时使用 httptest server 的 client）
func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *HTTPClient) {
		c.httpClient = hc
	}
}

// WithUserAgent 设置 User-Agent
func WithUserAgent(ua string) ClientOption {
	return func(c *HTTPClient) {
		c.userAgent = ua
	}
}

// WithAccessToken 设置 Bearer Token
func WithAccessToken(token string) ClientOption {
	return func(c *HTTPClient) {
		c.accessToken = token
	}
}

// WithOnUnauthorized 设置 401 回调（清除本地令牌并引导重新登录）
func WithOnUnauthorized(fn func()) ClientOption {
	return func(c *HTTPClient) {
		c.onUnauthorized = fn
	}
}

// WithTimeout 设置请求超时（秒）
func WithTimeout(seconds int) ClientOption {
	return func(c *HTTPClient) {
		if seconds > 0 {
			c.httpClient.Timeout = time.Duration(seconds) * time.Second
		}
	}
}

// WithProxy 设置 HTTP 代理
func WithProxy(proxyURL string) ClientOption {
	return func(c *HTTPClient) {
		if proxyURL == "" {
			return
		}
		if u, err := url.Parse(proxyURL); err == nil {
			c.httpClient.Transport = &http.Transport{
				Proxy: http.ProxyURL(u),
			}
		}
	}
}

// WithOAuthBaseURL 设置 OAuth 服务器地址（便于测试）
func WithOAuthBaseURL(rawURL string) ClientOption {
	return func(c *HTTPClient) {
		if u, err := url.Parse(rawURL); err == nil {
			c.oauthBaseURL = u
		}
	}
}

// NewClient 创建 HTTPClient，默认使用生产环境 Base URL
func NewClient(opts ...ClientOption) (*HTTPClient, error) {
	base, err := url.Parse(defaultBaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse base url: %w", err)
	}
	oauthBase, err := url.Parse(oauthBaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse oauth base url: %w", err)
	}

	c := &HTTPClient{
		baseURL:      base,
		oauthBaseURL: oauthBase,
		httpClient:   &http.Client{Timeout: defaultTimeout},
		userAgent:    BuildUserAgent(""),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// BaseURL 暴露 baseURL 供测试断言
func (c *HTTPClient) BaseURL() string {
	return c.baseURL.String()
}

// newRequest 构建带认证与 User-Agent 的 API 请求（api.bgm.tv）
func (c *HTTPClient) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	u := c.baseURL.JoinPath(path)
	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	c.setHeaders(req)
	return req, nil
}

// newOAuthRequest 构建 OAuth 请求（bgm.tv，不自动带 Bearer token）
// OAuth 端点的认证信息通过 form body 传递，不可与 Header 同时使用。
func (c *HTTPClient) newOAuthRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	u := c.oauthBaseURL.JoinPath(path)
	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("create oauth request: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// setHeaders 设置公共请求头
func (c *HTTPClient) setHeaders(req *http.Request) {
	req.Header.Set("User-Agent", c.userAgent)
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "application/json")
	}
	if c.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}
}

// getImageRedirect 获取图片重定向 URL（不跟随重定向）
func (c *HTTPClient) getImageRedirect(ctx context.Context, path string, imgType string) (string, error) {
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return "", err
	}
	q := req.URL.Query()
	q.Set("type", imgType)
	req.URL.RawQuery = q.Encode()

	// 不自动跟随重定向
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: c.httpClient.Timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("get image: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == 302 {
		return resp.Header.Get("Location"), nil
	}
	return "", &APIError{StatusCode: resp.StatusCode, Message: "image not found"}
}

// do 执行请求并解析 JSON 响应
func (c *HTTPClient) do(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		if resp.StatusCode == http.StatusUnauthorized && c.onUnauthorized != nil {
			c.onUnauthorized()
		}
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

// ---------------------------------------------------------------------------
// 自定义错误
// ---------------------------------------------------------------------------

// APIError 表示 API 返回的错误
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("bangumi api error (status=%d): %s", e.StatusCode, e.Message)
}
