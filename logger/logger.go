package logger

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/go-kratos/kratos/v2/log"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewLogger 动态创建日志实例，返回项目 Logger 接口。
//
// 内部通过工厂注册表查找后端（zap/logrus/fluent 等）。
// 工厂现在直接返回项目 Logger 接口，无需 Kratos 适配。
func NewLogger(cfg *conf.Logger) (Logger, error) {
	if cfg == nil {
		return nil, nil
	}

	if cfg.GetType() == "" || cfg.GetType() == string(Std) {
		return NewStdLogger(), nil
	}

	// normalize to lower case for lookup
	typ := Type(strings.ToLower(cfg.GetType()))

	f, ok := GetFactory(typ)
	if !ok {
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

// NewKratosLogger 动态创建日志实例，返回 Kratos log.Logger。
// 旧 API 兼容入口，内部通过 AsKratosLogger 适配项目 Logger。
// 新代码请使用 [NewLogger]。
func NewKratosLogger(cfg *conf.Logger) (log.Logger, error) {
	lg, err := NewLogger(cfg)
	if err != nil {
		return nil, err
	}
	if lg == nil {
		return nil, nil
	}
	return AsKratosLogger(lg), nil
}

// NewLoggerProvider 创建一个新的日志记录器提供者，返回项目 Logger 接口。
// 它会从 cfg 创建具体 logger（通过 NewLogger），并为 logger 附加一组标准字段
// （service.*, ts, caller, trace_id, span_id）。
// 实现是防御性的：当 cfg 或 appInfo 为空或 NewLogger 返回 nil/err 时，
// 会回退到标准控制台 logger。
func NewLoggerProvider(cfg *conf.Logger, appInfo *conf.AppInfo) Logger {
	var l Logger
	if cfg == nil || cfg.GetType() == "" {
		l = NewStdLogger()
	} else {
		if lg, err := NewLogger(cfg); err == nil && lg != nil {
			l = lg
		} else {
			l = NewStdLogger()
		}
	}

	// 通过 Kratos log.With 附加 Valuer 字段（ts/caller）。
	// 这些 Valuer 不依赖 ctx（DefaultTimestamp 返回当前时间，DefaultCaller 返回调用栈），
	// 所以能正确穿过 Kratos log.Logger 不携带 ctx 的限制。
	// trace_id/span_id 不使用 Valuer（因 Kratos Log 不接收 ctx 导致求值为空），
	// 而是通过 WrapTrace 从 Logger 接口的 ctx 参数直接提取。
	kl := AsKratosLogger(l)

	fields := []any{
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
	}

	if appInfo != nil {
		fields = append([]any{
			"service.id", appInfo.GetAppId(),
			"service.instance", appInfo.GetInstanceId(),
			"service.version", appInfo.GetVersion(),
		}, fields...)
	}

	kl = log.With(kl, fields...)
	l = FromKratosLogger(kl)

	// 自动注入 trace_id / span_id（从 Logger 接口的 ctx 参数提取）
	return WrapTrace(l)
}

// NewKratosStdLogger 创建一个 Kratos 内置的控制台日志记录器。
// 旧 API 兼容入口，新代码请使用 [NewStdLogger]（slog 后端，无 Kratos 依赖）。
func NewKratosStdLogger() log.Logger {
	return log.NewStdLogger(os.Stdout)
}
