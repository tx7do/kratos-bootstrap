package clickhouse

import (
	"context"
	"crypto/tls"
	"database/sql"
	"net/url"

	clickhouseV2 "github.com/ClickHouse/clickhouse-go/v2"
	driverV2 "github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	"github.com/go-kratos/kratos/v2/log"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/utils"
)

type Creator func() any

var compressionMap = map[string]clickhouseV2.CompressionMethod{
	"none":    clickhouseV2.CompressionNone,
	"zstd":    clickhouseV2.CompressionZSTD,
	"lz4":     clickhouseV2.CompressionLZ4,
	"lz4hc":   clickhouseV2.CompressionLZ4HC,
	"gzip":    clickhouseV2.CompressionGZIP,
	"deflate": clickhouseV2.CompressionDeflate,
	"br":      clickhouseV2.CompressionBrotli,
}

type Client struct {
	log *log.Helper

	conn clickhouseV2.Conn
	db   *sql.DB
}

func NewClient(logger log.Logger, cfg *conf.Bootstrap) (*Client, error) {
	c := &Client{
		log: log.NewHelper(log.With(logger, "module", "clickhouse-client")),
	}

	if err := c.createClickHouseClient(cfg); err != nil {
		return nil, err
	}

	return c, nil
}

// createClickHouseClient 创建ClickHouse客户端
func (c *Client) createClickHouseClient(cfg *conf.Bootstrap) error {
	if cfg.Data == nil || cfg.Data.Clickhouse == nil {
		return nil
	}

	opts := &clickhouseV2.Options{}

	if cfg.Data.Clickhouse.Dsn != nil {
		tmp, err := clickhouseV2.ParseDSN(cfg.Data.Clickhouse.GetDsn())
		if err != nil {
			c.log.Errorf("failed to parse clickhouse DSN: %v", err)
			return ErrInvalidDSN
		}
		opts = tmp
	}

	if cfg.Data.Clickhouse.Addresses != nil {
		opts.Addr = cfg.Data.Clickhouse.GetAddresses()
	}

	if cfg.Data.Clickhouse.Database != nil ||
		cfg.Data.Clickhouse.Username != nil ||
		cfg.Data.Clickhouse.Password != nil {
		opts.Auth = clickhouseV2.Auth{}

		if cfg.Data.Clickhouse.Database != nil {
			opts.Auth.Database = cfg.Data.Clickhouse.GetDatabase()
		}
		if cfg.Data.Clickhouse.Username != nil {
			opts.Auth.Username = cfg.Data.Clickhouse.GetUsername()
		}
		if cfg.Data.Clickhouse.Password != nil {
			opts.Auth.Password = cfg.Data.Clickhouse.GetPassword()
		}
	}

	if cfg.Data.Clickhouse.Debug != nil {
		opts.Debug = cfg.Data.Clickhouse.GetDebug()
	}

	if cfg.Data.Clickhouse.MaxOpenConns != nil {
		opts.MaxOpenConns = int(cfg.Data.Clickhouse.GetMaxOpenConns())
	}
	if cfg.Data.Clickhouse.MaxIdleConns != nil {
		opts.MaxIdleConns = int(cfg.Data.Clickhouse.GetMaxIdleConns())
	}

	if cfg.Data.Clickhouse.Tls != nil {
		var tlsCfg *tls.Config
		var err error

		if tlsCfg, err = utils.LoadServerTlsConfig(cfg.Server.Grpc.Tls); err != nil {
			panic(err)
		}

		if tlsCfg != nil {
			opts.TLS = tlsCfg
		}
	}

	if cfg.Data.Clickhouse.CompressionMethod != nil || cfg.Data.Clickhouse.CompressionLevel != nil {
		opts.Compression = &clickhouseV2.Compression{}

		if cfg.Data.Clickhouse.GetCompressionMethod() != "" {
			opts.Compression.Method = compressionMap[cfg.Data.Clickhouse.GetCompressionMethod()]
		}
		if cfg.Data.Clickhouse.CompressionLevel != nil {
			opts.Compression.Level = int(cfg.Data.Clickhouse.GetCompressionLevel())
		}
	}
	if cfg.Data.Clickhouse.MaxCompressionBuffer != nil {
		opts.MaxCompressionBuffer = int(cfg.Data.Clickhouse.GetMaxCompressionBuffer())
	}

	if cfg.Data.Clickhouse.DialTimeout != nil {
		opts.DialTimeout = cfg.Data.Clickhouse.GetDialTimeout().AsDuration()
	}
	if cfg.Data.Clickhouse.ReadTimeout != nil {
		opts.ReadTimeout = cfg.Data.Clickhouse.GetReadTimeout().AsDuration()
	}
	if cfg.Data.Clickhouse.ConnMaxLifetime != nil {
		opts.ConnMaxLifetime = cfg.Data.Clickhouse.GetConnMaxLifetime().AsDuration()
	}

	if cfg.Data.Clickhouse.HttpProxy != nil {
		proxyURL, err := url.Parse(cfg.Data.Clickhouse.GetHttpProxy())
		if err != nil {
			c.log.Errorf("failed to parse HTTP proxy URL: %v", err)
			return ErrInvalidProxyURL
		}

		opts.HTTPProxyURL = proxyURL
	}

	if cfg.Data.Clickhouse.ConnectionOpenStrategy != nil {
		strategy := clickhouseV2.ConnOpenInOrder
		switch cfg.Data.Clickhouse.GetConnectionOpenStrategy() {
		case "in_order":
			strategy = clickhouseV2.ConnOpenInOrder
		case "round_robin":
			strategy = clickhouseV2.ConnOpenRoundRobin
		case "random":
			strategy = clickhouseV2.ConnOpenRandom
		}
		opts.ConnOpenStrategy = strategy
	}

	if cfg.Data.Clickhouse.Scheme != nil {
		switch cfg.Data.Clickhouse.GetScheme() {
		case "http":
			opts.Protocol = clickhouseV2.HTTP
		case "https":
			opts.Protocol = clickhouseV2.HTTP
		default:
			opts.Protocol = clickhouseV2.Native
		}
	}

	if cfg.Data.Clickhouse.BlockBufferSize != nil {
		opts.BlockBufferSize = uint8(cfg.Data.Clickhouse.GetBlockBufferSize())
	}

	// 创建ClickHouse连接
	conn, err := clickhouseV2.Open(opts)
	if err != nil {
		c.log.Errorf("failed to create clickhouse client: %v", err)
		return ErrConnectionFailed
	}

	c.conn = conn

	return nil
}

