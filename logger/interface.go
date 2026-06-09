package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"go.opentelemetry.io/otel/trace"
)

// Level 日志级别。
//
// 与 Kratos log.Level 数值兼容，方便适配器零成本转换。
type Level int32

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// String 返回日志级别的可读字符串。
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger 项目统一日志接口。
//
// 设计目标：
//   - context 首参：支持 tracing / request-scoped 属性，取代 Kratos Valuer 机制
//   - 每级独立方法：API 清晰，编译期类型安全，IDE 友好
//   - msg 独立参数：不埋在 keyvals 中，调用更直观
//   - keyvals 结构化：交替 key-value 对，兼容 Kratos / slog 生态
//   - With() 返回新实例：可附加固定字段（如 module / service.*）
type Logger interface {
	// Debug 输出 DEBUG 级别日志。
	// ctx 用于提取 trace_id / span_id 等请求级属性，可为 nil。
	// keyvals 为交替的 key-value 对，例如 ("user", 42, "action", "login")。
	Debug(ctx context.Context, msg string, keyvals ...any)

	// Info 输出 INFO 级别日志。
	Info(ctx context.Context, msg string, keyvals ...any)

	// Warn 输出 WARN 级别日志。
	Warn(ctx context.Context, msg string, keyvals ...any)

	// Error 输出 ERROR 级别日志。
	Error(ctx context.Context, msg string, keyvals ...any)

	// With 返回附加了固定 key-value 对的新 Logger 实例。
	// 常用于区分模块：logger.With("module", "bootstrap")。
	With(keyvals ...any) Logger
}

// ============================================================================
// nopLogger — 零开销的静默实现
// ============================================================================

// nopLogger 丢弃所有日志记录。导入包时无副作用，调用方通过 [SetLogger] 注入实例。
type nopLogger struct{}

func (nopLogger) Debug(context.Context, string, ...any) {}
func (nopLogger) Info(context.Context, string, ...any)  {}
func (nopLogger) Warn(context.Context, string, ...any)  {}
func (nopLogger) Error(context.Context, string, ...any) {}
func (n nopLogger) With(...any) Logger                  { return n }

// Compile-time assertion: nopLogger implements Logger.
var _ Logger = nopLogger{}

// NopLogger 返回一个丢弃所有日志的 Logger 实例。
func NopLogger() Logger { return nopLogger{} }

// ============================================================================
// stdLogger — 基于标准库 slog 的实现（无 Kratos 依赖）
// ============================================================================

// stdLogger 使用标准库 log/slog 作为后端。
// 这是项目 Logger 的参考实现，不依赖任何第三方日志框架。
type stdLogger struct {
	l     *slog.Logger
	extra []any
}

// Compile-time assertion: stdLogger implements Logger.
var _ Logger = (*stdLogger)(nil)

// Debug forwards to slog.Logger.DebugContext.
func (s *stdLogger) Debug(ctx context.Context, msg string, keyvals ...any) {
	s.l.DebugContext(orBackground(ctx), msg, s.kvs(keyvals)...)
}

// Info forwards to slog.Logger.InfoContext.
func (s *stdLogger) Info(ctx context.Context, msg string, keyvals ...any) {
	s.l.InfoContext(orBackground(ctx), msg, s.kvs(keyvals)...)
}

// Warn forwards to slog.Logger.WarnContext.
func (s *stdLogger) Warn(ctx context.Context, msg string, keyvals ...any) {
	s.l.WarnContext(orBackground(ctx), msg, s.kvs(keyvals)...)
}

// Error forwards to slog.Logger.ErrorContext.
func (s *stdLogger) Error(ctx context.Context, msg string, keyvals ...any) {
	s.l.ErrorContext(orBackground(ctx), msg, s.kvs(keyvals)...)
}

// With 返回附加了指定 key-value 对的新 stdLogger 实例。
func (s *stdLogger) With(keyvals ...any) Logger {
	return &stdLogger{l: s.l.With(keyvals...), extra: mergeKvs(s.extra, keyvals)}
}

// kvs 合并 extra（来自 With）与本次调用的 keyvals。
func (s *stdLogger) kvs(keyvals []any) []any {
	if len(s.extra) == 0 {
		return keyvals
	}
	return mergeKvs(s.extra, keyvals)
}

// NewStdLogger 创建一个基于 slog 的 Logger，输出到 stderr，级别为 Info。
//
// 这是项目 Logger 接口的原生实现，不依赖 Kratos。
// 调用方需要不同格式/级别/输出位置时，可以：
//   - 构建自己的 *slog.Logger 并直接赋值给 stdLogger{L: myLogger}
//   - 或自行实现 Logger 接口
func NewStdLogger() Logger {
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})
	return &stdLogger{l: slog.New(h)}
}

// NewSlogLogger 使用指定的 *slog.Logger 创建 Logger。
func NewSlogLogger(l *slog.Logger) Logger {
	if l == nil {
		return nopLogger{}
	}
	return &stdLogger{l: l}
}

