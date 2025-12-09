package clickhouse

import (
	"crypto/tls"
	"database/sql"
	"errors"

	clickhouseV2 "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/go-kratos/kratos/v2/log"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/utils"

	clickhouseCrud "github.com/tx7do/go-crud/clickhouse"
)

type Client struct {
	log *log.Helper

	conn clickhouseV2.Conn
	db   *sql.DB
}

func NewClient(logger log.Logger, cfg *conf.Bootstrap) (*clickhouseCrud.Client, error) {
	if cfg.Data == nil || cfg.Data.Clickhouse == nil {
		return nil, errors.New("clickhouse config is nil")
	}

	var options []clickhouseCrud.Option

	if logger != nil {
		options = append(options, clickhouseCrud.WithLogger(logger))
	}
	if cfg.Data.Clickhouse.Dsn != nil {
		options = append(options, clickhouseCrud.WithDsn(cfg.Data.Clickhouse.GetDsn()))
	}
	if cfg.Data.Clickhouse.Addresses != nil {
		options = append(options, clickhouseCrud.WithAddresses(cfg.Data.Clickhouse.GetAddresses()...))
	}
	if cfg.Data.Clickhouse.Database != nil {
		options = append(options, clickhouseCrud.WithDatabase(cfg.Data.Clickhouse.GetDatabase()))
	}
	if cfg.Data.Clickhouse.Username != nil {
		options = append(options, clickhouseCrud.WithUsername(cfg.Data.Clickhouse.GetUsername()))
	}
	if cfg.Data.Clickhouse.Password != nil {
		options = append(options, clickhouseCrud.WithPassword(cfg.Data.Clickhouse.GetPassword()))
	}
	if cfg.Data.Clickhouse.Debug != nil {
		options = append(options, clickhouseCrud.WithDebug(cfg.Data.Clickhouse.GetDebug()))
	}
	if cfg.Data.Clickhouse.MaxOpenConns != nil {
		options = append(options, clickhouseCrud.WithMaxOpenConns(int(cfg.Data.Clickhouse.GetMaxOpenConns())))
	}
	if cfg.Data.Clickhouse.MaxIdleConns != nil {
		options = append(options, clickhouseCrud.WithMaxIdleConns(int(cfg.Data.Clickhouse.GetMaxIdleConns())))
	}

	if cfg.Data.Clickhouse.Tls != nil {
		var tlsCfg *tls.Config
		var err error

		if tlsCfg, err = utils.LoadServerTlsConfig(cfg.Server.Grpc.Tls); err != nil {
			panic(err)
		}

		if tlsCfg != nil {
			options = append(options, clickhouseCrud.WithTLSConfig(tlsCfg))
		}
	}

	if cfg.Data.Clickhouse.GetCompressionMethod() != "" {
		options = append(options, clickhouseCrud.WithCompressionMethod(cfg.Data.Clickhouse.GetCompressionMethod()))
	}
	if cfg.Data.Clickhouse.CompressionLevel != nil {
		options = append(options, clickhouseCrud.WithCompressionLevel(int(cfg.Data.Clickhouse.GetCompressionLevel())))
	}

	if cfg.Data.Clickhouse.MaxCompressionBuffer != nil {
		options = append(options, clickhouseCrud.WithMaxCompressionBuffer(int(cfg.Data.Clickhouse.GetMaxCompressionBuffer())))
	}

	if cfg.Data.Clickhouse.DialTimeout != nil {
		options = append(options, clickhouseCrud.WithDialTimeout(cfg.Data.Clickhouse.GetDialTimeout().AsDuration()))
	}
	if cfg.Data.Clickhouse.ReadTimeout != nil {
		options = append(options, clickhouseCrud.WithReadTimeout(cfg.Data.Clickhouse.GetReadTimeout().AsDuration()))
	}
	if cfg.Data.Clickhouse.ConnMaxLifetime != nil {
		options = append(options, clickhouseCrud.WithConnMaxLifetime(cfg.Data.Clickhouse.GetConnMaxLifetime().AsDuration()))
	}

	if cfg.Data.Clickhouse.HttpProxy != nil {
		options = append(options, clickhouseCrud.WithHttpProxy(cfg.Data.Clickhouse.GetHttpProxy()))
	}

	if cfg.Data.Clickhouse.ConnectionOpenStrategy != nil {
		options = append(options, clickhouseCrud.WithConnectionOpenStrategy(cfg.Data.Clickhouse.GetConnectionOpenStrategy()))
	}

	if cfg.Data.Clickhouse.Scheme != nil {
		options = append(options, clickhouseCrud.WithScheme(cfg.Data.Clickhouse.GetScheme()))
	}

	if cfg.Data.Clickhouse.BlockBufferSize != nil {
		options = append(options, clickhouseCrud.WithBlockBufferSize(int(cfg.Data.Clickhouse.GetBlockBufferSize())))
	}

	c, err := clickhouseCrud.NewClient(options...)

	return c, err
}
