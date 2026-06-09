package zerolog

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	bLogger "github.com/tx7do/kratos-bootstrap/logger"
)

// Compile-time assertion: Logger implements bLogger.Logger.
var _ bLogger.Logger = (*Logger)(nil)

type Logger struct {
	log *zerolog.Logger
}

func NewZerologLogger(logger *zerolog.Logger) *Logger {
	return &Logger{log: logger}
}

// Debug 输出 DEBUG 级别日志。
func (l *Logger) Debug(_ context.Context, msg string, keyvals ...any) {
	l.logEvent(l.log.Debug(), msg, keyvals)
}

// Info 输出 INFO 级别日志。
func (l *Logger) Info(_ context.Context, msg string, keyvals ...any) {
	l.logEvent(l.log.Info(), msg, keyvals)
}

// Warn 输出 WARN 级别日志。
func (l *Logger) Warn(_ context.Context, msg string, keyvals ...any) {
	l.logEvent(l.log.Warn(), msg, keyvals)
}

// Error 输出 ERROR 级别日志。
func (l *Logger) Error(_ context.Context, msg string, keyvals ...any) {
	l.logEvent(l.log.Error(), msg, keyvals)
}

// With 返回附加了指定 key-value 对的新 Logger 实例。
func (l *Logger) With(keyvals ...any) bLogger.Logger {
	ctx := l.log.With()
	for i := 0; i+1 < len(keyvals); i += 2 {
		ctx = ctx.Any(fmt.Sprint(keyvals[i]), keyvals[i+1])
	}
	newLog := ctx.Logger()
	return &Logger{log: &newLog}
}

// logEvent 统一处理 zerolog event 的字段附加和消息发送。
func (l *Logger) logEvent(event *zerolog.Event, msg string, keyvals []any) {
	for i := 0; i+1 < len(keyvals); i += 2 {
		event = event.Any(fmt.Sprint(keyvals[i]), keyvals[i+1])
	}
	event.Msg(msg)
}
