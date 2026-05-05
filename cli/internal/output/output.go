// Package output 提供双格式输出：txt（AI Agent 优化）和 json。
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Format 输出格式
type Format string

const (
	FormatTxt  Format = "txt"
	FormatJSON Format = "json"
)

// Writer 封装输出格式和写入目标
type Writer struct {
	Format Format
	Out    io.Writer
}

// NewWriter 创建输出器
func NewWriter(f Format) *Writer {
	return &Writer{Format: f, Out: os.Stdout}
}

// Print 输出 data，根据格式自动选择 txt 或 json
// 若 data 实现了 fmt.Stringer，txt 模式调用 String()
// 若 data 实现了 JSON() ([]byte, error)，json 模式调用 JSON()
// 否则 json 模式使用 encoding/json 序列化
func (w *Writer) Print(data interface{}) error {
	if w.Format == FormatJSON {
		return w.printJSON(data)
	}
	return w.printTxt(data)
}

func (w *Writer) printTxt(data interface{}) error {
	switch v := data.(type) {
	case string:
		_, err := fmt.Fprintln(w.Out, v)
		return err
	case fmt.Stringer:
		_, err := fmt.Fprintln(w.Out, v.String())
		return err
	default:
		return w.printJSON(data) // 回退到 json
	}
}

func (w *Writer) printJSON(data interface{}) error {
	encoder := json.NewEncoder(w.Out)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// ReadFileOrValue 读取文件内容或直接返回值
// 若 filePath 不为空，读取文件内容；否则返回 value
func ReadFileOrValue(filePath, value string) (string, error) {
	if filePath != "" {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("读取文件 %s 失败: %w", filePath, err)
		}
		return string(data), nil
	}
	return value, nil
}

// AddFileFlag 为 cobra 命令添加 xxx-file 和 xxx 两个互斥标志
// 返回值：fileFlagPtr, valueFlagPtr
//func AddFileFlag(cmd *cobra.Command, flagName, desc string) (*string, *string) {
//	filePtr := cmd.Flags().String(flagName+"-file", "", desc+"（从文件读取，支持多行/长文本）")
//	valuePtr := cmd.Flags().String(flagName, "", desc)
//	return filePtr, valuePtr
//}
