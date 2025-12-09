package influxdb

import (
	"errors"

	"github.com/go-kratos/kratos/v2/log"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"

	influxdbCrud "github.com/tx7do/go-crud/influxdb"
)

func NewClient(logger log.Logger, cfg *conf.Bootstrap) (*influxdbCrud.Client, error) {
	if cfg.Data == nil || cfg.Data.Influxdb == nil {
		return nil, errors.New("influxdb config is nil")
	}

	var options []influxdbCrud.Option

	if logger != nil {
		options = append(options, influxdbCrud.WithLogger(logger))
	}

	if cfg.Data.Influxdb.GetHost() != "" {
		options = append(options, influxdbCrud.WithHost(cfg.Data.Influxdb.GetHost()))
	}
	if cfg.Data.Influxdb.GetToken() != "" {
		options = append(options, influxdbCrud.WithToken(cfg.Data.Influxdb.GetToken()))
	}
	if cfg.Data.Influxdb.GetDatabase() != "" {
		options = append(options, influxdbCrud.WithDatabase(cfg.Data.Influxdb.GetDatabase()))
	}
	if cfg.Data.Influxdb.GetOrganization() != "" {
		options = append(options, influxdbCrud.WithOrganization(cfg.Data.Influxdb.GetOrganization()))
	}
	if cfg.Data.Influxdb.WriteTimeout != nil {
		options = append(options, influxdbCrud.WithWriteTimeout(cfg.Data.Influxdb.GetWriteTimeout().AsDuration()))
	}
	if cfg.Data.Influxdb.QueryTimeout != nil {
		options = append(options, influxdbCrud.WithQueryTimeout(cfg.Data.Influxdb.GetQueryTimeout().AsDuration()))
	}
	if cfg.Data.Influxdb.IdleConnectionTimeout != nil {
		options = append(options, influxdbCrud.WithIdleConnectionTimeout(cfg.Data.Influxdb.GetIdleConnectionTimeout().AsDuration()))
	}
	if cfg.Data.Influxdb.MaxIdleConnections != nil {
		options = append(options, influxdbCrud.WithMaxIdleConnections(int(cfg.Data.Influxdb.GetMaxIdleConnections())))
	}
	if cfg.Data.Influxdb.AuthScheme != nil {
		options = append(options, influxdbCrud.WithAuthScheme(cfg.Data.Influxdb.GetAuthScheme()))
	}

	return influxdbCrud.NewClient(options...)
}
