package logger

import (
	"os"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	"github.com/tx7do/kratos-bootstrap/logger/fluent"
	"github.com/tx7do/kratos-bootstrap/logger/logrus"
	"github.com/tx7do/kratos-bootstrap/logger/tencent"
	"github.com/tx7do/kratos-bootstrap/logger/zap"

	"github.com/tx7do/kratos-bootstrap/config"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewLogger 创建一个新的日志记录器
func NewLogger(cfg *conf.Logger) log.Logger {
	if cfg == nil {
		return NewStdLogger()
	}

	switch Type(cfg.Type) {
	default:
		fallthrough
	case Std:
		return NewStdLogger()
	case Fluent:
		return fluent.NewLogger(cfg)
	case Zap:
		return zap.NewLogger(cfg)
	case Logrus:
		return logrus.NewLogger(cfg)
	case Aliyun:
		return nil
	case Tencent:
		return tencent.NewLogger(cfg)
	}
}

// NewLoggerProvider 创建一个新的日志记录器提供者
func NewLoggerProvider(cfg *conf.Logger, serviceInfo *config.ServiceInfo) log.Logger {
	l := NewLogger(cfg)

	return log.With(
		l,
		"service.id", serviceInfo.Id,
		"service.name", serviceInfo.Name,
		"service.version", serviceInfo.Version,
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
	)
}

// NewStdLogger 创建一个新的日志记录器 - Kratos内置，控制台输出
func NewStdLogger() log.Logger {
	l := log.NewStdLogger(os.Stdout)
	return l
}
