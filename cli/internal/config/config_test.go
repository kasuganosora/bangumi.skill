package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func useTempDirs(t *testing.T) (userDir, legacyDir string) {
	t.Helper()
	userDir = t.TempDir()
	legacyDir = t.TempDir()
	prevConfig, prevLegacy := configDirFunc, legacyDirFunc
	configDirFunc = func() (string, error) { return userDir, nil }
	legacyDirFunc = func() (string, error) { return legacyDir, nil }
	t.Cleanup(func() {
		configDirFunc = prevConfig
		legacyDirFunc = prevLegacy
	})
	return userDir, legacyDir
}

func TestLoadToken_MigratesLegacyFile(t *testing.T) {
	userDir, legacyDir := useTempDirs(t)
	t.Setenv(envToken, "")

	legacy := filepath.Join(legacyDir, tokenFileName)
	require.NoError(t, os.WriteFile(legacy, []byte(`{"access_token":"legacy-token","user_id":7}`), 0o600))

	td, err := LoadToken()
	require.NoError(t, err)
	require.NotNil(t, td)
	assert.Equal(t, "legacy-token", td.AccessToken)
	assert.Equal(t, 7, td.UserID)

	migrated, err := os.ReadFile(filepath.Join(userDir, tokenFileName))
	require.NoError(t, err)
	assert.Contains(t, string(migrated), "legacy-token")
}

func TestLoadToken_EnvOverridesFile(t *testing.T) {
	userDir, _ := useTempDirs(t)
	require.NoError(t, os.WriteFile(filepath.Join(userDir, tokenFileName), []byte(`{"access_token":"file-token"}`), 0o600))
	t.Setenv(envToken, "env-token")

	td, err := LoadToken()
	require.NoError(t, err)
	assert.Equal(t, "env-token", td.AccessToken)
	assert.True(t, TokenFromEnv())
}

func TestSaveAndDeleteToken(t *testing.T) {
	userDir, legacyDir := useTempDirs(t)
	t.Setenv(envToken, "")

	require.NoError(t, os.WriteFile(filepath.Join(legacyDir, tokenFileName), []byte(`{"access_token":"old"}`), 0o600))
	require.NoError(t, SaveToken(&TokenData{AccessToken: "new", UserID: 3}))

	td, err := LoadToken()
	require.NoError(t, err)
	assert.Equal(t, "new", td.AccessToken)

	require.NoError(t, DeleteToken())
	td, err = LoadToken()
	require.NoError(t, err)
	assert.Nil(t, td)
	_, err = os.Stat(filepath.Join(userDir, tokenFileName))
	assert.True(t, os.IsNotExist(err))
	_, err = os.Stat(filepath.Join(legacyDir, tokenFileName))
	assert.True(t, os.IsNotExist(err))
}

func TestLoadConfig_MigratesLegacyFile(t *testing.T) {
	userDir, legacyDir := useTempDirs(t)
	require.NoError(t, os.WriteFile(filepath.Join(legacyDir, configFileName), []byte(`{"proxy":"http://127.0.0.1:1081","timeout":15}`), 0o600))

	cfg, err := LoadConfig()
	require.NoError(t, err)
	assert.Equal(t, "http://127.0.0.1:1081", cfg.Proxy)
	assert.Equal(t, 15, cfg.Timeout)

	migrated, err := os.ReadFile(filepath.Join(userDir, configFileName))
	require.NoError(t, err)
	assert.Contains(t, string(migrated), "1081")
}
