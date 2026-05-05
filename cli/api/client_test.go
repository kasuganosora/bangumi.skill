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
// Mock 实现 —— 用于上层代码的单元测试
// ---------------------------------------------------------------------------

// MockClient 实现 Client 接口，所有方法可自定义返回
type MockClient struct {
	GetCalendarFunc                     func(ctx context.Context) ([]CalendarItem, error)
	SearchSubjectByKeywordsFunc         func(ctx context.Context, keywords string, typ SubjectType, rg string, start, max int) (*Paged[SubjectSmall], error)
	SearchSubjectsFunc                  func(ctx context.Context, req SearchSubjectRequest, limit, offset int) (*Paged[Subject], error)
	SearchCharactersFunc                func(ctx context.Context, req SearchCharacterRequest, limit, offset int) (*Paged[CharacterDetail], error)
	SearchPersonsFunc                   func(ctx context.Context, req SearchPersonRequest, limit, offset int) (*Paged[PersonDetail], error)
	GetSubjectsFunc                     func(ctx context.Context, typ SubjectType, cat SubjectCategory, sort string, year, month, limit, offset int) (*Paged[Subject], error)
	GetSubjectByIDFunc                  func(ctx context.Context, id int) (*Subject, error)
	GetSubjectImageFunc                 func(ctx context.Context, id int, imgType string) (string, error)
	GetSubjectPersonsFunc               func(ctx context.Context, id int) ([]RelatedPerson, error)
	GetSubjectCharactersFunc            func(ctx context.Context, id int) ([]RelatedCharacter, error)
	GetSubjectRelationsFunc             func(ctx context.Context, id int) ([]V0SubjectRelation, error)
	GetEpisodesFunc                     func(ctx context.Context, subjectID int, typ *EpType, limit, offset int) (*Paged[EpisodeDetail], error)
	GetEpisodeByIDFunc                  func(ctx context.Context, id int) (*EpisodeDetail, error)
	GetCharacterByIDFunc                func(ctx context.Context, id int) (*CharacterDetail, error)
	GetCharacterImageFunc               func(ctx context.Context, id int, imgType string) (string, error)
	GetCharacterSubjectsFunc            func(ctx context.Context, id int) ([]V0RelatedSubject, error)
	GetCharacterPersonsFunc             func(ctx context.Context, id int) ([]CharacterPerson, error)
	CollectCharacterFunc                func(ctx context.Context, id int) error
	UncollectCharacterFunc              func(ctx context.Context, id int) error
	GetPersonByIDFunc                   func(ctx context.Context, id int) (*PersonDetail, error)
	GetPersonImageFunc                  func(ctx context.Context, id int, imgType string) (string, error)
	GetPersonSubjectsFunc               func(ctx context.Context, id int) ([]V0RelatedSubject, error)
	GetPersonCharactersFunc             func(ctx context.Context, id int) ([]CharacterPerson, error)
	CollectPersonFunc                   func(ctx context.Context, id int) error
	UncollectPersonFunc                 func(ctx context.Context, id int) error
	GetUserByNameFunc                   func(ctx context.Context, username string) (*UserDetail, error)
	GetUserAvatarFunc                   func(ctx context.Context, username string, imgType string) (string, error)
	GetMeFunc                           func(ctx context.Context) (*UserDetail, error)
	GetUserCollectionsFunc              func(ctx context.Context, username string, st *SubjectType, ct *SubjectCollectionType, limit, offset int) (*Paged[UserSubjectCollection], error)
	GetUserSubjectCollectionFunc        func(ctx context.Context, username string, subjectID int) (*UserSubjectCollection, error)
	UpdateUserSubjectCollectionFunc     func(ctx context.Context, subjectID int, req UserSubjectCollectionUpdate) error
	GetUserSubjectEpisodeCollectionFunc func(ctx context.Context, subjectID int, limit, offset int) (*Paged[UserEpisodeCollection], error)
	UpdateUserEpisodeCollectionFunc     func(ctx context.Context, episodeID int, typ EpisodeCollectionType) error
	GetUserCharacterCollectionsFunc     func(ctx context.Context, username string) ([]UserCharacterCollection, error)
	UpdateUserCharacterCollectionFunc   func(ctx context.Context, characterID int) error
	GetUserPersonCollectionsFunc        func(ctx context.Context, username string) ([]UserPersonCollection, error)
	UpdateUserPersonCollectionFunc      func(ctx context.Context, personID int) error
	GetSubjectRevisionsFunc             func(ctx context.Context, limit, offset int) (*Paged[Revision], error)
	GetSubjectRevisionByIDFunc          func(ctx context.Context, id int) (*DetailedRevision, error)
	GetCharacterRevisionsFunc           func(ctx context.Context, limit, offset int) (*Paged[Revision], error)
	GetCharacterRevisionByIDFunc        func(ctx context.Context, id int) (*DetailedRevision, error)
	GetPersonRevisionsFunc              func(ctx context.Context, limit, offset int) (*Paged[Revision], error)
	GetPersonRevisionByIDFunc           func(ctx context.Context, id int) (*DetailedRevision, error)
	GetEpisodeRevisionsFunc             func(ctx context.Context, limit, offset int) (*Paged[Revision], error)
	GetEpisodeRevisionByIDFunc          func(ctx context.Context, id int) (*DetailedRevision, error)
	NewIndexFunc                        func(ctx context.Context, req NewIndexRequest) (*Index, error)
	GetIndexByIDFunc                    func(ctx context.Context, id int) (*Index, error)
	EditIndexByIDFunc                   func(ctx context.Context, id int, req UpdateIndexRequest) (*Index, error)
	GetIndexSubjectsFunc                func(ctx context.Context, id int, limit, offset int) (*Paged[IndexSubject], error)
	AddIndexSubjectFunc                 func(ctx context.Context, indexID int, req AddIndexSubjectRequest) error
	EditIndexSubjectFunc                func(ctx context.Context, indexID, subjectID int, req EditIndexSubjectRequest) error
	DeleteIndexSubjectFunc              func(ctx context.Context, indexID, subjectID int) error
	CollectIndexFunc                    func(ctx context.Context, id int) error
	UncollectIndexFunc                  func(ctx context.Context, id int) error
}

