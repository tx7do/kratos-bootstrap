package ent

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"

	"github.com/go-kratos/kratos/v2/log"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"

	entCrud "github.com/tx7do/go-crud/entgo"
)

// NewEntClient 创建Ent ORM数据库客户端
func NewEntClient[T entCrud.EntClientInterface](cfg *conf.Bootstrap, l *log.Helper, db T) *entCrud.EntClient[T] {
	if cfg.Data == nil || cfg.Data.Database == nil {
		l.Warn("database config is nil")
		return nil
	}

	drv, err := entCrud.CreateDriver(
		cfg.Data.Database.GetDriver(),
		cfg.Data.Database.GetSource(),
		cfg.Data.Database.GetEnableTrace(),
		cfg.Data.Database.GetEnableMetrics(),
	)
	if err != nil {
		l.Fatalf("failed opening connection to db: %v", err)
		return nil
	}

	wrapperClient := entCrud.NewEntClient(db, drv)

	if cfg.Data.Database.MaxIdleConnections != nil &&
		cfg.Data.Database.MaxOpenConnections != nil &&
		cfg.Data.Database.ConnectionMaxLifetime != nil {
		wrapperClient.SetConnectionOption(
			int(cfg.Data.Database.GetMaxIdleConnections()),
			int(cfg.Data.Database.GetMaxOpenConnections()),
			cfg.Data.Database.GetConnectionMaxLifetime().AsDuration(),
		)
	}

	return wrapperClient
}
