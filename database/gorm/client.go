package gorm

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	bLogger "github.com/tx7do/kratos-bootstrap/logger"
	gwLog "github.com/tx7do/go-wind/log"

	gormCrud "github.com/tx7do/go-crud/gorm"
)

// NewGormClient 创建GORM数据库客户端
// l 为 bootstrap.Context.NewLoggerHelper 产出的 *logger.Helper
// (bootstrap v0.1.17 起为 tx7do logger.Helper,不再是 kratos log.Helper)。
func NewGormClient(cfg *conf.Bootstrap, l *bLogger.Helper, migrates []any) (*gormCrud.Client, error) {
	if cfg == nil || cfg.Data == nil || cfg.Data.Database == nil {
		log.Warn("[GORM] database config is nil")
		return nil, nil
	}

	var options []gormCrud.Option

	if l != nil {
		options = append(options, gormCrud.WithLogger(loggerBridge{h: l}))
	}

	if len(migrates) > 0 {
		options = append(options, gormCrud.WithAutoMigrate(migrates...))
	}

	if cfg.Data.Database.GetDriver() != "" {
		options = append(options, gormCrud.WithDriverName(cfg.Data.Database.GetDriver()))
	}
	if cfg.Data.Database.GetSource() != "" {
		options = append(options, gormCrud.WithDSN(cfg.Data.Database.GetSource()))
	}

	options = append(options, gormCrud.WithEnableMigrate(cfg.Data.Database.GetMigrate()))
	options = append(options, gormCrud.WithEnableTrace(cfg.Data.Database.GetEnableTrace()))
	options = append(options, gormCrud.WithEnableMetrics(cfg.Data.Database.GetEnableMetrics()))
	//options = append(options, gormCrud.WithEnableDbResolver(cfg.Data.Database.GetEnableDbResolver()))

	if cfg.Data.Database.MaxIdleConnections != nil {
		options = append(options, gormCrud.WithMaxIdleConns(int(cfg.Data.Database.GetMaxIdleConnections())))
	}
	if cfg.Data.Database.MaxOpenConnections != nil {
		options = append(options, gormCrud.WithMaxOpenConns(int(cfg.Data.Database.GetMaxOpenConnections())))
	}
	if cfg.Data.Database.ConnectionMaxLifetime != nil {
		options = append(options, gormCrud.WithConnMaxLifetime(cfg.Data.Database.GetConnectionMaxLifetime().AsDuration()))
	}

	db, err := gormCrud.NewClient(options...)
	if err != nil {
		log.Errorf("[GORM] failed opening connection to db: %v", err)
		return nil, err
	}

	return db, nil
}

// loggerBridge 把 bootstrap logger.Helper 适配为 go-wind/log.Logger——
// 方法签名 1:1 转发(bLogger 四方法即 ctx 首参格式化);With 在两侧
// 均无 Helper 语义,返回自身。
type loggerBridge struct {
	h *bLogger.Helper
}

func (b loggerBridge) Debug(ctx context.Context, msg string, args ...any) {
	b.h.Debugf(ctx, msg, args...)
}

func (b loggerBridge) Info(ctx context.Context, msg string, args ...any) {
	b.h.Infof(ctx, msg, args...)
}

func (b loggerBridge) Warn(ctx context.Context, msg string, args ...any) {
	b.h.Warnf(ctx, msg, args...)
}

func (b loggerBridge) Error(ctx context.Context, msg string, args ...any) {
	b.h.Errorf(ctx, msg, args...)
}

// Enabled 恒真:级别过滤交给底层 logger,这里不做二次门控。
func (b loggerBridge) Enabled(_ gwLog.Level) bool {
	return true
}

func (b loggerBridge) With(args ...any) gwLog.Logger {
	return b
}
