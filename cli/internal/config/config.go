// Package config 管理 CLI 配置：token 存储、proxy 设置。
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	tokenFileName  = "token.json"
	configFileName = "config.json"
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

// TokenPath 返回与二进制同目录的 token.json 路径
func TokenPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("获取可执行文件路径失败: %w", err)
	}
	return filepath.Join(filepath.Dir(exe), tokenFileName), nil
}

// LoadToken 从 token.json 读取 token
func LoadToken() (*TokenData, error) {
	path, err := TokenPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // 无 token 文件，返回 nil
		}
		return nil, fmt.Errorf("读取 token 文件失败: %w", err)
	}
	var td TokenData
	if err := json.Unmarshal(data, &td); err != nil {
		return nil, fmt.Errorf("解析 token 文件失败: %w", err)
	}
	if td.AccessToken == "" {
		return nil, nil
	}
	return &td, nil
}

// SaveToken 保存 token 到 token.json
func SaveToken(td *TokenData) error {
	path, err := TokenPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(td, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化 token 失败: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("写入 token 文件失败: %w", err)
	}
	return nil
}

// DeleteToken 删除 token.json
func DeleteToken() error {
	path, err := TokenPath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除 token 文件失败: %w", err)
	}
	return nil
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

// ConfigPath 返回二进制同目录的 config.json 路径
func ConfigPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("获取可执行文件路径失败: %w", err)
	}
	return filepath.Join(filepath.Dir(exe), configFileName), nil
}

// LoadConfig 读取 config.json
func LoadConfig() (*ConfigData, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &ConfigData{}, nil
		}
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}
	var cfg ConfigData
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}
	return &cfg, nil
}

// SaveConfig 保存 config.json
func SaveConfig(cfg *ConfigData) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}
	return nil
}
