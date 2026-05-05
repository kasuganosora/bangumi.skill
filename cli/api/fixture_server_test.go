package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

// ===========================================================================
// FixtureServerSuite — 使用真实 API 响应作为 mock 数据的集成测试
// ===========================================================================

type FixtureServerSuite struct {
	suite.Suite
	server *httptest.Server
	client *HTTPClient
}

var fixtureRoutes = map[string]string{
	"/calendar":                  "calendar.json",
	"/v0/search/subjects":        "search_subjects.json",
	"/v0/search/characters":      "search_characters.json",
	"/v0/search/persons":         "search_persons.json",
	"/v0/subjects":               "subjects_anime.json",
	"/v0/subjects/12":            "subject_12.json",
	"/v0/subjects/12/persons":    "subject_12_persons.json",
	"/v0/subjects/12/characters": "subject_12_characters.json",
	"/v0/subjects/12/subjects":   "subject_12_relations.json",
	"/v0/episodes":               "episodes_12.json",
	"/v0/episodes/1027":          "episode_1027.json",
	"/v0/characters/67":          "character_67.json",
	"/v0/characters/67/subjects": "character_67_subjects.json",
	"/v0/characters/67/persons":  "character_67_persons.json",
	"/v0/persons/1":              "person_1.json",
	"/v0/persons/1/subjects":     "person_1_subjects.json",
	"/v0/persons/1/characters":   "person_1_characters.json",
	"/v0/users/sai":              "user_sai.json",
	"/v0/me":                     "me.json",
	"/v0/users/sai/collections":  "user_sai_collections.json",
	"/v0/indices/1":              "index_1.json",
	"/v0/indices/1/subjects":     "index_1_subjects.json",
	"/search/subject/EVA":        "legacy_search.json",
}

func (s *FixtureServerSuite) SetupSuite() {
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// route matching: /v0/subjects/12 → /v0/subjects/12, /v0/characters/67 → /v0/characters/67
		// 同时也匹配带 query 的路径
		fixtureName, ok := fixtureRoutes[path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		s.serveFixture(w, fixtureName)
	}))
}

func (s *FixtureServerSuite) TearDownSuite() {
	s.server.Close()
}

func (s *FixtureServerSuite) SetupTest() {
	var err error
	s.client, err = NewClient(WithBaseURL(s.server.URL))
	s.Require().NoError(err)
}

