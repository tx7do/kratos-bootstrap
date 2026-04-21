package elasticsearch

import (
	"errors"
	"os"

	"github.com/go-kratos/kratos/v2/log"
	elasticsearchCrud "github.com/tx7do/go-crud/elasticsearch"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func NewClient(logger log.Logger, cfg *conf.Bootstrap, opts ...elasticsearchCrud.Option) (*elasticsearchCrud.Client, error) {
	if cfg.Data == nil || cfg.Data.Elasticsearch == nil {
		return nil, errors.New("elasticsearch config is nil")
	}

	var options []elasticsearchCrud.Option

	if logger != nil {
		options = append(options, elasticsearchCrud.WithLogger(logger))
	}
	if len(cfg.Data.Elasticsearch.GetAddresses()) > 0 {
		options = append(options, elasticsearchCrud.WithAddresses(cfg.Data.Elasticsearch.GetAddresses()...))
	}
	if cfg.Data.Elasticsearch.GetUsername() != "" {
		options = append(options, elasticsearchCrud.WithUsername(cfg.Data.Elasticsearch.GetUsername()))
	}
	if cfg.Data.Elasticsearch.GetPassword() != "" {
		options = append(options, elasticsearchCrud.WithPassword(cfg.Data.Elasticsearch.GetPassword()))
	}

	if cfg.Data.Elasticsearch.EnableMetrics != nil {
		options = append(options, elasticsearchCrud.WithEnableMetrics(cfg.Data.Elasticsearch.GetEnableMetrics()))
	}
	if cfg.Data.Elasticsearch.EnableDebugLogger != nil {
		options = append(options, elasticsearchCrud.WithEnableDebugLogger(cfg.Data.Elasticsearch.GetEnableDebugLogger()))
	}
	if cfg.Data.Elasticsearch.EnableCompatibilityMode != nil {
		options = append(options, elasticsearchCrud.WithEnableCompatibilityMode(cfg.Data.Elasticsearch.GetEnableCompatibilityMode()))
	}
	if cfg.Data.Elasticsearch.DisableMetaHeader != nil {
		options = append(options, elasticsearchCrud.WithDisableMetaHeader(cfg.Data.Elasticsearch.GetDisableMetaHeader()))
	}
	if cfg.Data.Elasticsearch.DiscoverNodesOnStart != nil {
		options = append(options, elasticsearchCrud.WithDiscoverNodesOnStart(cfg.Data.Elasticsearch.GetDiscoverNodesOnStart()))
	}
	if cfg.Data.Elasticsearch.DiscoverNodesInterval != nil {
		options = append(options, elasticsearchCrud.WithDiscoverNodesInterval(cfg.Data.Elasticsearch.GetDiscoverNodesInterval().AsDuration()))
	}

	if cfg.Data.Elasticsearch.DisableRetry != nil {
		options = append(options, elasticsearchCrud.WithDisableRetry(cfg.Data.Elasticsearch.GetDisableRetry()))
	}
	if cfg.Data.Elasticsearch.MaxRetries != nil {
		options = append(options, elasticsearchCrud.WithMaxRetries(int(cfg.Data.Elasticsearch.GetMaxRetries())))
	}
	if cfg.Data.Elasticsearch.CompressRequestBody != nil {
		options = append(options, elasticsearchCrud.WithCompressRequestBody(cfg.Data.Elasticsearch.GetCompressRequestBody()))
	}
	if cfg.Data.Elasticsearch.CompressRequestBodyLevel != nil {
		options = append(options, elasticsearchCrud.WithCompressRequestBodyLevel(int(cfg.Data.Elasticsearch.GetCompressRequestBodyLevel())))
	}
	if cfg.Data.Elasticsearch.PoolCompressor != nil {
		options = append(options, elasticsearchCrud.WithPoolCompressor(cfg.Data.Elasticsearch.GetPoolCompressor()))
	}

	if cfg.Data.Elasticsearch.CloudId != nil {
		options = append(options, elasticsearchCrud.WithCloudID(cfg.Data.Elasticsearch.GetCloudId()))
	}
	if cfg.Data.Elasticsearch.GetApiKey() != "" {
		options = append(options, elasticsearchCrud.WithAPIKey(cfg.Data.Elasticsearch.GetApiKey()))
	}
	if cfg.Data.Elasticsearch.GetServiceToken() != "" {
		options = append(options, elasticsearchCrud.WithServiceToken(cfg.Data.Elasticsearch.GetServiceToken()))
	}
	if cfg.Data.Elasticsearch.GetCertificateFingerprint() != "" {
		options = append(options, elasticsearchCrud.WithCertificateFingerprint(cfg.Data.Elasticsearch.GetCertificateFingerprint()))
	}

	if cfg.Data.Opensearch.Tls != nil {
		if caData, err := loadCACertData(cfg.Data.Opensearch.Tls); err != nil {
			return nil, err
		} else {
			options = append(options, elasticsearchCrud.WithCACert(caData))
		}
	}

	if opts != nil {
		options = append(options, opts...)
	}

	return elasticsearchCrud.NewElasticsearchClient(options...)
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
