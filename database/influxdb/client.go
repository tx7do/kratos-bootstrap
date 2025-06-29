package influxdb

import (
	"context"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"github.com/go-kratos/kratos/v2/log"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

type Client struct {
	cli *influxdb3.Client

	log *log.Helper
}

func NewClient(logger log.Logger, cfg *conf.Bootstrap) (*Client, error) {
	c := &Client{
		log: log.NewHelper(log.With(logger, "module", "influxdb-client")),
	}

	if err := c.createInfluxdbClient(cfg); err != nil {
		return nil, err
	}

	return c, nil
}

// createInfluxdbClient 创建InfluxDB客户端
func (c *Client) createInfluxdbClient(cfg *conf.Bootstrap) error {
	if cfg.Data == nil || cfg.Data.Influxdb == nil {
		return nil
	}

	client, err := influxdb3.New(influxdb3.ClientConfig{
		Host:         cfg.Data.Influxdb.GetHost(),
		Token:        cfg.Data.Influxdb.GetToken(),
		Database:     cfg.Data.Influxdb.GetDatabase(),
		Organization: cfg.Data.Influxdb.GetOrganization(),
	})
	if err != nil {
		c.log.Errorf("failed to create influxdb client: %v", err)
		return err
	}

	c.cli = client

	return nil
}

// Close 关闭InfluxDB客户端
func (c *Client) Close() {
	if c.cli == nil {
		c.log.Warn("influxdb client is nil, nothing to close")
		return
	}

	if err := c.cli.Close(); err != nil {
		c.log.Errorf("failed to close influxdb client: %v", err)
	} else {
		c.log.Info("influxdb client closed successfully")
	}
}

// Query 查询数据
func (c *Client) Query(ctx context.Context, query string) (*influxdb3.QueryIterator, error) {
	if c.cli == nil {
		return nil, ErrInfluxDBClientNotInitialized
	}

	result, err := c.cli.Query(
		ctx,
		query,
		influxdb3.WithQueryType(influxdb3.InfluxQL),
	)
	if err != nil {
		c.log.Errorf("failed to query data: %v", err)
		return nil, ErrInfluxDBQueryFailed
	}

	return result, nil
}

func (c *Client) QueryWithParams(
	ctx context.Context,
	table string,
	filters map[string]interface{},
	operators map[string]string,
	fields []string,
) (*influxdb3.QueryIterator, error) {
	if c.cli == nil {
		return nil, ErrInfluxDBClientNotInitialized
	}

	query := BuildQueryWithParams(table, filters, operators, fields)
	result, err := c.cli.Query(
		ctx,
		query,
		influxdb3.WithQueryType(influxdb3.InfluxQL),
	)
	if err != nil {
		c.log.Errorf("failed to query data: %v", err)
		return nil, ErrInfluxDBQueryFailed
	}

	return result, nil
}

// Insert 插入数据
func (c *Client) Insert(ctx context.Context, point *influxdb3.Point) error {
	if c.cli == nil {
		return ErrInfluxDBClientNotInitialized
	}
	if point == nil {
		return ErrInvalidPoint
	}

	points := []*influxdb3.Point{point}
	if err := c.cli.WritePoints(ctx, points); err != nil {
		c.log.Errorf("failed to insert data: %v", err)
		return ErrInsertFailed
	}

	return nil
}

// BatchInsert 批量插入数据
func (c *Client) BatchInsert(ctx context.Context, points []*influxdb3.Point) error {
	if c.cli == nil {
		return ErrInfluxDBClientNotInitialized
	}

	if len(points) == 0 {
		return ErrNoPointsToInsert
	}

	if err := c.cli.WritePoints(ctx, points); err != nil {
		c.log.Errorf("failed to batch insert data: %v", err)
		return ErrBatchInsertFailed
	}

	return nil
}
