package fluent

import (
	"github.com/go-kratos/kratos/v2/log"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/logger"
)

func init() {
	_ = logger.Register(logger.Fluent, func(cfg *conf.Logger) (log.Logger, error) {
		return NewLogger(cfg)
	})
}

// NewLogger 创建一个新的日志记录器 - Fluent
func NewLogger(cfg *conf.Logger) (log.Logger, error) {
	if cfg == nil || cfg.Fluent == nil {
		return nil, nil
	}

	wrapped, err := NewFluentLogger(cfg.Fluent.Endpoint)
	if err != nil {
		return nil, err
	}
	return wrapped, nil
}
