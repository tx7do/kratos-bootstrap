package zerolog

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

// NewStdoutWriter 返回默认的 stdout writer。
func NewStdoutWriter() io.Writer {
	return os.Stdout
}

// NewConsoleWriter 返回 zerolog.ConsoleWriter，timeFormat 为空时使用 zerolog 时间格式。
func NewConsoleWriter(timeFormat string) io.Writer {
	if timeFormat == "" {
		timeFormat = zerolog.TimeFieldFormat
	}
	return zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: timeFormat,
	}
}

// NewFileWriter 打开或创建文件并返回 writer，会创建父目录，mode 默认 0644。
func NewFileWriter(path string, mode os.FileMode) (io.Writer, error) {
	if path == "" {
		return nil, errors.New("path is empty")
	}
	dir := filepath.Dir(path)
	if dir != "." && dir != "/" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}
	if mode == 0 {
		mode = 0o644
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, mode)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// NewLumberjackWriter 使用 lumberjack 提供滚动写入支持。
func NewLumberjackWriter(path string, maxSizeMB, maxBackups, maxAge int, compress bool) io.Writer {
	return &lumberjack.Logger{
		Filename:   path,
		MaxSize:    maxSizeMB,  // megabytes
		MaxBackups: maxBackups, // number of backups
		MaxAge:     maxAge,     // days
		Compress:   compress,
	}
}

// NewMultiWriter 合并多个 writer。
func NewMultiWriter(writers ...io.Writer) io.Writer {
	return io.MultiWriter(writers...)
}

// NewWriter 根据 kind 创建不同类型的 io.Writer。
// 支持的 kinds:
//
//	"stdout"       -> NewStdoutWriter()
//	"console"      -> NewConsoleWriter (params["timeFormat"] string)
//	"file"         -> NewFileWriter (path 必填，params 可选: "mode" os.FileMode)
//	"lumberjack"   -> NewLumberjackWriter (path 必填, params: "maxSizeMB" int, "maxBackups" int, "maxAge" int, "compress" bool)
//	"multi"        -> NewMultiWriter (params["writers"] []io.Writer)
//
// params 可为 nil。
func NewWriter(kind, path string, params map[string]any) (io.Writer, error) {
	switch kind {
	case "stdout":
		return NewStdoutWriter(), nil

	case "console":
		tf := ""
		if params != nil {
			if v, ok := params["timeFormat"].(string); ok {
				tf = v
			}
		}
		return NewConsoleWriter(tf), nil

	case "file":
		mode := os.FileMode(0)
		if params != nil {
			if v, ok := params["mode"].(os.FileMode); ok {
				mode = v
			}
		}
		return NewFileWriter(path, mode)

	case "lumberjack":
		if path == "" {
			return nil, errors.New("path is required for lumberjack")
		}
		maxSize := 100
		maxBackups := 7
		maxAge := 30
		compress := false
		if params != nil {
			if v, ok := params["maxSizeMB"].(int); ok {
				maxSize = v
			}
			if v, ok := params["maxBackups"].(int); ok {
				maxBackups = v
			}
			if v, ok := params["maxAge"].(int); ok {
				maxAge = v
			}
			if v, ok := params["compress"].(bool); ok {
				compress = v
			}
		}
		return NewLumberjackWriter(path, maxSize, maxBackups, maxAge, compress), nil

	case "multi":
		if params == nil {
			return nil, errors.New("params[\"writers\"] is required for multi")
		}
		wrs, ok := params["writers"].([]io.Writer)
		if !ok {
			return nil, errors.New("params[\"writers\"] must be []io.Writer")
		}
		return NewMultiWriter(wrs...), nil

	default:
		return nil, errors.New("unsupported writer kind")
	}
}
