// Package log 基于 slog 的日志封装。
// 默认静默（仅 ERROR 级别输出），通过 SetLevel 控制可见度。
package log

import (
	"context"
	"io"
	"log/slog"
	"os"
)

// Level 日志级别
type Level = slog.Level

const (
	LevelDebug  = slog.LevelDebug
	LevelInfo   = slog.LevelInfo
	LevelWarn   = slog.LevelWarn
	LevelError  = slog.LevelError
	LevelSilent = slog.Level(12) // 高于所有级别，完全静默
)

var (
	defaultLogger = slog.New(slog.NewTextHandler(io.Discard, nil))
)

// SetLevel 设置全局日志级别。LevelSilent 完全静默。
func SetLevel(lvl Level) {
	h := defaultLogger.Handler()
	if lvl == LevelSilent {
		defaultLogger = slog.New(slog.NewTextHandler(io.Discard, nil))
		return
	}
	// 保持原有 handler，仅替换级别
	defaultLogger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: lvl,
	}))
	_ = h // 旧的 handler 将被 GC
}

// SetDefault 替换默认 logger（用于测试注入）
func SetDefault(logger *slog.Logger) {
	defaultLogger = logger
}

// Default 返回默认 logger
func Default() *slog.Logger {
	return defaultLogger
}

// Debug 输出调试日志
func Debug(msg string, args ...any) {
	defaultLogger.Debug(msg, args...)
}

// Info 输出信息日志
func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}

// Warn 输出警告日志
func Warn(msg string, args ...any) {
	defaultLogger.Warn(msg, args...)
}

// Error 输出错误日志
func Error(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
}

// DebugContext 带 context 的调试日志
func DebugContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.DebugContext(ctx, msg, args...)
}

// InfoContext 带 context 的信息日志
func InfoContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.InfoContext(ctx, msg, args...)
}

// WarnContext 带 context 的警告日志
func WarnContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.WarnContext(ctx, msg, args...)
}

// ErrorContext 带 context 的错误日志
func ErrorContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.ErrorContext(ctx, msg, args...)
}
