package logger

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewLogger 动态创建日志实例
func NewLogger(cfg *conf.Logger) (log.Logger, error) {
	if cfg == nil {
		return nil, nil
	}

	if cfg.GetType() == "" || cfg.GetType() == string(Std) {
		return NewStdLogger(), nil
	}

	// normalize to lower case for lookup
	typ := Type(strings.ToLower(cfg.GetType()))
	norm := Type(strings.ToLower(string(typ)))

	f, ok := GetFactory(norm)
	if !ok {
		// prepare available list for helpful error
		available := ListFactories()
		strs := make([]string, 0, len(available))
		for _, t := range available {
			strs = append(strs, string(t))
		}
		sort.Strings(strs)
		return nil, fmt.Errorf("unsupported logger type: %s; available: %v", typ, strs)
	}

	lg, err := f(cfg)
	if err != nil {
		return nil, fmt.Errorf("create logger %s: %w", typ, err)
	}
	if lg == nil {
		return nil, fmt.Errorf("logger factory %s returned nil logger", typ)
	}
	return lg, nil
}

// NewLoggerProvider 创建一个新的日志记录器提供者
// 它会从 cfg 创建具体 logger（通过 NewLogger），并为 logger 附加一组标准字段（service.*, ts, caller, trace_id, span_id）。
// 实现是防御性的：当 cfg 或 appInfo 为空或 NewLogger 返回 nil/err 时，会回退到标准控制台 logger。
func NewLoggerProvider(cfg *conf.Logger, appInfo *conf.AppInfo) log.Logger {
	var l log.Logger
	if cfg == nil || cfg.GetType() == "" {
		l = NewStdLogger()
	} else {
		// try to create logger by type via factory
		if lg, err := NewLogger(cfg); err == nil && lg != nil {
			l = lg
		} else {
			l = NewStdLogger()
		}
	}

	// build base fields - always include timestamp, caller, trace/span ids
	fields := []interface{}{
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
	}

	// attach service fields only if appInfo is provided
	if appInfo != nil {
		fields = append([]interface{}{
			"service.id", appInfo.GetAppId(),
			"service.name", appInfo.GetName(),
			"service.version", appInfo.GetVersion(),
		}, fields...)
	}

	return log.With(l, fields...)
}

// NewStdLogger 创建一个新的日志记录器 - Kratos内置，控制台输出
func NewStdLogger() log.Logger {
	l := log.NewStdLogger(os.Stdout)
	return l
}
