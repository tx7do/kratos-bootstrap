package script_engine

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	scriptEngine "github.com/tx7do/go-scripts"
)

// Compile-time assertion: KratosLogger implements scriptEngine.Logger.
var _ scriptEngine.Logger = (*KratosLogger)(nil)

// KratosLogger 将 Kratos log.Logger 适配为 go-scripts Logger 接口。
//
// 使得 go-scripts 内部日志（引擎初始化、脚本执行、热更新等）能够
// 无缝接入 Kratos 统一日志体系。
//
// 使用方式:
//
//	import "github.com/go-kratos/kratos/v2/log"
//	import scriptEngine "github.com/tx7do/go-scripts"
//
//	kratosLogger := log.NewStdLogger(os.Stdout)
//	scriptEngine.SetLogger(&script_engine.KratosLogger{Logger: kratosLogger})
type KratosLogger struct {
	// Logger 是 Kratos 日志实例。
	Logger log.Logger
}

// Debug 以 DEBUG 级别记录日志。
func (l *KratosLogger) Debug(ctx context.Context, msg string, args ...any) {
	if l.Logger == nil {
		return
	}
	keyvals := buildKeyVals(msg, args)
	_ = l.Logger.Log(log.LevelDebug, keyvals...)
}

// Info 以 INFO 级别记录日志。
func (l *KratosLogger) Info(ctx context.Context, msg string, args ...any) {
	if l.Logger == nil {
		return
	}
	keyvals := buildKeyVals(msg, args)
	_ = l.Logger.Log(log.LevelInfo, keyvals...)
}

// Warn 以 WARN 级别记录日志。
func (l *KratosLogger) Warn(ctx context.Context, msg string, args ...any) {
	if l.Logger == nil {
		return
	}
	keyvals := buildKeyVals(msg, args)
	_ = l.Logger.Log(log.LevelWarn, keyvals...)
}

// Error 以 ERROR 级别记录日志。
func (l *KratosLogger) Error(ctx context.Context, msg string, args ...any) {
	if l.Logger == nil {
		return
	}
	keyvals := buildKeyVals(msg, args)
	_ = l.Logger.Log(log.LevelError, keyvals...)
}

// With 返回附加了指定键值对的新 KratosLogger 实例。
// 常用于区分模块，例如 logger.With("module", "lua")。
func (l *KratosLogger) With(args ...any) scriptEngine.Logger {
	if l.Logger == nil {
		return l
	}
	return &KratosLogger{Logger: log.With(l.Logger, args...)}
}

// SetLogger 设置 go-scripts 全局日志为 Kratos Logger。
// 传入 nil 可恢复为默认的静默日志。
//
// 便捷封装，等同于:
//
//	scriptEngine.SetLogger(&KratosLogger{Logger: kratosLogger})
func SetLogger(kratosLogger log.Logger) {
	if kratosLogger == nil {
		scriptEngine.SetLogger(nil)
		return
	}
	scriptEngine.SetLogger(&KratosLogger{Logger: kratosLogger})
}

// buildKeyVals 将 go-scripts 风格的 (msg, args...) 转换为 Kratos 风格的 keyvals。
// go-scripts: msg + [k1, v1, k2, v2, ...]
// kratos:     ["msg", msg, k1, v1, k2, v2, ...]
func buildKeyVals(msg string, args []any) []any {
	keyvals := make([]any, 0, 2+len(args))
	keyvals = append(keyvals, "msg", msg)
	// args 已是 alternating key-value pairs
	keyvals = append(keyvals, args...)
	return keyvals
}
