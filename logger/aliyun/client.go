package aliyun

import (
	"errors"

	"github.com/go-kratos/kratos/v2/log"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/logger"
)

func init() {
	_ = logger.Register(logger.Aliyun, func(cfg *conf.Logger) (log.Logger, error) {
		return NewLogger(cfg)
	})
}

// NewLogger 创建一个新的日志记录器 - Aliyun
func NewLogger(cfg *conf.Logger) (log.Logger, error) {
	if cfg == nil || cfg.Aliyun == nil {
		return nil, nil
	}

	// basic validation of required fields
	if cfg.Aliyun.Project == "" || cfg.Aliyun.Endpoint == "" || cfg.Aliyun.AccessKey == "" || cfg.Aliyun.AccessSecret == "" {
		return nil, errors.New("aliyun config invalid")
	}

	wrapped, err := NewAliyunLog(
		WithProject(cfg.Aliyun.Project),
		WithEndpoint(cfg.Aliyun.Endpoint),
		WithAccessKey(cfg.Aliyun.AccessKey),
		WithAccessSecret(cfg.Aliyun.AccessSecret),
	)
	if err != nil {
		// creation failed, return nil so caller can fallback
		return nil, err
	}

	return wrapped, nil
}
