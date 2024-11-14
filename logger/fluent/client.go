package fluent

import (
	fluentLogger "github.com/go-kratos/kratos/contrib/log/fluent/v2"
	"github.com/go-kratos/kratos/v2/log"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewLogger 创建一个新的日志记录器 - Fluent
func NewLogger(cfg *conf.Logger) log.Logger {
	if cfg == nil || cfg.Fluent == nil {
		return nil
	}

	wrapped, err := fluentLogger.NewLogger(cfg.Fluent.Endpoint)
	if err != nil {
		panic("create fluent logger failed")
		return nil
	}
	return wrapped
}
