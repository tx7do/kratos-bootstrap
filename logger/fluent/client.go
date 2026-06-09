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

	wrapped, err := NewFluentLogger(cfg.Fluent.Endpoint)
	if err != nil {
		return nil, err
	}
	return wrapped, nil
}
