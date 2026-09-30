package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPatchUserEpisodeCollections_Batches(t *testing.T) {
	var bodies []map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "/v0/users/-/collections/12/episodes", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		raw, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		var body map[string]interface{}
		require.NoError(t, json.Unmarshal(raw, &body))
		bodies = append(bodies, body)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client, err := NewClient(WithBaseURL(srv.URL), WithAccessToken("token"))
	require.NoError(t, err)

	ids := make([]int, 101)
	for i := range ids {
		ids[i] = i + 1
	}
	err = client.PatchUserEpisodeCollections(context.Background(), 12, ids, EpCollectionDone)
	require.NoError(t, err)
	require.Len(t, bodies, 2)
	assert.Len(t, bodies[0]["episode_id"], 100)
	assert.Len(t, bodies[1]["episode_id"], 1)
	assert.EqualValues(t, 2, bodies[0]["type"])
}

func TestPatchUserEpisodeCollections_Empty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("empty batch should not call the API")
	}))
	defer srv.Close()

	client, err := NewClient(WithBaseURL(srv.URL))
	require.NoError(t, err)
	require.NoError(t, client.PatchUserEpisodeCollections(context.Background(), 12, nil, EpCollectionDone))
}
