package doris

import (
	"errors"

	"github.com/go-kratos/kratos/v2/log"
	dorisCrud "github.com/tx7do/go-crud/doris"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func NewClient(logger log.Logger, cfg *conf.Bootstrap) (*dorisCrud.Client, error) {
	if cfg.Data == nil || cfg.Data.Doris == nil {
		return nil, errors.New("doris config is nil")
	}

	var options []dorisCrud.Option

	if logger != nil {
		options = append(options, dorisCrud.WithLogger(log.NewHelper(logger)))
	}
	if len(cfg.Data.Doris.GetDsn()) != 0 {
		options = append(options, dorisCrud.WithDSN(cfg.Data.Doris.GetDsn()))
	}

	if cfg.Data.Doris.MaxOpenConnections != nil {
		options = append(options, dorisCrud.WithMaxOpenConns(int(cfg.Data.Doris.GetMaxOpenConnections())))
	}
	if cfg.Data.Doris.MaxIdleConnections != nil {
		options = append(options, dorisCrud.WithMaxIdleConns(int(cfg.Data.Doris.GetMaxIdleConnections())))
	}

	if cfg.Data.Doris.ConnectionMaxLifetime != nil {
		options = append(options, dorisCrud.WithConnMaxLifetime(cfg.Data.Doris.GetConnectionMaxLifetime().AsDuration()))
	}

	if cfg.Data.Doris.StreamLoad != nil {
		if len(cfg.Data.Doris.StreamLoad.Endpoint) != 0 {
			options = append(options, dorisCrud.WithStreamLoadEndpoint(cfg.Data.Doris.StreamLoad.GetEndpoint()))
		}
		if len(cfg.Data.Doris.StreamLoad.GetUsername()) != 0 && len(cfg.Data.Doris.StreamLoad.GetPassword()) != 0 {
			options = append(options, dorisCrud.WithStreamLoadAuth(cfg.Data.Doris.StreamLoad.GetUsername(), cfg.Data.Doris.StreamLoad.GetPassword()))
		}
		if cfg.Data.Doris.StreamLoad.Timeout != nil {
			options = append(options, dorisCrud.WithStreamLoadTimeout(cfg.Data.Doris.StreamLoad.GetTimeout().AsDuration()))
		}
		if cfg.Data.Doris.StreamLoad.Method != nil {
			options = append(options, dorisCrud.WithStreamLoadMethod(cfg.Data.Doris.StreamLoad.GetMethod()))
		}
	}

	c, err := dorisCrud.NewClient(options...)

	return c, err
}
