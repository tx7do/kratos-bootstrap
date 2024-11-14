package clickhouse

import (
	"crypto/tls"

	"github.com/ClickHouse/clickhouse-go/v2"

	"github.com/go-kratos/kratos/v2/log"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/utils"
)

func NewClickHouseClient(cfg *conf.Bootstrap, l *log.Helper) clickhouse.Conn {
	if cfg.Data == nil || cfg.Data.Clickhouse == nil {
		l.Warn("ClickHouse config is nil")
		return nil
	}

	options := &clickhouse.Options{
		Addr: []string{cfg.Data.Clickhouse.Address},
		Auth: clickhouse.Auth{
			Database: cfg.Data.Clickhouse.Database,
			Username: cfg.Data.Clickhouse.Username,
			Password: cfg.Data.Clickhouse.Password,
		},
		Debug:           cfg.Data.Clickhouse.Debug,
		DialTimeout:     cfg.Data.Clickhouse.DialTimeout.AsDuration(),
		MaxOpenConns:    int(cfg.Data.Clickhouse.MaxOpenConns),
		MaxIdleConns:    int(cfg.Data.Clickhouse.MaxIdleConns),
		ConnMaxLifetime: cfg.Data.Clickhouse.ConnMaxLifeTime.AsDuration(),
	}

	// 设置ssl
	if cfg.Data.Clickhouse.Tls != nil {
		var tlsCfg *tls.Config
		var err error

		if cfg.Data.Clickhouse.Tls.File != nil {
			if tlsCfg, err = utils.LoadServerTlsConfigFile(
				cfg.Data.Clickhouse.Tls.File.GetKeyPath(),
				cfg.Data.Clickhouse.Tls.File.GetCertPath(),
				cfg.Data.Clickhouse.Tls.File.GetCaPath(),
				cfg.Data.Clickhouse.Tls.InsecureSkipVerify,
			); err != nil {
				panic(err)
			}
		}
		if tlsCfg == nil && cfg.Data.Clickhouse.Tls.Config != nil {
			if tlsCfg, err = utils.LoadServerTlsConfig(
				cfg.Data.Clickhouse.Tls.Config.GetKeyPem(),
				cfg.Data.Clickhouse.Tls.Config.GetCertPem(),
				cfg.Data.Clickhouse.Tls.Config.GetCaPem(),
				cfg.Data.Clickhouse.Tls.InsecureSkipVerify,
			); err != nil {
				panic(err)
			}
		}

		if tlsCfg != nil {
			options.TLS = tlsCfg
		}
	}

	conn, err := clickhouse.Open(options)
	if err != nil {
		l.Fatalf("failed opening connection to clickhouse: %v", err)
		return nil
	}

	return conn
}
