package gorm

import (
	"github.com/go-kratos/kratos/v2/log"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"

	gormCrud "github.com/tx7do/go-crud/gorm"
)

// NewGormClient 创建GORM数据库客户端
func NewGormClient(cfg *conf.Bootstrap, l *log.Helper, migrates []interface{}) *gormCrud.Client {
	if cfg.Data == nil || cfg.Data.Database == nil {
		l.Warn("database config is nil")
		return nil
	}

	var options []gormCrud.Option

	if l != nil {
		options = append(options, gormCrud.WithLogger(l))
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
		l.Fatalf("failed opening connection to db: %v", err)
		return nil
	}

	return db
}