func (s *FixtureServerSuite) serveFixture(w http.ResponseWriter, name string) {
	path := filepath.Join("testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"fixture not found: ` + name + `"}`))
		return
	}
	// 验证有效 JSON
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

func TestFixtureServerSuite(t *testing.T) {
	suite.Run(t, new(FixtureServerSuite))
}

// ===========================================================================
// 测试用例 — 验证所有 fixture 能正确反序列化
// ===========================================================================

func (s *FixtureServerSuite) TestGetCalendar() {
	items, err := s.client.GetCalendar(s.T().Context())
	s.Require().NoError(err)
	s.Require().NotEmpty(items)
	s.Assert().NotEmpty(items[0].Weekday.CN)
}

func (s *FixtureServerSuite) TestSearchSubjects() {
	r, err := s.client.SearchSubjects(s.T().Context(), SearchSubjectRequest{}, 5, 0)
	s.Require().NoError(err)
	s.Require().NotEmpty(r.Data)
}

func (s *FixtureServerSuite) TestSearchCharacters() {
	r, err := s.client.SearchCharacters(s.T().Context(), SearchCharacterRequest{}, 5, 0)
	s.Require().NoError(err)
	s.Require().NotEmpty(r.Data)
}

func (s *FixtureServerSuite) TestSearchPersons() {
	r, err := s.client.SearchPersons(s.T().Context(), SearchPersonRequest{}, 5, 0)
	s.Require().NoError(err)
	s.Require().NotEmpty(r.Data)
}

func (s *FixtureServerSuite) TestGetSubjects() {
	r, err := s.client.GetSubjects(s.T().Context(), SubjectAnime, 0, "", 0, 0, 5, 0)
	s.Require().NoError(err)
	s.Require().NotEmpty(r.Data)
}

func (s *FixtureServerSuite) TestGetSubjectByID() {
	sj, err := s.client.GetSubjectByID(s.T().Context(), 12)
	s.Require().NoError(err)
	s.Assert().Equal(12, sj.ID)
}

func (s *FixtureServerSuite) TestGetSubjectPersons() {
	r, err := s.client.GetSubjectPersons(s.T().Context(), 12)
	s.Require().NoError(err)
	s.Require().NotEmpty(r)
}

func (s *FixtureServerSuite) TestGetSubjectCharacters() {
	r, err := s.client.GetSubjectCharacters(s.T().Context(), 12)
	s.Require().NoError(err)
	s.Require().NotEmpty(r)
}

func (s *FixtureServerSuite) TestGetSubjectRelations() {
	r, err := s.client.GetSubjectRelations(s.T().Context(), 12)
	s.Require().NoError(err)
	s.Require().NotEmpty(r)
}

func (s *FixtureServerSuite) TestGetEpisodes() {
	r, err := s.client.GetEpisodes(s.T().Context(), 12, nil, 100, 0)
	s.Require().NoError(err)
	s.Require().NotEmpty(r.Data)
}

func (s *FixtureServerSuite) TestGetEpisodeByID() {
	e, err := s.client.GetEpisodeByID(s.T().Context(), 1027)
	s.Require().NoError(err)
	s.Assert().Equal(1027, e.ID)
}

func (s *FixtureServerSuite) TestGetCharacterByID() {
	c, err := s.client.GetCharacterByID(s.T().Context(), 67)
	s.Require().NoError(err)
	s.Assert().Equal(67, c.ID)
}

func (s *FixtureServerSuite) TestGetCharacterSubjects() {
	r, err := s.client.GetCharacterSubjects(s.T().Context(), 67)
	s.Require().NoError(err)
	s.Require().NotEmpty(r)
}

func (s *FixtureServerSuite) TestGetCharacterPersons() {
	r, err := s.client.GetCharacterPersons(s.T().Context(), 67)
	s.Require().NoError(err)
	s.Require().NotEmpty(r)
}

func (s *FixtureServerSuite) TestGetPersonByID() {
	p, err := s.client.GetPersonByID(s.T().Context(), 1)
	s.Require().NoError(err)
	s.Assert().Equal(1, p.ID)
}

func (s *FixtureServerSuite) TestGetPersonSubjects() {
	r, err := s.client.GetPersonSubjects(s.T().Context(), 1)
	s.Require().NoError(err)
	s.Require().NotEmpty(r)
}

func (s *FixtureServerSuite) TestGetPersonCharacters() {
	r, err := s.client.GetPersonCharacters(s.T().Context(), 1)
	s.Require().NoError(err)
	s.Require().NotEmpty(r)
}

func (s *FixtureServerSuite) TestGetUserByName() {
	u, err := s.client.GetUserByName(s.T().Context(), "sai")
	s.Require().NoError(err)
	s.Assert().Equal("test_user", u.Username)
	s.Assert().Equal(100000, u.ID)
}

func (s *FixtureServerSuite) TestGetMe() {
	u, err := s.client.GetMe(s.T().Context())
	s.Require().NoError(err)
	s.Assert().Equal("demo_user", u.Username)
	s.Assert().Equal(999999, u.ID)
}

func (s *FixtureServerSuite) TestGetUserCollections() {
	r, err := s.client.GetUserCollections(s.T().Context(), "sai", nil, nil, 5, 0)
	s.Require().NoError(err)
	s.Require().NotEmpty(r.Data)
}

func (s *FixtureServerSuite) TestGetIndexByID() {
	idx, err := s.client.GetIndexByID(s.T().Context(), 1)
	s.Require().NoError(err)
	s.Require().NotEmpty(idx.Title)
}

func (s *FixtureServerSuite) TestIndexSubjects() {
	r, err := s.client.GetIndexSubjects(s.T().Context(), 1, 5, 0)
	s.Require().NoError(err)
	s.Require().NotEmpty(r.Data)
}

func (s *FixtureServerSuite) TestLegacySearch() {
	r, err := s.client.SearchSubjectByKeywords(s.T().Context(), "EVA", 0, "small", 0, 5)
	s.Require().NoError(err)
	s.Require().NotEmpty(r.Data)
	s.Assert().Equal(128, r.Total)
}
