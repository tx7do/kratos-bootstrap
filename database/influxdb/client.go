package influxdb

import (
	"github.com/InfluxCommunity/influxdb3-go/influxdb3"

	"github.com/go-kratos/kratos/v2/log"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func NewInfluxClient(cfg *conf.Bootstrap, l *log.Helper) *influxdb3.Client {
	if cfg.Data == nil || cfg.Data.Influxdb == nil {
		l.Warn("influxdb config is nil")
		return nil
	}

	client, err := influxdb3.New(influxdb3.ClientConfig{
		Host:         cfg.Data.Influxdb.Address,
		Token:        cfg.Data.Influxdb.Token,
		Database:     cfg.Data.Influxdb.Bucket,
		Organization: cfg.Data.Influxdb.Organization,
	})
	if err != nil {
		l.Fatalf("failed opening connection to influxdb: %v", err)
		return nil
	}

	return client
}
