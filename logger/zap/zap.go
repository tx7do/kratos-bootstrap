package zap

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	bLogger "github.com/tx7do/kratos-bootstrap/logger"
)

// Compile-time assertion: Logger implements bLogger.Logger.
var _ bLogger.Logger = (*Logger)(nil)

type Logger struct {
	log *zap.Logger
}

func NewZapLogger(zlog *zap.Logger) *Logger {
	return &Logger{log: zlog}
}

// Debug 输出 DEBUG 级别日志。
func (l *Logger) Debug(_ context.Context, msg string, keyvals ...any) {
	l.log.Debug(msg, l.toFields(keyvals)...)
}

// Info 输出 INFO 级别日志。
func (l *Logger) Info(_ context.Context, msg string, keyvals ...any) {
	l.log.Info(msg, l.toFields(keyvals)...)
}

// Warn 输出 WARN 级别日志。
func (l *Logger) Warn(_ context.Context, msg string, keyvals ...any) {
	l.log.Warn(msg, l.toFields(keyvals)...)
}

// Error 输出 ERROR 级别日志。
func (l *Logger) Error(_ context.Context, msg string, keyvals ...any) {
	l.log.Error(msg, l.toFields(keyvals)...)
}

// With 返回附加了指定 key-value 对的新 Logger 实例。
// 使用 zap 原生的 With() 机制。
func (l *Logger) With(keyvals ...any) bLogger.Logger {
	return &Logger{log: l.log.With(l.toFields(keyvals)...)}
}

// toFields 将交替的 key-value 对转换为 zap.Field 切片。
func (l *Logger) toFields(keyvals []any) []zap.Field {
	fields := make([]zap.Field, 0, len(keyvals)/2)
	for i := 0; i+1 < len(keyvals); i += 2 {
		fields = append(fields, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}
	return fields
}

func (l *Logger) Sync() error {
	return l.log.Sync()
}

func (l *Logger) Close() error {
	return l.Sync()
}
