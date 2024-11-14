package tencent

import (
	tencentLogger "github.com/go-kratos/kratos/contrib/log/tencent/v2"
	"github.com/go-kratos/kratos/v2/log"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewLogger 创建一个新的日志记录器 - Tencent
func NewLogger(cfg *conf.Logger) log.Logger {
	if cfg == nil || cfg.Tencent == nil {
		return nil
	}

	wrapped, err := tencentLogger.NewLogger(
		tencentLogger.WithTopicID(cfg.Tencent.TopicId),
		tencentLogger.WithEndpoint(cfg.Tencent.Endpoint),
		tencentLogger.WithAccessKey(cfg.Tencent.AccessKey),
		tencentLogger.WithAccessSecret(cfg.Tencent.AccessSecret),
	)
	if err != nil {
		panic(err)
		return nil
	}
	return wrapped
}
