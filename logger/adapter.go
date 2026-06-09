package logger

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// msgKey 返回 Kratos keyvals 中标识消息体的 key。
// 与 Kratos log.DefaultMessageKey 保持一致。
func msgKey() string {
	return log.DefaultMessageKey
}

// ============================================================================
// Forward 适配器：项目 Logger → Kratos log.Logger
//
// 当需要把项目 Logger 传递给 kratos.App、log.NewHelper 等
// 需要 log.Logger 的场景时，使用 AsKratosLogger 包装。
// ============================================================================

// kratosAdapter 将项目 Logger 适配为 Kratos log.Logger。
type kratosAdapter struct {
	l Logger
}

// Compile-time assertion: kratosAdapter implements log.Logger.
var _ log.Logger = (*kratosAdapter)(nil)

// Log 实现 Kratos log.Logger 接口。
//
// 从 Kratos 风格的 keyvals 中提取 msg，将剩余 key-value 对转发给项目 Logger。
// 同时支持 Kratos Valuer 求值。
func (a *kratosAdapter) Log(level log.Level, keyvals ...any) error {
	all := keyvals

	// 提取 msg 和过滤 keyvals
	msg := ""
	userKvs := make([]any, 0, len(all))
	for i := 0; i+1 < len(all); i += 2 {
		key, ok := all[i].(string)
		if !ok {
			continue
		}
		// 求值 Valuer
		val := all[i+1]
		if v, ok := val.(log.Valuer); ok {
			val = v(context.Background())
		}
		if key == msgKey() {
			if s, ok := val.(string); ok {
				msg = s
			}
			continue
		}
		userKvs = append(userKvs, key, val)
	}

	ctx := context.Background()
	switch level {
	case log.LevelDebug:
		a.l.Debug(ctx, msg, userKvs...)
	case log.LevelInfo:
		a.l.Info(ctx, msg, userKvs...)
	case log.LevelWarn:
		a.l.Warn(ctx, msg, userKvs...)
	case log.LevelError:
		a.l.Error(ctx, msg, userKvs...)
	case log.LevelFatal:
		// Logger 接口没有 Fatal 方法，降级为 Error
		a.l.Error(ctx, msg, userKvs...)
	default:
		a.l.Info(ctx, msg, userKvs...)
	}
	return nil
}

// AsKratosLogger 将项目 Logger 适配为 Kratos log.Logger。
//
// 用于兼容仍需要 log.Logger 的场景，例如：
//
//	kratosApp := kratos.New(kratos.Logger(logger.AsKratosLogger(myLogger)))
//	helper := log.NewHelper(logger.AsKratosLogger(myLogger))
//
// 如果传入的 Logger 已经是包装的 Kratos logger（通过 FromKratosLogger 创建），
// 会直接解包返回底层的 log.Logger，避免双重包装。
func AsKratosLogger(l Logger) log.Logger {
	if l == nil {
		return nil
	}
	// 如果已经是包装的 kratos logger，直接解包
	if w, ok := l.(*kratosWrapper); ok {
		return w.l
	}
	return &kratosAdapter{l: l}
}

// ============================================================================
// Reverse 适配器：Kratos log.Logger → 项目 Logger
//
// 当需要把现有的 Kratos log.Logger 实现（zap/logrus/fluent 等）
// 接入项目 Logger 体系时，使用 FromKratosLogger 包装。
// ============================================================================

// kratosWrapper 包装 Kratos log.Logger 为项目 Logger。
type kratosWrapper struct {
	l log.Logger
}

// Compile-time assertion: kratosWrapper implements Logger.
var _ Logger = (*kratosWrapper)(nil)

// Debug 以 DEBUG 级别记录日志。
func (w *kratosWrapper) Debug(ctx context.Context, msg string, keyvals ...any) {
	w.log(ctx, log.LevelDebug, msg, keyvals)
}

// Info 以 INFO 级别记录日志。
func (w *kratosWrapper) Info(ctx context.Context, msg string, keyvals ...any) {
	w.log(ctx, log.LevelInfo, msg, keyvals)
}

// Warn 以 WARN 级别记录日志。
func (w *kratosWrapper) Warn(ctx context.Context, msg string, keyvals ...any) {
	w.log(ctx, log.LevelWarn, msg, keyvals)
}

// Error 以 ERROR 级别记录日志。
func (w *kratosWrapper) Error(ctx context.Context, msg string, keyvals ...any) {
	w.log(ctx, log.LevelError, msg, keyvals)
}

// With 返回附加了指定 key-value 对的新 Logger 实例。
func (w *kratosWrapper) With(keyvals ...any) Logger {
	// 使用 Kratos log.With 包装底层 logger
	newL := log.With(w.l, keyvals...)
	return &kratosWrapper{l: newL}
}

// log 内部方法，统一处理 keyvals 构建 + Valuer 求值 + 转发。
func (w *kratosWrapper) log(ctx context.Context, level log.Level, msg string, keyvals []any) {
	// 构建 Kratos 风格的 keyvals: [msgKey, msg, k1, v1, ...]
	all := make([]any, 0, 2+len(keyvals))
	all = append(all, msgKey(), msg)
	all = append(all, keyvals...)

	// 求值 Kratos Valuer（如 DefaultTimestamp、tracing.TraceID()）
	// 使用请求 context 求值，以正确提取 trace_id 等信息
	evalValuers(ctx, all)

	_ = w.l.Log(level, all...)
}

// evalValuers 遍历 keyvals 中的值，若为 log.Valuer 则用 ctx 求值替换。
func evalValuers(ctx context.Context, keyvals []any) {
	if ctx == nil {
		ctx = context.Background()
	}
	for i := 1; i < len(keyvals); i += 2 {
		if v, ok := keyvals[i].(log.Valuer); ok {
			keyvals[i] = v(ctx)
		}
	}
}

// FromKratosLogger 包装 Kratos log.Logger 为项目 Logger。
//
// 使用方式：
//
//	// 现有的 Kratos zap logger
//	kratosL, _ := zap.NewLogger(cfg)
//	// 包装为项目 Logger
//	myLogger := logger.FromKratosLogger(kratosL)
func FromKratosLogger(l log.Logger) Logger {
	if l == nil {
		return nopLogger{}
	}
	return &kratosWrapper{l: l}
}
