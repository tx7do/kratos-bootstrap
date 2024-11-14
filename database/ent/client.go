package ent

import (
	_ "github.com/go-sql-driver/mysql"

	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"

	"github.com/go-kratos/kratos/v2/log"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewEntClient 创建Ent ORM数据库客户端
func NewEntClient[T ClientInterface](cfg *conf.Bootstrap, l *log.Helper, db T) *ClientWrapper[T] {
	if cfg.Data == nil || cfg.Data.Database == nil {
		l.Warn("database config is nil")
		return nil
	}

	drv, err := CreateDriver(
		cfg.Data.Database.GetDriver(),
		cfg.Data.Database.GetSource(),
		cfg.Data.Database.GetEnableTrace(),
		cfg.Data.Database.GetEnableMetrics(),
	)
	if err != nil {
		l.Fatalf("failed opening connection to db: %v", err)
		return nil
	}

	wrapperClient := NewEntClientWrapper(db, drv)

	wrapperClient.SetConnectionOption(
		int(cfg.Data.Database.GetMaxIdleConnections()),
		int(cfg.Data.Database.GetMaxOpenConnections()),
		cfg.Data.Database.GetConnectionMaxLifetime().AsDuration(),
	)

	return wrapperClient
}
