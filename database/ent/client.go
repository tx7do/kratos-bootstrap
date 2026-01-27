package ent

import (
	"entgo.io/ent/dialect/sql"

	"github.com/go-kratos/kratos/v2/log"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"

	entCrud "github.com/tx7do/go-crud/entgo"
)

// DbCreator 定义创建Ent ORM数据库客户端的函数类型
type DbCreator[T entCrud.EntClientInterface] func(drv *sql.Driver) T

// NewEntClient 创建Ent ORM数据库客户端
func NewEntClient[T entCrud.EntClientInterface](cfg *conf.Bootstrap, dbCreator DbCreator[T]) *entCrud.EntClient[T] {
	if cfg.Data == nil || cfg.Data.Database == nil {
		log.Warn("[ENT] database config is nil")
		return nil
	}

	if dbCreator == nil {
		log.Warn("[ENT] dbCreator is nil")
		return nil
	}

	drv, err := entCrud.CreateDriver(
		cfg.Data.Database.GetDriver(),
		cfg.Data.Database.GetSource(),
		cfg.Data.Database.GetEnableTrace(),
		cfg.Data.Database.GetEnableMetrics(),
	)
	if err != nil {
		log.Fatalf("[ENT] failed opening connection to db: %v", err)
		return nil
	}

	db := dbCreator(drv)

	wrapperClient := entCrud.NewEntClient(db, drv)
	if wrapperClient == nil {
		log.Fatalf("[ENT] failed creating ent client")
		return nil
	}

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
