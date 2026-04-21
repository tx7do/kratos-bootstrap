package opensearch

import (
	"errors"
	"os"

	"github.com/go-kratos/kratos/v2/log"
	opensearchCrud "github.com/tx7do/go-crud/opensearch"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func NewClient(logger log.Logger, cfg *conf.Bootstrap, opts ...opensearchCrud.Option) (*opensearchCrud.Client, error) {
	if cfg.Data == nil || cfg.Data.Opensearch == nil {
		return nil, errors.New("opensearch config is nil")
	}

	var options []opensearchCrud.Option

	if logger != nil {
		options = append(options, opensearchCrud.WithLogger(logger))
	}
	if len(cfg.Data.Opensearch.GetAddresses()) > 0 {
		options = append(options, opensearchCrud.WithAddresses(cfg.Data.Opensearch.GetAddresses()...))
	}
	if cfg.Data.Opensearch.GetUsername() != "" {
		options = append(options, opensearchCrud.WithUsername(cfg.Data.Opensearch.GetUsername()))
	}
	if cfg.Data.Opensearch.GetPassword() != "" {
		options = append(options, opensearchCrud.WithPassword(cfg.Data.Opensearch.GetPassword()))
	}

	if cfg.Data.Opensearch.RetryOnStatus != nil {
		var retryOnStatus []int
		for _, status := range cfg.Data.Opensearch.GetRetryOnStatus() {
			retryOnStatus = append(retryOnStatus, int(status))
		}
		options = append(options, opensearchCrud.WithRetryOnStatus(retryOnStatus...))
	}
	if cfg.Data.Opensearch.DisableRetry != nil {
		options = append(options, opensearchCrud.WithDisableRetry(cfg.Data.Opensearch.GetDisableRetry()))
	}
	if cfg.Data.Opensearch.EnableRetryOnTimeout != nil {
		options = append(options, opensearchCrud.WithEnableRetryOnTimeout(cfg.Data.Opensearch.GetEnableRetryOnTimeout()))
	}
	if cfg.Data.Opensearch.MaxRetries != nil {
		options = append(options, opensearchCrud.WithMaxRetries(int(cfg.Data.Opensearch.GetMaxRetries())))
	}

	if cfg.Data.Opensearch.CompressRequestBody != nil {
		options = append(options, opensearchCrud.WithCompressRequestBody(cfg.Data.Opensearch.GetCompressRequestBody()))
	}

	if cfg.Data.Opensearch.DiscoverNodesOnStart != nil {
		options = append(options, opensearchCrud.WithDiscoverNodesOnStart(cfg.Data.Opensearch.GetDiscoverNodesOnStart()))
	}
	if cfg.Data.Opensearch.DiscoverNodesInterval != nil {
		options = append(options, opensearchCrud.WithDiscoverNodesInterval(cfg.Data.Opensearch.GetDiscoverNodesInterval().AsDuration()))
	}

	if cfg.Data.Opensearch.EnableMetrics != nil {
		options = append(options, opensearchCrud.WithEnableMetrics(cfg.Data.Opensearch.GetEnableMetrics()))
	}
	if cfg.Data.Opensearch.EnableDebugLogger != nil {
		options = append(options, opensearchCrud.WithEnableDebugLogger(cfg.Data.Opensearch.GetEnableDebugLogger()))
	}

	if cfg.Data.Opensearch.Tls != nil {
		if caData, err := loadCACertData(cfg.Data.Opensearch.Tls); err != nil {
			return nil, err
		} else {
			options = append(options, opensearchCrud.WithCACert(caData))
		}
	}

	if opts != nil {
		options = append(options, opts...)
	}

	return opensearchCrud.NewOpenSearchClient(options...)
}

func loadCACertData(cfg *conf.TLS) ([]byte, error) {
	if cfg == nil {
		return nil, nil
	}

	if cfg.File != nil {
		caPath := cfg.File.GetCaPath()
		if caPath == "" {
			return nil, errors.New("CA path is empty")
		}
		caData, err := os.ReadFile(caPath)
		if err != nil {
			return nil, err
		}
		return caData, nil
	} else if cfg.Config != nil {
		return cfg.Config.CaPem, nil
	}

	return nil, errors.New("invalid TLS config: no CA certificate found")
}
