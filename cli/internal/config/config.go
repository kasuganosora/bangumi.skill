// Package config 管理 CLI 配置：token 存储、proxy 设置。
//
// 令牌和配置写在用户配置目录（Linux 为 ~/.config/bangumi/）。
// 仍会读取二进制旁边的旧文件，读到后迁到新目录。
// 环境变量 BANGUMI_TOKEN 优先于文件中的令牌。
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	tokenFileName  = "token.json"
	configFileName = "config.json"
	appDirName     = "bangumi"
	envToken       = "BANGUMI_TOKEN"
)

// TokenData 存储的 token 信息
type TokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserID       int    `json:"user_id"`
	ExpiresIn    int    `json:"expires_in"`
}

// ConfigData CLI 全局配置
type ConfigData struct {
	Proxy   string `json:"proxy,omitempty"`
	Timeout int    `json:"timeout,omitempty"` // 秒，0 表示未设置
}

// configDirFunc 返回配置目录，测试可替换。
var configDirFunc = defaultConfigDir

// legacyDirFunc 返回旧版「二进制所在目录」，测试可替换。
var legacyDirFunc = defaultLegacyDir

func defaultConfigDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("获取用户配置目录失败: %w", err)
	}
	return filepath.Join(base, appDirName), nil
}

func defaultLegacyDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("获取可执行文件路径失败: %w", err)
	}
	return filepath.Dir(exe), nil
}

func ensureConfigDir() (string, error) {
	dir, err := configDirFunc()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", fmt.Errorf("创建配置目录失败: %w", err)
	}
	return dir, nil
}

// ConfigDir 返回用户配置目录，必要时创建。
func ConfigDir() (string, error) {
	return ensureConfigDir()
}

// TokenFromEnv 报告当前是否由环境变量提供令牌。
func TokenFromEnv() bool {
	return strings.TrimSpace(os.Getenv(envToken)) != ""
}

// TokenPath 返回用户配置目录下的 token.json 路径。
func TokenPath() (string, error) {
	dir, err := ensureConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, tokenFileName), nil
}

func legacyTokenPath() (string, error) {
	dir, err := legacyDirFunc()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, tokenFileName), nil
}

// LoadToken 读取令牌。优先 BANGUMI_TOKEN，其次用户配置目录，最后是二进制旁边的旧文件。
func LoadToken() (*TokenData, error) {
	if tok := strings.TrimSpace(os.Getenv(envToken)); tok != "" {
		return &TokenData{AccessToken: tok}, nil
	}
	path, err := TokenPath()
	if err != nil {
		return nil, err
	}
	td, found, err := readToken(path)
	if err != nil || found {
		return td, err
	}
	legacy, err := legacyTokenPath()
	if err != nil {
		return nil, nil
	}
	td, found, err = readToken(legacy)
	if err != nil || !found {
		return td, err
	}
	if saveErr := SaveToken(td); saveErr != nil {
		return td, nil
	}
	return td, nil
}

func readToken(path string) (*TokenData, bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("读取 token 文件失败: %w", err)
	}
	var td TokenData
	if err := json.Unmarshal(data, &td); err != nil {
		return nil, false, fmt.Errorf("解析 token 文件失败: %w", err)
	}
	if td.AccessToken == "" {
		return nil, false, nil
	}
	return &td, true, nil
}

// SaveToken 保存 token 到用户配置目录。
func SaveToken(td *TokenData) error {
	path, err := TokenPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(td, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化 token 失败: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("写入 token 文件失败: %w", err)
	}
	return nil
}

// DeleteToken 删除用户配置目录和旧位置的 token.json。
func DeleteToken() error {
	var errs []error
	paths := []func() (string, error){TokenPath, legacyTokenPath}
	for _, pathFn := range paths {
		path, err := pathFn()
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			errs = append(errs, fmt.Errorf("删除 token 文件失败: %w", err))
		}
	}
	return errors.Join(errs...)
}

// RequireToken 要求 token，无 token 时输出提示并返回错误
func RequireToken() (*TokenData, error) {
	td, err := LoadToken()
	if err != nil {
		return nil, err
	}
	if td == nil || td.AccessToken == "" {
		return nil, fmt.Errorf(
			"未设置个人令牌\n\n请先申请个人令牌: https://next.bgm.tv/demo/access-token\n然后运行: bangumi auth login --token <你的令牌>",
		)
	}
	return td, nil
}

// ---------------------------------------------------------------------------
// config.json (proxy 等全局配置)
// ---------------------------------------------------------------------------

// ConfigPath 返回用户配置目录下的 config.json 路径。
func ConfigPath() (string, error) {
	dir, err := ensureConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, configFileName), nil
}

func legacyConfigPath() (string, error) {
	dir, err := legacyDirFunc()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, configFileName), nil
}

// LoadConfig 读取 config.json。用户目录没有文件时，尝试旧位置并迁过来。
func LoadConfig() (*ConfigData, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}
	cfg, found, err := readConfig(path)
	if err != nil || found {
		return cfg, err
	}
	legacy, err := legacyConfigPath()
	if err != nil {
		return &ConfigData{}, nil
	}
	cfg, found, err = readConfig(legacy)
	if err != nil || !found {
		if err != nil {
			return nil, err
		}
		return &ConfigData{}, nil
	}
	if saveErr := SaveConfig(cfg); saveErr != nil {
		return cfg, nil
	}
	return cfg, nil
}

func readConfig(path string) (*ConfigData, bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &ConfigData{}, false, nil
		}
		return nil, false, fmt.Errorf("读取配置文件失败: %w", err)
	}
	var cfg ConfigData
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, false, fmt.Errorf("解析配置文件失败: %w", err)
	}
	return &cfg, true, nil
}

// SaveConfig 保存 config.json 到用户配置目录。
func SaveConfig(cfg *ConfigData) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}
	return nil
}