// 各方法代理到对应 Func 字段

func (m *MockClient) GetCalendar(ctx context.Context) ([]CalendarItem, error) {
	if m.GetCalendarFunc != nil {
		return m.GetCalendarFunc(ctx)
	}
	return nil, fmt.Errorf("GetCalendar not mocked")
}
func (m *MockClient) SearchSubjectByKeywords(ctx context.Context, k string, t SubjectType, rg string, s, mx int) (*Paged[SubjectSmall], error) {
	if m.SearchSubjectByKeywordsFunc != nil {
		return m.SearchSubjectByKeywordsFunc(ctx, k, t, rg, s, mx)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) SearchSubjects(ctx context.Context, r SearchSubjectRequest, l, o int) (*Paged[Subject], error) {
	if m.SearchSubjectsFunc != nil {
		return m.SearchSubjectsFunc(ctx, r, l, o)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) SearchCharacters(ctx context.Context, r SearchCharacterRequest, l, o int) (*Paged[CharacterDetail], error) {
	if m.SearchCharactersFunc != nil {
		return m.SearchCharactersFunc(ctx, r, l, o)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) SearchPersons(ctx context.Context, r SearchPersonRequest, l, o int) (*Paged[PersonDetail], error) {
	if m.SearchPersonsFunc != nil {
		return m.SearchPersonsFunc(ctx, r, l, o)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) GetSubjects(ctx context.Context, t SubjectType, c SubjectCategory, s string, y, mo, l, o int) (*Paged[Subject], error) {
	if m.GetSubjectsFunc != nil {
		return m.GetSubjectsFunc(ctx, t, c, s, y, mo, l, o)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) GetSubjectByID(ctx context.Context, id int) (*Subject, error) {
	if m.GetSubjectByIDFunc != nil {
		return m.GetSubjectByIDFunc(ctx, id)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) GetSubjectImage(ctx context.Context, id int, it string) (string, error) {
	if m.GetSubjectImageFunc != nil {
		return m.GetSubjectImageFunc(ctx, id, it)
	}
	return "", fmt.Errorf("not mocked")
}
func (m *MockClient) GetSubjectPersons(ctx context.Context, id int) ([]RelatedPerson, error) {
	if m.GetSubjectPersonsFunc != nil {
		return m.GetSubjectPersonsFunc(ctx, id)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) GetSubjectCharacters(ctx context.Context, id int) ([]RelatedCharacter, error) {
	if m.GetSubjectCharactersFunc != nil {
		return m.GetSubjectCharactersFunc(ctx, id)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) GetSubjectRelations(ctx context.Context, id int) ([]V0SubjectRelation, error) {
	if m.GetSubjectRelationsFunc != nil {
		return m.GetSubjectRelationsFunc(ctx, id)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) GetEpisodes(ctx context.Context, sid int, t *EpType, l, o int) (*Paged[EpisodeDetail], error) {
	if m.GetEpisodesFunc != nil {
		return m.GetEpisodesFunc(ctx, sid, t, l, o)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) GetEpisodeByID(ctx context.Context, id int) (*EpisodeDetail, error) {
	if m.GetEpisodeByIDFunc != nil {
		return m.GetEpisodeByIDFunc(ctx, id)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) GetCharacterByID(ctx context.Context, id int) (*CharacterDetail, error) {
	if m.GetCharacterByIDFunc != nil {
		return m.GetCharacterByIDFunc(ctx, id)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) GetCharacterImage(ctx context.Context, id int, it string) (string, error) {
	if m.GetCharacterImageFunc != nil {
		return m.GetCharacterImageFunc(ctx, id, it)
	}
	return "", fmt.Errorf("not mocked")
}
func (m *MockClient) GetCharacterSubjects(ctx context.Context, id int) ([]V0RelatedSubject, error) {
	if m.GetCharacterSubjectsFunc != nil {
		return m.GetCharacterSubjectsFunc(ctx, id)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) GetCharacterPersons(ctx context.Context, id int) ([]CharacterPerson, error) {
	if m.GetCharacterPersonsFunc != nil {
		return m.GetCharacterPersonsFunc(ctx, id)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) CollectCharacter(ctx context.Context, id int) error {
	if m.CollectCharacterFunc != nil {
		return m.CollectCharacterFunc(ctx, id)
	}
	return fmt.Errorf("not mocked")
}
func (m *MockClient) UncollectCharacter(ctx context.Context, id int) error {
	if m.UncollectCharacterFunc != nil {
		return m.UncollectCharacterFunc(ctx, id)
	}
	return fmt.Errorf("not mocked")
}
func (m *MockClient) GetPersonByID(ctx context.Context, id int) (*PersonDetail, error) {
	if m.GetPersonByIDFunc != nil {
		return m.GetPersonByIDFunc(ctx, id)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) GetPersonImage(ctx context.Context, id int, it string) (string, error) {
	if m.GetPersonImageFunc != nil {
		return m.GetPersonImageFunc(ctx, id, it)
	}
	return "", fmt.Errorf("not mocked")
}
func (m *MockClient) GetPersonSubjects(ctx context.Context, id int) ([]V0RelatedSubject, error) {
	if m.GetPersonSubjectsFunc != nil {
		return m.GetPersonSubjectsFunc(ctx, id)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) GetPersonCharacters(ctx context.Context, id int) ([]CharacterPerson, error) {
	if m.GetPersonCharactersFunc != nil {
		return m.GetPersonCharactersFunc(ctx, id)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) CollectPerson(ctx context.Context, id int) error {
	if m.CollectPersonFunc != nil {
		return m.CollectPersonFunc(ctx, id)
	}
	return fmt.Errorf("not mocked")
}
func (m *MockClient) UncollectPerson(ctx context.Context, id int) error {
	if m.UncollectPersonFunc != nil {
		return m.UncollectPersonFunc(ctx, id)
	}
	return fmt.Errorf("not mocked")
}
func (m *MockClient) GetUserByName(ctx context.Context, u string) (*UserDetail, error) {
	if m.GetUserByNameFunc != nil {
		return m.GetUserByNameFunc(ctx, u)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) GetUserAvatar(ctx context.Context, u, it string) (string, error) {
	if m.GetUserAvatarFunc != nil {
		return m.GetUserAvatarFunc(ctx, u, it)
	}
	return "", fmt.Errorf("not mocked")
}
func (m *MockClient) GetMe(ctx context.Context) (*UserDetail, error) {
	if m.GetMeFunc != nil {
		return m.GetMeFunc(ctx)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) GetUserCollections(ctx context.Context, u string, st *SubjectType, ct *SubjectCollectionType, l, o int) (*Paged[UserSubjectCollection], error) {
	if m.GetUserCollectionsFunc != nil {
		return m.GetUserCollectionsFunc(ctx, u, st, ct, l, o)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) GetUserSubjectCollection(ctx context.Context, u string, sid int) (*UserSubjectCollection, error) {
	if m.GetUserSubjectCollectionFunc != nil {
		return m.GetUserSubjectCollectionFunc(ctx, u, sid)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) UpdateUserSubjectCollection(ctx context.Context, sid int, r UserSubjectCollectionUpdate) error {
	if m.UpdateUserSubjectCollectionFunc != nil {
		return m.UpdateUserSubjectCollectionFunc(ctx, sid, r)
	}
	return fmt.Errorf("not mocked")
}
func (m *MockClient) GetUserSubjectEpisodeCollection(ctx context.Context, sid, l, o int) (*Paged[UserEpisodeCollection], error) {
	if m.GetUserSubjectEpisodeCollectionFunc != nil {
		return m.GetUserSubjectEpisodeCollectionFunc(ctx, sid, l, o)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) UpdateUserEpisodeCollection(ctx context.Context, eid int, t EpisodeCollectionType) error {
	if m.UpdateUserEpisodeCollectionFunc != nil {
		return m.UpdateUserEpisodeCollectionFunc(ctx, eid, t)
	}
	return fmt.Errorf("not mocked")
}
func (m *MockClient) GetUserCharacterCollections(ctx context.Context, u string) ([]UserCharacterCollection, error) {
	if m.GetUserCharacterCollectionsFunc != nil {
		return m.GetUserCharacterCollectionsFunc(ctx, u)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) UpdateUserCharacterCollection(ctx context.Context, cid int) error {
	if m.UpdateUserCharacterCollectionFunc != nil {
		return m.UpdateUserCharacterCollectionFunc(ctx, cid)
	}
	return fmt.Errorf("not mocked")
}
func (m *MockClient) GetUserPersonCollections(ctx context.Context, u string) ([]UserPersonCollection, error) {
	if m.GetUserPersonCollectionsFunc != nil {
		return m.GetUserPersonCollectionsFunc(ctx, u)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) UpdateUserPersonCollection(ctx context.Context, pid int) error {
	if m.UpdateUserPersonCollectionFunc != nil {
		return m.UpdateUserPersonCollectionFunc(ctx, pid)
	}
	return fmt.Errorf("not mocked")
}
func (m *MockClient) GetSubjectRevisions(ctx context.Context, l, o int) (*Paged[Revision], error) {
	if m.GetSubjectRevisionsFunc != nil {
		return m.GetSubjectRevisionsFunc(ctx, l, o)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) GetSubjectRevisionByID(ctx context.Context, id int) (*DetailedRevision, error) {
	if m.GetSubjectRevisionByIDFunc != nil {
		return m.GetSubjectRevisionByIDFunc(ctx, id)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) GetCharacterRevisions(ctx context.Context, l, o int) (*Paged[Revision], error) {
	if m.GetCharacterRevisionsFunc != nil {
		return m.GetCharacterRevisionsFunc(ctx, l, o)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) GetCharacterRevisionByID(ctx context.Context, id int) (*DetailedRevision, error) {
	if m.GetCharacterRevisionByIDFunc != nil {
		return m.GetCharacterRevisionByIDFunc(ctx, id)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) GetPersonRevisions(ctx context.Context, l, o int) (*Paged[Revision], error) {
	if m.GetPersonRevisionsFunc != nil {
		return m.GetPersonRevisionsFunc(ctx, l, o)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) GetPersonRevisionByID(ctx context.Context, id int) (*DetailedRevision, error) {
	if m.GetPersonRevisionByIDFunc != nil {
		return m.GetPersonRevisionByIDFunc(ctx, id)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) GetEpisodeRevisions(ctx context.Context, l, o int) (*Paged[Revision], error) {
	if m.GetEpisodeRevisionsFunc != nil {
		return m.GetEpisodeRevisionsFunc(ctx, l, o)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) GetEpisodeRevisionByID(ctx context.Context, id int) (*DetailedRevision, error) {
	if m.GetEpisodeRevisionByIDFunc != nil {
		return m.GetEpisodeRevisionByIDFunc(ctx, id)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) NewIndex(ctx context.Context, r NewIndexRequest) (*Index, error) {
	if m.NewIndexFunc != nil {
		return m.NewIndexFunc(ctx, r)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) GetIndexByID(ctx context.Context, id int) (*Index, error) {
	if m.GetIndexByIDFunc != nil {
		return m.GetIndexByIDFunc(ctx, id)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) EditIndexByID(ctx context.Context, id int, r UpdateIndexRequest) (*Index, error) {
	if m.EditIndexByIDFunc != nil {
		return m.EditIndexByIDFunc(ctx, id, r)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) GetIndexSubjects(ctx context.Context, id, l, o int) (*Paged[IndexSubject], error) {
	if m.GetIndexSubjectsFunc != nil {
		return m.GetIndexSubjectsFunc(ctx, id, l, o)
	}
	return nil, fmt.Errorf("not mocked")
}
func (m *MockClient) AddIndexSubject(ctx context.Context, iid int, r AddIndexSubjectRequest) error {
	if m.AddIndexSubjectFunc != nil {
		return m.AddIndexSubjectFunc(ctx, iid, r)
	}
	return fmt.Errorf("not mocked")
}
func (m *MockClient) EditIndexSubject(ctx context.Context, iid, sid int, r EditIndexSubjectRequest) error {
	if m.EditIndexSubjectFunc != nil {
		return m.EditIndexSubjectFunc(ctx, iid, sid, r)
	}
	return fmt.Errorf("not mocked")
}
func (m *MockClient) DeleteIndexSubject(ctx context.Context, iid, sid int) error {
	if m.DeleteIndexSubjectFunc != nil {
		return m.DeleteIndexSubjectFunc(ctx, iid, sid)
	}
	return fmt.Errorf("not mocked")
}
func (m *MockClient) CollectIndex(ctx context.Context, id int) error {
	if m.CollectIndexFunc != nil {
		return m.CollectIndexFunc(ctx, id)
	}
	return fmt.Errorf("not mocked")
}
func (m *MockClient) UncollectIndex(ctx context.Context, id int) error {
	if m.UncollectIndexFunc != nil {
		return m.UncollectIndexFunc(ctx, id)
	}
	return fmt.Errorf("not mocked")
}

// 编译期验证 MockClient 实现了 Client 接口
var _ Client = (*MockClient)(nil)

// ---------------------------------------------------------------------------
// HTTP 客户端集成测试（使用 httptest 模拟真实 API）
// ---------------------------------------------------------------------------

type HTTPClientSuite struct {
	suite.Suite
	server *httptest.Server
	client *HTTPClient
}

func (s *HTTPClientSuite) SetupSuite() {
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/calendar":
			// 支持通过 query 参数模拟错误
			if r.URL.Query().Get("error") == "401" {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error":"unauthorized"}`))
				return
			}
			s.handleCalendar(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func (s *HTTPClientSuite) TearDownSuite() {
	s.server.Close()
}

func (s *HTTPClientSuite) SetupTest() {
	var err error
	s.client, err = NewClient(
		WithBaseURL(s.server.URL),
	)
	s.Require().NoError(err)
}

func (s *HTTPClientSuite) handleCalendar(w http.ResponseWriter, r *http.Request) {
	// 验证请求头
	if ua := r.Header.Get("User-Agent"); ua == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := []CalendarItem{
		{
			Weekday: Weekday{EN: "Mon", CN: "星期一", JA: "月耀日", ID: 1},
			Items: []SubjectSmall{
				{
					ID:   12,
					Name: "ちょびっツ",
					Type: SubjectAnime,
				},
			},
		},
		{
			Weekday: Weekday{EN: "Tue", CN: "星期二", JA: "火耀日", ID: 2},
			Items:   []SubjectSmall{},
		},
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func TestHTTPClientSuite(t *testing.T) {
	suite.Run(t, new(HTTPClientSuite))
}

// ---------------------------------------------------------------------------
// 测试用例
// ---------------------------------------------------------------------------

func (s *HTTPClientSuite) TestGetCalendar_Success() {
	items, err := s.client.GetCalendar(context.Background())

	s.Require().NoError(err)
	s.Require().Len(items, 2)

	// 验证第一条
	s.Assert().Equal("Mon", items[0].Weekday.EN)
	s.Assert().Equal("星期一", items[0].Weekday.CN)
	s.Assert().Equal(1, items[0].Weekday.ID)

	s.Require().Len(items[0].Items, 1)
	s.Assert().Equal(12, items[0].Items[0].ID)
	s.Assert().Equal("ちょびっツ", items[0].Items[0].Name)
	s.Assert().Equal(SubjectAnime, items[0].Items[0].Type)

	// 验证第二条（空列表）
	s.Assert().Equal("Tue", items[1].Weekday.EN)
	s.Assert().Len(items[1].Items, 0)
}

func (s *HTTPClientSuite) TestGetCalendar_ErrorResponse() {
	// 创建一个始终返回 401 的独立 server
	errServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"unauthorized"}`))
	}))
	defer errServer.Close()

	client, err := NewClient(WithBaseURL(errServer.URL))
	s.Require().NoError(err)

	_, err = client.GetCalendar(context.Background())
	s.Require().Error(err)

	var apiErr *APIError
	s.Require().ErrorAs(err, &apiErr)
	s.Assert().Equal(http.StatusUnauthorized, apiErr.StatusCode)
}

func (s *HTTPClientSuite) TestGetCalendar_Timeout() {
	// 创建一个永远阻塞的 server，但能被 Close 打断
	blockCh := make(chan struct{})
	slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-blockCh // block until server closes
	}))
	defer func() {
		close(blockCh) // 释放 handler
		slowServer.Close()
	}()

	client, err := NewClient(
		WithBaseURL(slowServer.URL),
		WithHTTPClient(&http.Client{Timeout: 1}), // 1ns, almost instant timeout
	)
	s.Require().NoError(err)

	_, err = client.GetCalendar(context.Background())
	s.Require().Error(err)
}

func (s *HTTPClientSuite) TestNewClient_WithOptions() {
	client, err := NewClient(
		WithBaseURL("https://example.com"),
		WithUserAgent("test-agent/1.0"),
		WithAccessToken("test-token"),
	)
	s.Require().NoError(err)

	s.Assert().Equal("https://example.com", client.BaseURL())
	s.Assert().Equal("test-agent/1.0", client.userAgent)
	s.Assert().Equal("test-token", client.accessToken)
}

func (s *HTTPClientSuite) TestNewClient_DefaultValues() {
	client, err := NewClient()
	s.Require().NoError(err)

	s.Assert().Equal("https://api.bgm.tv", client.BaseURL())
	s.Assert().Equal(BuildUserAgent(""), client.userAgent)
	s.Assert().Empty(client.accessToken)
}

// ---------------------------------------------------------------------------
// Mock 测试 —— 验证上层调用方可以使用 MockClient
// ---------------------------------------------------------------------------

func TestMockClient_GetCalendar(t *testing.T) {
	expected := []CalendarItem{
		{Weekday: Weekday{EN: "Wed", CN: "星期三", JA: "水耀日", ID: 3}},
	}

	mock := &MockClient{
		GetCalendarFunc: func(ctx context.Context) ([]CalendarItem, error) {
			return expected, nil
		},
	}

	items, err := mock.GetCalendar(context.Background())

	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "Wed", items[0].Weekday.EN)
}

func TestMockClient_GetCalendar_Error(t *testing.T) {
	mock := &MockClient{
		GetCalendarFunc: func(ctx context.Context) ([]CalendarItem, error) {
			return nil, fmt.Errorf("network error")
		},
	}

	_, err := mock.GetCalendar(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "network error")
}

func TestMockClient_NotMocked(t *testing.T) {
	mock := &MockClient{}
	_, err := mock.GetCalendar(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not mocked")
}

// ---------------------------------------------------------------------------
// APIError 测试
// ---------------------------------------------------------------------------

func TestAPIError_Error(t *testing.T) {
	err := &APIError{StatusCode: 404, Message: "not found"}
	assert.Contains(t, err.Error(), "404")
	assert.Contains(t, err.Error(), "not found")
}

// ---------------------------------------------------------------------------
// MonoSource 自定义 JSON 解析测试
// ---------------------------------------------------------------------------

func TestMonoSource_Unmarshal_String(t *testing.T) {
	data := []byte(`"anidb.net"`)
	var ms MonoSource
	err := json.Unmarshal(data, &ms)
	require.NoError(t, err)
	assert.Equal(t, []string{"anidb.net"}, ms.Values)
}

func TestMonoSource_Unmarshal_Array(t *testing.T) {
	data := []byte(`["anidb.net", "wikipedia.org"]`)
	var ms MonoSource
	err := json.Unmarshal(data, &ms)
	require.NoError(t, err)
	assert.Equal(t, []string{"anidb.net", "wikipedia.org"}, ms.Values)
}

func TestMonoSource_Marshal_Single(t *testing.T) {
	ms := MonoSource{Values: []string{"anidb.net"}}
	data, err := json.Marshal(ms)
	require.NoError(t, err)
	assert.Equal(t, `"anidb.net"`, string(data))
}

func TestMonoSource_Marshal_Array(t *testing.T) {
	ms := MonoSource{Values: []string{"a", "b"}}
	data, err := json.Marshal(ms)
	require.NoError(t, err)
	assert.Equal(t, `["a","b"]`, string(data))
}

// ---------------------------------------------------------------------------
// 表驱动测试
// ---------------------------------------------------------------------------

func TestAPIError_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		err      *APIError
		contains string
	}{
		{"401", &APIError{401, "unauthorized"}, "401"},
		{"403", &APIError{403, "forbidden"}, "403"},
		{"500", &APIError{500, "server error"}, "server error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Contains(t, tt.err.Error(), tt.contains)
		})
	}
}
