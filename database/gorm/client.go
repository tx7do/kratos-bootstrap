package gorm

import (
	"gorm.io/driver/bigquery"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"

	"gorm.io/plugin/opentelemetry/tracing"
	"gorm.io/plugin/prometheus"

	"gorm.io/gorm"

	"github.com/go-kratos/kratos/v2/log"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewGormClient 创建GORM数据库客户端
func NewGormClient(cfg *conf.Bootstrap, l *log.Helper, migrates []interface{}) *gorm.DB {
	if cfg.Data == nil || cfg.Data.Database == nil {
		l.Warn("database config is nil")
		return nil
	}

	var driver gorm.Dialector
	switch cfg.Data.Database.Driver {
	default:
		fallthrough
	case "mysql":
		driver = mysql.Open(cfg.Data.Database.Source)
		break
	case "postgres":
		driver = postgres.Open(cfg.Data.Database.Source)
		break
	case "clickhouse":
		driver = clickhouse.Open(cfg.Data.Database.Source)
		break
	case "sqlite":
		driver = sqlite.Open(cfg.Data.Database.Source)
		break
	case "sqlserver":
		driver = sqlserver.Open(cfg.Data.Database.Source)
		break
	case "bigquery":
		driver = bigquery.Open(cfg.Data.Database.Source)
		break
	}

	db, err := gorm.Open(driver, &gorm.Config{})
	if err != nil {
		l.Fatalf("failed opening connection to db: %v", err)
		return nil
	}

	// 运行数据库迁移工具
	if cfg.Data.Database.Migrate {
		if err = db.AutoMigrate(
			migrates...,
		); err != nil {
			l.Fatalf("failed creating schema resources: %v", err)
			return nil
		}
	}

	if cfg.Data.Database.GetEnableTrace() {
		if err = db.Use(tracing.NewPlugin()); err != nil {
			l.Fatalf("failed enable trace: %v", err)
			return nil
		}
	}

	if cfg.Data.Database.GetEnableMetrics() {
		db.Use(prometheus.New(prometheus.Config{
			DBName:          "db1",                       // `DBName` as metrics label
			RefreshInterval: 15,                          // refresh metrics interval (default 15 seconds)
			PushAddr:        "prometheus pusher address", // push metrics if `PushAddr` configured
			StartServer:     true,                        // start http server to expose metrics
			HTTPServerPort:  8080,                        // configure http server port, default port 8080 (if you have configured multiple instances, only the first `HTTPServerPort` will be used to start server)
			MetricsCollector: []prometheus.MetricsCollector{
				&prometheus.MySQL{VariableNames: []string{"Threads_running"}},
			},
			Labels: map[string]string{
				"instance": "127.0.0.1", // config custom labels if necessary
			},
		}))

	}

	sqlDB, err := db.DB()
	if sqlDB != nil {
		if cfg.Data.Database.MaxIdleConnections != nil {
			sqlDB.SetMaxIdleConns(int(cfg.Data.Database.GetMaxIdleConnections()))
		}
		if cfg.Data.Database.MaxOpenConnections != nil {
			sqlDB.SetMaxOpenConns(int(cfg.Data.Database.GetMaxOpenConnections()))
		}
		if cfg.Data.Database.ConnectionMaxLifetime != nil {
			sqlDB.SetConnMaxLifetime(cfg.Data.Database.GetConnectionMaxLifetime().AsDuration())
		}
	}

	return db
}