// ============================================================================
// traceLogger — 自动注入 trace_id / span_id
//
// 替代 Kratos Valuer 机制。项目 Logger 接口携带 ctx，可以直接从中提取
// OTel trace 信息，不受 Kratos log.Logger 不接收 ctx 的限制。
// ============================================================================

// traceLogger 包装 Logger，在每条日志中自动注入 trace_id / span_id。
type traceLogger struct {
	Logger
}

var _ Logger = (*traceLogger)(nil)

// WrapTrace 用 trace_id/span_id 自动注入包装一个 Logger。
// 重复包装是安全的（会检测并跳过）。
func WrapTrace(l Logger) Logger {
	if l == nil {
		return nil
	}
	if _, ok := l.(*traceLogger); ok {
		return l
	}
	return &traceLogger{Logger: l}
}

func (t *traceLogger) Debug(ctx context.Context, msg string, keyvals ...any) {
	t.Logger.Debug(ctx, msg, appendTrace(ctx, keyvals)...)
}

func (t *traceLogger) Info(ctx context.Context, msg string, keyvals ...any) {
	t.Logger.Info(ctx, msg, appendTrace(ctx, keyvals)...)
}

func (t *traceLogger) Warn(ctx context.Context, msg string, keyvals ...any) {
	t.Logger.Warn(ctx, msg, appendTrace(ctx, keyvals)...)
}

func (t *traceLogger) Error(ctx context.Context, msg string, keyvals ...any) {
	t.Logger.Error(ctx, msg, appendTrace(ctx, keyvals)...)
}

func (t *traceLogger) With(keyvals ...any) Logger {
	return &traceLogger{Logger: t.Logger.With(keyvals...)}
}

// appendTrace 从 ctx 提取 trace_id/span_id 并前置到 keyvals。
func appendTrace(ctx context.Context, keyvals []any) []any {
	if ctx == nil {
		return keyvals
	}
	traceID, spanID := extractTrace(ctx)
	if traceID == "" && spanID == "" {
		return keyvals
	}
	out := make([]any, 0, len(keyvals)+4)
	if traceID != "" {
		out = append(out, "trace_id", traceID)
	}
	if spanID != "" {
		out = append(out, "span_id", spanID)
	}
	return append(out, keyvals...)
}

// extractTrace 从 ctx 中提取 OTel trace_id 和 span_id。
func extractTrace(ctx context.Context) (traceID, spanID string) {
	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() {
		return "", ""
	}
	if spanCtx.HasTraceID() {
		traceID = spanCtx.TraceID().String()
	}
	if spanCtx.HasSpanID() {
		spanID = spanCtx.SpanID().String()
	}
	return
}

// ============================================================================
// Helper — 便捷日志助手
// ============================================================================

// Helper 提供格式化和简化调用方法。
// 等价于 Kratos log.Helper，但不依赖 Kratos。
type Helper struct {
	Logger
}

// NewHelper 创建日志助手。
func NewHelper(l Logger) *Helper {
	if l == nil {
		l = nopLogger{}
	}
	return &Helper{Logger: l}
}

// Debugf 格式化输出 DEBUG 日志。
func (h *Helper) Debugf(ctx context.Context, format string, args ...any) {
	h.Debug(ctx, fmt.Sprintf(format, args...))
}

// Infof 格式化输出 INFO 日志。
func (h *Helper) Infof(ctx context.Context, format string, args ...any) {
	h.Info(ctx, fmt.Sprintf(format, args...))
}

// Warnf 格式化输出 WARN 日志。
func (h *Helper) Warnf(ctx context.Context, format string, args ...any) {
	h.Warn(ctx, fmt.Sprintf(format, args...))
}

// Errorf 格式化输出 ERROR 日志。
func (h *Helper) Errorf(ctx context.Context, format string, args ...any) {
	h.Error(ctx, fmt.Sprintf(format, args...))
}

// ============================================================================
// 全局 Logger 状态
// ============================================================================

var (
	globalMu  sync.RWMutex
	globalLog Logger = nopLogger{}
)

// SetLogger 设置包级全局 Logger。传 nil 恢复为静默 nopLogger。
// 并发安全。
func SetLogger(l Logger) {
	globalMu.Lock()
	defer globalMu.Unlock()
	if l == nil {
		globalLog = nopLogger{}
		return
	}
	globalLog = l
}

// GetLogger 返回当前包级 Logger。
func GetLogger() Logger {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return globalLog
}

// ============================================================================
// 内部工具函数
// ============================================================================

// orBackground 确保 context 非 nil。
func orBackground(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

// mergeKvs 合并两组 key-value 对。
func mergeKvs(a, b []any) []any {
	if len(a) == 0 {
		return b
	}
	if len(b) == 0 {
		return a
	}
	out := make([]any, 0, len(a)+len(b))
	out = append(out, a...)
	out = append(out, b...)
	return out
}
