package fluent

import (
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	bLogger "github.com/tx7do/kratos-bootstrap/logger"
)

func init() {
	_ = bLogger.Register(bLogger.Fluent, func(cfg *conf.Logger) (bLogger.Logger, error) {
		return NewLogger(cfg)
	})
}

// NewLogger 创建一个新的日志记录器 - Fluent
func NewLogger(cfg *conf.Logger) (bLogger.Logger, error) {
	if cfg == nil || cfg.Fluent == nil {
		return nil, nil
	}

	var opts []Option
	in := cfg.Fluent

	if in.WriteTimeout != nil {
		opts = append(opts, WithWriteTimeout(in.GetWriteTimeout().AsDuration()))
	}
	if in.BufferLimit != nil {
		opts = append(opts, WithBufferLimit(int(in.GetBufferLimit())))
	}
	if in.RetryWait != nil {
		opts = append(opts, WithRetryWait(int(in.GetRetryWait().AsDuration().Milliseconds())))
	}
	if in.MaxRetry != nil {
		opts = append(opts, WithMaxRetry(int(in.GetMaxRetry())))
	}
	if in.MaxRetryWait != nil {
		opts = append(opts, WithMaxRetryWait(int(in.GetMaxRetryWait().AsDuration().Milliseconds())))
	}
	if in.TagPrefix != nil && in.GetTagPrefix() != "" {
		opts = append(opts, WithTagPrefix(in.GetTagPrefix()))
	}
	if in.Async != nil {
		opts = append(opts, WithAsync(in.GetAsync()))
	}
	if in.ForceStopAsyncSend != nil {
		opts = append(opts, WithForceStopAsyncSend(in.GetForceStopAsyncSend()))
	}

	wrapped, err := NewFluentLogger(in.Endpoint, opts...)
	if err != nil {
		return nil, err
	}
	return wrapped, nil
}
