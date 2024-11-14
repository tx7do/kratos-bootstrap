package aliyun

import (
	aliyunLogger "github.com/go-kratos/kratos/contrib/log/aliyun/v2"
	"github.com/go-kratos/kratos/v2/log"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewLogger 创建一个新的日志记录器 - Aliyun
func NewLogger(cfg *conf.Logger) log.Logger {
	if cfg == nil || cfg.Aliyun == nil {
		return nil
	}

	wrapped := aliyunLogger.NewAliyunLog(
		aliyunLogger.WithProject(cfg.Aliyun.Project),
		aliyunLogger.WithEndpoint(cfg.Aliyun.Endpoint),
		aliyunLogger.WithAccessKey(cfg.Aliyun.AccessKey),
		aliyunLogger.WithAccessSecret(cfg.Aliyun.AccessSecret),
	)
	return wrapped
}
