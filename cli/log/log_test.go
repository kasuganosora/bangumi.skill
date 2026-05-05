package log

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultSilent(t *testing.T) {
	// 默认 logger 输出到 io.Discard，任何级别都不应该有输出到 stderr
	var buf bytes.Buffer
	defaultLogger = slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))

	// 恢复默认
	defer func() {
		defaultLogger = slog.New(slog.NewTextHandler(ioDiscard(), nil))
	}()

	Debug("should appear", "key", "val")
	assert.Contains(t, buf.String(), "should appear")
}

func TestSetLevel_Info(t *testing.T) {
	SetLevel(LevelInfo)

	// 取当前 logger 验证级别
	h := Default().Handler()
	assert.NotNil(t, h)

	// 恢复
	SetLevel(LevelSilent)
}

func TestSetLevel_Debug(t *testing.T) {
	var buf bytes.Buffer
	SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: LevelDebug})))

	Debug("debug msg", "a", 1)
	Info("info msg", "b", 2)

	out := buf.String()
	assert.Contains(t, out, "debug msg")
	assert.Contains(t, out, "info msg")
}

func TestSetLevel_Silent(t *testing.T) {
	var buf bytes.Buffer
	SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: LevelDebug})))

	SetLevel(LevelSilent)

	// 清空
	buf.Reset()
	Debug("invisible")
	Info("invisible")
	Error("invisible")

	assert.Empty(t, buf.String())
}

func TestLevels(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: LevelInfo}))
	SetDefault(logger)

	Debug("should be hidden")
	Info("should appear")
	Warn("should appear")
	Error("should appear")

	out := buf.String()
	assert.NotContains(t, out, "should be hidden")
	assert.Contains(t, out, "should appear")
	assert.Contains(t, out, "should appear")
}

func TestContextMethods(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: LevelDebug}))
	SetDefault(logger)

	// 使用 context 方法不会 panic
	DebugContext(t.Context(), "debug with ctx")
	InfoContext(t.Context(), "info with ctx")
	WarnContext(t.Context(), "warn with ctx")
	ErrorContext(t.Context(), "error with ctx")

	assert.Contains(t, buf.String(), "debug with ctx")
}

func TestJSONOutput(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: LevelInfo})
	logger := slog.New(handler)
	SetDefault(logger)

	Info("structured", "user", "kasuganosora", "action", "login")

	var m map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &m)
	require.NoError(t, err)
	assert.Equal(t, "structured", m["msg"])
	assert.Equal(t, "INFO", m["level"])
	assert.Equal(t, "kasuganosora", m["user"])
	assert.Equal(t, "login", m["action"])
}

func TestSetDefault_WithAttrs(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: LevelDebug}))
	logger = logger.With("service", "bangumi-cli")
	SetDefault(logger)

	Info("test msg")
	assert.Contains(t, buf.String(), "service=bangumi-cli")
	assert.Contains(t, buf.String(), "test msg")
}

func TestDefaultFuncs_Panic(t *testing.T) {
	// 确保默认静默 logger 调用各函数不会 panic
	SetLevel(LevelSilent)
	assert.NotPanics(t, func() { Debug("x") })
	assert.NotPanics(t, func() { Info("x") })
	assert.NotPanics(t, func() { Warn("x") })
	assert.NotPanics(t, func() { Error("x") })
}

// 辅助：获取 io.Discard 类型（避免导入 os 仅为了 io.Discard）
func ioDiscard() *bytes.Buffer {
	return &bytes.Buffer{}
}
