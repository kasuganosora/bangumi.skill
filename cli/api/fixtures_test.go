//go:build fixtures
// +build fixtures

// 用法: go test -tags=fixtures -run TestGenerateFixtures -count=1 ./api/
// 使用真实 API Key 调用所有读端点，将响应保存为 mock 数据。

package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

const testAPIKey = "feQM4wUhAnWH6XDePR4tMKZmwNILgcj7ul5wkg9Z"
const testOutDir = "testdata"

var ctx = context.Background()

func TestGenerateFixtures(t *testing.T) {
	client, err := NewClient(
		WithAccessToken(testAPIKey),
		WithProxy("http://127.0.0.1:1081"),
	)
	if err != nil {
		t.Fatalf("create client: %v", err)
	}
	_ = os.MkdirAll(testOutDir, 0755)

	// ── Calendar ──
	doGet(t, client, "calendar", func() (interface{}, error) {
		return client.GetCalendar(ctx)
	})

	// ── Search ──
	doGet(t, client, "search_subjects", func() (interface{}, error) {
		return client.SearchSubjects(ctx, SearchSubjectRequest{Keyword: "EVA", Sort: "rank"}, 5, 0)
	})
	doGet(t, client, "search_characters", func() (interface{}, error) {
		return client.SearchCharacters(ctx, SearchCharacterRequest{Keyword: "saber"}, 5, 0)
	})
	doGet(t, client, "search_persons", func() (interface{}, error) {
		return client.SearchPersons(ctx, SearchPersonRequest{Keyword: "花泽香菜"}, 5, 0)
	})

	// ── Subject ──
	doGet(t, client, "subjects_anime", func() (interface{}, error) {
		return client.GetSubjects(ctx, SubjectAnime, 0, "rank", 0, 0, 5, 0)
	})
	doGet(t, client, "subject_12", func() (interface{}, error) {
		return client.GetSubjectByID(ctx, 12)
	})
	doGet(t, client, "subject_12_persons", func() (interface{}, error) {
		return client.GetSubjectPersons(ctx, 12)
	})
	doGet(t, client, "subject_12_characters", func() (interface{}, error) {
		return client.GetSubjectCharacters(ctx, 12)
	})
	doGet(t, client, "subject_12_relations", func() (interface{}, error) {
		return client.GetSubjectRelations(ctx, 12)
	})

	// ── Episode ──
	doGet(t, client, "episodes_12", func() (interface{}, error) {
		return client.GetEpisodes(ctx, 12, nil, 100, 0)
	})
	doGet(t, client, "episode_1027", func() (interface{}, error) {
		return client.GetEpisodeByID(ctx, 1027)
	})

	// ── Character ──
	doGet(t, client, "character_67", func() (interface{}, error) {
		return client.GetCharacterByID(ctx, 67)
	})
	doGet(t, client, "character_67_subjects", func() (interface{}, error) {
		return client.GetCharacterSubjects(ctx, 67)
	})
	doGet(t, client, "character_67_persons", func() (interface{}, error) {
		return client.GetCharacterPersons(ctx, 67)
	})

	// ── Person ──
	doGet(t, client, "person_1", func() (interface{}, error) {
		return client.GetPersonByID(ctx, 1)
	})
	doGet(t, client, "person_1_subjects", func() (interface{}, error) {
		return client.GetPersonSubjects(ctx, 1)
	})
	doGet(t, client, "person_1_characters", func() (interface{}, error) {
		return client.GetPersonCharacters(ctx, 1)
	})

	// ── User ──
	doGet(t, client, "user_sai", func() (interface{}, error) {
		return client.GetUserByName(ctx, "sai")
	})
	doGet(t, client, "me", func() (interface{}, error) {
		return client.GetMe(ctx)
	})

	// ── Collections ──
	doGet(t, client, "user_sai_collections", func() (interface{}, error) {
		return client.GetUserCollections(ctx, "sai", nil, nil, 5, 0)
	})

	// ── Revisions ── (需要 entity ID，暂用 subject_id=12)
	doGet(t, client, "subject_revisions", func() (interface{}, error) {
		return nil, nil // skip: needs subject_id query param
	})
	doGet(t, client, "character_revisions", func() (interface{}, error) {
		return nil, nil // skip: needs character_id query param
	})

	// ── Index ──
	doGet(t, client, "index_1", func() (interface{}, error) {
		return client.GetIndexByID(ctx, 1)
	})
	doGet(t, client, "index_1_subjects", func() (interface{}, error) {
		return client.GetIndexSubjects(ctx, 1, 5, 0)
	})

	// ── Legacy Search ── (使用原始 HTTP 响应，因为其响应格式特殊: {results, list})
	doRawGet(t, client, "legacy_search", "/search/subject/EVA?responseGroup=small&max_results=5")

	t.Logf("✅ 所有 fixtures 已保存到 %s/", testOutDir)
}

// doRawGet 用真实 token 直接请求原始 HTTP 端点并保存响应体（用于格式特殊的端点）
func doRawGet(t *testing.T, client *HTTPClient, name string, uri string) {
	t.Helper()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, client.baseURL.String()+"/"+uri, nil)
	req.Header.Set("User-Agent", client.userAgent)
	req.Header.Set("Authorization", "Bearer "+testAPIKey)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		t.Errorf("❌ raw %s: %v", name, err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	path := filepath.Join(testOutDir, name+".json")
	if err := os.WriteFile(path, body, 0644); err != nil {
		t.Errorf("write raw %s: %v", name, err)
		return
	}
	t.Logf("  ✅ raw %s (%d bytes)", name, len(body))
}

func doGet(t *testing.T, client *HTTPClient, name string, fn func() (interface{}, error)) {
	t.Helper()
	data, err := fn()
	if err != nil {
		t.Errorf("❌ %s: %v", name, err)
		return
	}
	if data == nil {
		t.Logf("  ⏭️ %s (skipped)", name)
		return
	}
	path := filepath.Join(testOutDir, name+".json")
	b, _ := json.MarshalIndent(data, "", "  ")
	if err := os.WriteFile(path, b, 0644); err != nil {
		t.Errorf("写入 %s 失败: %v", name, err)
		return
	}
	t.Logf("  ✅ %s (%d bytes)", name, len(b))
}
