package tencent

import (
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	bLogger "github.com/tx7do/kratos-bootstrap/logger"
)

func init() {
	_ = bLogger.Register(bLogger.Tencent, func(cfg *conf.Logger) (bLogger.Logger, error) {
		return NewLogger(cfg)
	})
}

// NewLogger 创建一个新的日志记录器 - Tencent
func NewLogger(cfg *conf.Logger) (bLogger.Logger, error) {
	if cfg == nil || cfg.Tencent == nil {
		return nil, nil
	}

	wrapped, err := NewTencentLogger(
		WithTopicID(cfg.Tencent.TopicId),
		WithEndpoint(cfg.Tencent.Endpoint),
		WithAccessKey(cfg.Tencent.AccessKey),
		WithAccessSecret(cfg.Tencent.AccessSecret),
	)
	if err != nil {
		return nil, err
	}
	return wrapped, nil
}
