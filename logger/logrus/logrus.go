package logrus

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	bLogger "github.com/tx7do/kratos-bootstrap/logger"
)

// Compile-time assertion: Logger implements bLogger.Logger.
var _ bLogger.Logger = (*Logger)(nil)

type Logger struct {
	log    *logrus.Logger
	fields logrus.Fields
}

func NewLogrusLogger(logger *logrus.Logger) *Logger {
	return &Logger{log: logger, fields: make(logrus.Fields)}
}

// Debug 输出 DEBUG 级别日志。
func (l *Logger) Debug(_ context.Context, msg string, keyvals ...any) {
	l.logEntry(logrus.DebugLevel, msg, keyvals)
}

// Info 输出 INFO 级别日志。
func (l *Logger) Info(_ context.Context, msg string, keyvals ...any) {
	l.logEntry(logrus.InfoLevel, msg, keyvals)
}

// Warn 输出 WARN 级别日志。
func (l *Logger) Warn(_ context.Context, msg string, keyvals ...any) {
	l.logEntry(logrus.WarnLevel, msg, keyvals)
}

// Error 输出 ERROR 级别日志。
func (l *Logger) Error(_ context.Context, msg string, keyvals ...any) {
	l.logEntry(logrus.ErrorLevel, msg, keyvals)
}

// With 返回附加了指定 key-value 对的新 Logger 实例。
func (l *Logger) With(keyvals ...any) bLogger.Logger {
	newFields := make(logrus.Fields)
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range toFields(keyvals) {
		newFields[k] = v
	}
	return &Logger{log: l.log, fields: newFields}
}

// logEntry 统一处理 logrus 日志输出。
func (l *Logger) logEntry(level logrus.Level, msg string, keyvals []any) {
	if level > l.log.Level {
		return
	}
	entry := l.log.WithFields(l.fields)
	extra := toFields(keyvals)
	if len(extra) > 0 {
		entry = entry.WithFields(extra)
	}
	entry.Log(level, msg)
}

// toFields 将交替的 key-value 对转换为 logrus.Fields。
func toFields(keyvals []any) logrus.Fields {
	fields := make(logrus.Fields)
	for i := 0; i+1 < len(keyvals); i += 2 {
		fields[fmt.Sprint(keyvals[i])] = keyvals[i+1]
	}
	return fields
}
