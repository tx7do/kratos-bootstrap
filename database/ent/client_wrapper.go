package ent

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"

	"github.com/XSAM/otelsql"

	entSql "entgo.io/ent/dialect/sql"
)

type ClientInterface interface {
	Close() error
}

type ClientWrapper[T ClientInterface] struct {
	db  T
	drv *entSql.Driver
}

func NewEntClientWrapper[T ClientInterface](db T, drv *entSql.Driver) *ClientWrapper[T] {
	return &ClientWrapper[T]{
		db:  db,
		drv: drv,
	}
}

func (c *ClientWrapper[T]) Client() T {
	return c.db
}

func (c *ClientWrapper[T]) Driver() *entSql.Driver {
	return c.drv
}

func (c *ClientWrapper[T]) DB() *sql.DB {
	return c.drv.DB()
}

// Close 关闭数据库连接
func (c *ClientWrapper[T]) Close() error {
	return c.db.Close()
}

// Query 查询数据
func (c *ClientWrapper[T]) Query(ctx context.Context, query string, args, v any) error {
	return c.Driver().Query(ctx, query, args, v)
}

func (c *ClientWrapper[T]) Exec(ctx context.Context, query string, args, v any) error {
	return c.Driver().Exec(ctx, query, args, v)
}

// SetConnectionOption 设置连接配置
func (c *ClientWrapper[T]) SetConnectionOption(maxIdleConnections, maxOpenConnections int, connMaxLifetime time.Duration) {
	// 连接池中最多保留的空闲连接数量
	c.DB().SetMaxIdleConns(maxIdleConnections)
	// 连接池在同一时间打开连接的最大数量
	c.DB().SetMaxOpenConns(maxOpenConnections)
	// 连接可重用的最大时间长度
	c.DB().SetConnMaxLifetime(connMaxLifetime)
}

func driverNameToSemConvKeyValue(driverName string) attribute.KeyValue {
	switch driverName {
	case "mariadb":
		return semconv.DBSystemMariaDB
	case "mysql":
		return semconv.DBSystemMySQL
	case "postgresql":
		return semconv.DBSystemPostgreSQL
	case "sqlite":
		return semconv.DBSystemSqlite
	default:
		return semconv.DBSystemKey.String(driverName)
	}
}

// CreateDriver 创建数据库驱动
func CreateDriver(driverName, dsn string, enableTrace, enableMetrics bool) (*entSql.Driver, error) {
	var db *sql.DB
	var drv *entSql.Driver
	var err error

	if enableTrace {
		// Connect to database
		if db, err = otelsql.Open(driverName, dsn, otelsql.WithAttributes(
			driverNameToSemConvKeyValue(driverName),
		)); err != nil {
			return nil, errors.New(fmt.Sprintf("failed opening connection to db: %v", err))
		}

		drv = entSql.OpenDB(driverName, db)
	} else {
		if drv, err = entSql.Open(driverName, dsn); err != nil {
			return nil, errors.New(fmt.Sprintf("failed opening connection to db: %v", err))
		}

		db = drv.DB()
	}

	// Register DB stats to meter
	if enableMetrics {
		err = otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(
			driverNameToSemConvKeyValue(driverName),
		))
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed register otel meter: %v", err))
		}
	}

	return drv, nil
}