// Close 关闭ClickHouse客户端连接
func (c *Client) Close() {
	if c.conn == nil {
		c.log.Warn("clickhouse client is already closed or not initialized")
		return
	}

	if err := c.conn.Close(); err != nil {
		c.log.Errorf("failed to close clickhouse client: %v", err)
	} else {
		c.log.Info("clickhouse client closed successfully")
	}
}

// GetServerVersion 获取ClickHouse服务器版本
func (c *Client) GetServerVersion() string {
	if c.conn == nil {
		c.log.Error("clickhouse client is not initialized")
		return ""
	}

	version, err := c.conn.ServerVersion()
	if err != nil {
		c.log.Errorf("failed to get server version: %v", err)
		return ""
	} else {
		c.log.Infof("ClickHouse server version: %s", version)
		return version.String()
	}
}

// CheckConnection 检查ClickHouse客户端连接是否正常
func (c *Client) CheckConnection(ctx context.Context) error {
	if c.conn == nil {
		c.log.Error("clickhouse client is not initialized")
		return ErrClientNotInitialized
	}

	if err := c.conn.Ping(ctx); err != nil {
		c.log.Errorf("ping failed: %v", err)
		return ErrPingFailed
	}

	c.log.Info("clickhouse client connection is healthy")
	return nil
}

// Query 执行查询并返回结果
func (c *Client) Query(ctx context.Context, creator Creator, results *[]any, query string, args ...interface{}) error {
	if c.conn == nil {
		c.log.Error("clickhouse client is not initialized")
		return ErrClientNotInitialized
	}
	if creator == nil {
		c.log.Error("creator function cannot be nil")
		return ErrCreatorFunctionNil
	}

	rows, err := c.conn.Query(ctx, query, args...)
	if err != nil {
		c.log.Errorf("query failed: %v", err)
		return ErrQueryExecutionFailed
	}
	defer func(rows driverV2.Rows) {
		if err = rows.Close(); err != nil {
			c.log.Errorf("failed to close rows: %v", err)
		}
	}(rows)

	for rows.Next() {
		row := creator()
		if err = rows.ScanStruct(row); err != nil {
			c.log.Errorf("failed to scan row: %v", err)
			return ErrRowScanFailed
		}
		*results = append(*results, row)
	}

	// 检查是否有未处理的错误
	if rows.Err() != nil {
		c.log.Errorf("Rows iteration error: %v", rows.Err())
		return ErrRowsIterationError
	}

	return nil
}

// QueryRow 执行查询并返回单行结果
func (c *Client) QueryRow(ctx context.Context, dest any, query string, args ...interface{}) error {
	row := c.conn.QueryRow(ctx, query, args...)
	if row == nil {
		c.log.Error("query row returned nil")
		return ErrRowNotFound
	}

	if err := row.ScanStruct(dest); err != nil {
		c.log.Errorf("")
		return ErrRowScanFailed
	}

	return nil
}

// Select 封装 SELECT 子句
func (c *Client) Select(ctx context.Context, dest any, query string, args ...interface{}) error {
	if c.conn == nil {
		c.log.Error("clickhouse client is not initialized")
		return ErrClientNotInitialized
	}

	err := c.conn.Select(ctx, dest, query, args...)
	if err != nil {
		c.log.Errorf("select failed: %v", err)
		return ErrQueryExecutionFailed
	}

	return nil
}

// Exec 执行非查询语句
func (c *Client) Exec(ctx context.Context, query string, args ...interface{}) error {
	if c.conn == nil {
		c.log.Error("clickhouse client is not initialized")
		return ErrClientNotInitialized
	}

	if err := c.conn.Exec(ctx, query, args...); err != nil {
		c.log.Errorf("exec failed: %v", err)
		return ErrExecutionFailed
	}

	return nil
}

// AsyncInsert 异步插入数据
func (c *Client) AsyncInsert(ctx context.Context, query string, wait bool, args ...interface{}) error {
	if c.conn == nil {
		c.log.Error("clickhouse client is not initialized")
		return ErrClientNotInitialized
	}

	if err := c.conn.AsyncInsert(ctx, query, wait, args...); err != nil {
		c.log.Errorf("exec failed: %v", err)
		return ErrAsyncInsertFailed
	}

	return nil
}

// BatchInsert 批量插入数据
func (c *Client) BatchInsert(ctx context.Context, query string, data [][]interface{}) error {
	batch, err := c.conn.PrepareBatch(ctx, query)
	if err != nil {
		c.log.Errorf("failed to prepare batch: %v", err)
		return ErrBatchPrepareFailed
	}

	for _, row := range data {
		if err := batch.Append(row...); err != nil {
			c.log.Errorf("failed to append data: %v", err)
			return ErrBatchAppendFailed
		}
	}

	if err = batch.Send(); err != nil {
		c.log.Errorf("failed to send batch: %v", err)
		return ErrBatchSendFailed
	}

	return nil
}
