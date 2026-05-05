package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// ---------------------------------------------------------------------------
// OAuth Suite
// ---------------------------------------------------------------------------

type OAuthSuite struct {
	suite.Suite
	server *httptest.Server
	client *HTTPClient
	cfg    OAuthConfig
}

func (s *OAuthSuite) SetupSuite() {
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/oauth/access_token" && r.Method == http.MethodPost:
			s.handleAccessToken(w, r)
		case r.URL.Path == "/oauth/token_status" && r.Method == http.MethodPost:
			s.handleTokenStatus(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func (s *OAuthSuite) TearDownSuite() {
	s.server.Close()
}

func (s *OAuthSuite) SetupTest() {
	var err error
	s.client, err = NewClient(
		WithOAuthBaseURL(s.server.URL),
	)
	s.Require().NoError(err)

	s.cfg = OAuthConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURI:  "https://example.com/callback",
	}
}

func (s *OAuthSuite) handleAccessToken(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	grantType := r.FormValue("grant_type")

	// 模拟无效 code
	if r.FormValue("code") == "invalid" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"invalid_grant","error_description":"Invalid code"}`))
		return
	}

	switch grantType {
	case "authorization_code":
		s.writeToken(w, "access-code-token", "refresh-code-token")
	case "refresh_token":
		s.writeToken(w, "access-refreshed-token", "refresh-new-token")
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (s *OAuthSuite) handleTokenStatus(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	status := TokenStatus{
		AccessToken: r.FormValue("access_token"),
		ClientID:    "test-client-id",
		Expires:     1728000000,
		UserID:      42,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}

func (s *OAuthSuite) writeToken(w http.ResponseWriter, access, refresh string) {
	token := Token{
		AccessToken:  access,
		ExpiresIn:    604800,
		TokenType:    "Bearer",
		Scope:        nil,
		RefreshToken: refresh,
		UserID:       42,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(token)
}

func TestOAuthSuite(t *testing.T) {
	suite.Run(t, new(OAuthSuite))
}

// ---------------------------------------------------------------------------
// AuthorizeURL
// ---------------------------------------------------------------------------

func (s *OAuthSuite) TestAuthorizeURL() {
	url := s.client.AuthorizeURL(s.cfg, "random-state")

	s.Assert().Contains(url, "/oauth/authorize")
	s.Assert().Contains(url, "client_id=test-client-id")
	s.Assert().Contains(url, "response_type=code")
	s.Assert().Contains(url, "redirect_uri=https%3A%2F%2Fexample.com%2Fcallback")
	s.Assert().Contains(url, "state=random-state")
}

func (s *OAuthSuite) TestAuthorizeURL_NoRedirectURI() {
	cfg := OAuthConfig{ClientID: "id-only"}
	url := s.client.AuthorizeURL(cfg, "")

	s.Assert().Contains(url, "client_id=id-only")
	s.Assert().NotContains(url, "redirect_uri")
	s.Assert().NotContains(url, "state")
}

// ---------------------------------------------------------------------------
// GenerateState
// ---------------------------------------------------------------------------

func TestGenerateState(t *testing.T) {
	state1, err := GenerateState()
	require.NoError(t, err)
	assert.Len(t, state1, 32) // 16 bytes → 32 hex chars

	state2, err := GenerateState()
	require.NoError(t, err)
	assert.NotEqual(t, state1, state2, "states should be random")
}

// ---------------------------------------------------------------------------
// ExchangeCode
// ---------------------------------------------------------------------------

func (s *OAuthSuite) TestExchangeCode_Success() {
	token, err := s.client.ExchangeCode(context.Background(), s.cfg, "valid-code", "test-state")

	s.Require().NoError(err)
	s.Assert().Equal("access-code-token", token.AccessToken)
	s.Assert().Equal(604800, token.ExpiresIn)
	s.Assert().Equal("Bearer", token.TokenType)
	s.Assert().Equal("refresh-code-token", token.RefreshToken)
	s.Assert().Equal(42, token.UserID)
}

func (s *OAuthSuite) TestExchangeCode_InvalidCode() {
	_, err := s.client.ExchangeCode(context.Background(), s.cfg, "invalid", "")

	s.Require().Error(err)
	var apiErr *APIError
	s.Require().ErrorAs(err, &apiErr)
	s.Assert().Equal(http.StatusBadRequest, apiErr.StatusCode)
}

// ---------------------------------------------------------------------------
// RefreshAccessToken
// ---------------------------------------------------------------------------

func (s *OAuthSuite) TestRefreshAccessToken_Success() {
	token, err := s.client.RefreshAccessToken(context.Background(), s.cfg, "old-refresh-token")

	s.Require().NoError(err)
	s.Assert().Equal("access-refreshed-token", token.AccessToken)
	s.Assert().Equal("refresh-new-token", token.RefreshToken)
}

// ---------------------------------------------------------------------------
// GetTokenStatus
// ---------------------------------------------------------------------------

func (s *OAuthSuite) TestGetTokenStatus_Success() {
	status, err := s.client.GetTokenStatus(context.Background(), "some-access-token")

	s.Require().NoError(err)
	s.Assert().Equal("some-access-token", status.AccessToken)
	s.Assert().Equal("test-client-id", status.ClientID)
	s.Assert().Equal(int64(1728000000), status.Expires)
	s.Assert().Equal(42, status.UserID)
}

// ---------------------------------------------------------------------------
// Token JSON 反序列化
// ---------------------------------------------------------------------------

func TestToken_UnmarshalJSON(t *testing.T) {
	raw := `{"access_token":"tok","expires_in":3600,"token_type":"Bearer","scope":null,"refresh_token":"ref","user_id":1}`

	var token Token
	err := json.Unmarshal([]byte(raw), &token)

	require.NoError(t, err)
	assert.Equal(t, "tok", token.AccessToken)
	assert.Equal(t, 3600, token.ExpiresIn)
	assert.Equal(t, "ref", token.RefreshToken)
	assert.Equal(t, 1, token.UserID)
}

// ---------------------------------------------------------------------------
// OAuth 错误处理
// ---------------------------------------------------------------------------

func (s *OAuthSuite) TestExchangeCode_ServerError() {
	errServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer errServer.Close()

	client, err := NewClient(WithOAuthBaseURL(errServer.URL))
	s.Require().NoError(err)

	_, err = client.ExchangeCode(context.Background(), s.cfg, "x", "")
	s.Require().Error(err)
}

// ---------------------------------------------------------------------------
// 表驱动测试
// ---------------------------------------------------------------------------

func TestAuthorizeURL_TableDriven(t *testing.T) {
	client, err := NewClient(WithOAuthBaseURL("https://bgm.tv"))
	require.NoError(t, err)

	tests := []struct {
		name   string
		cfg    OAuthConfig
		state  string
		checks []string // 期望 URL 中包含的内容
	}{
		{
			name:  "full config",
			cfg:   OAuthConfig{ClientID: "abc", ClientSecret: "sec", RedirectURI: "https://x.com/cb"},
			state: "s1",
			checks: []string{
				"client_id=abc",
				"response_type=code",
				"redirect_uri=https%3A%2F%2Fx.com%2Fcb",
				"state=s1",
			},
		},
		{
			name:  "minimal config",
			cfg:   OAuthConfig{ClientID: "min"},
			state: "",
			checks: []string{
				"client_id=min",
				"response_type=code",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := client.AuthorizeURL(tt.cfg, tt.state)
			for _, check := range tt.checks {
				assert.Contains(t, u, check)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Mock Client 接口扩展
// ---------------------------------------------------------------------------

// 验证 MockClient 仍满足 Client 接口
func TestMockClient_OAuth_Compatibility(t *testing.T) {
	mock := &MockClient{
		GetCalendarFunc: func(ctx context.Context) ([]CalendarItem, error) {
			return nil, fmt.Errorf("mock")
		},
	}
	_, err := mock.GetCalendar(context.Background())
	assert.Error(t, err)
}
