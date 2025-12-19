package elasticsearch

import (
	"crypto/tls"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
	elasticsearchCrud "github.com/tx7do/go-crud/elasticsearch"
	tlsUtils "github.com/tx7do/go-utils/tls"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func NewClient(logger log.Logger, cfg *conf.Bootstrap) (*elasticsearchCrud.Client, error) {
	if cfg.Data == nil || cfg.Data.ElasticSearch == nil {
		return nil, errors.New("elasticsearch config is nil")
	}

	var options []elasticsearchCrud.Option

	if logger != nil {
		options = append(options, elasticsearchCrud.WithLogger(logger))
	}
	if cfg.Data.ElasticSearch.GetAddresses() != nil {
		options = append(options, elasticsearchCrud.WithAddresses(cfg.Data.ElasticSearch.GetAddresses()...))
	}
	if cfg.Data.ElasticSearch.GetUsername() != "" {
		options = append(options, elasticsearchCrud.WithUsername(cfg.Data.ElasticSearch.GetUsername()))
	}
	if cfg.Data.ElasticSearch.GetPassword() != "" {
		options = append(options, elasticsearchCrud.WithPassword(cfg.Data.ElasticSearch.GetPassword()))
	}
	if cfg.Data.ElasticSearch.Tls != nil {
		var tlsCfg *tls.Config
		var err error

		if tlsCfg, err = loadServerTlsConfig(cfg.Server.Grpc.Tls); err != nil {
			return nil, err
		}

		if tlsCfg != nil {
			options = append(options, elasticsearchCrud.WithTLSConfig(tlsCfg))
		}
	}
	if cfg.Data.ElasticSearch.GetEnableMetrics() {
		options = append(options, elasticsearchCrud.WithEnableMetrics(cfg.Data.ElasticSearch.GetEnableMetrics()))
	}
	if cfg.Data.ElasticSearch.GetEnableDebugLogger() {
		options = append(options, elasticsearchCrud.WithEnableDebugLogger(cfg.Data.ElasticSearch.GetEnableDebugLogger()))
	}
	if cfg.Data.ElasticSearch.GetEnableCompatibilityMode() {
		options = append(options, elasticsearchCrud.WithEnableCompatibilityMode(cfg.Data.ElasticSearch.GetEnableCompatibilityMode()))
	}
	if cfg.Data.ElasticSearch.GetDisableMetaHeader() {
		options = append(options, elasticsearchCrud.WithDisableMetaHeader(cfg.Data.ElasticSearch.GetDisableMetaHeader()))
	}
	if cfg.Data.ElasticSearch.GetDiscoverNodesOnStart() {
		options = append(options, elasticsearchCrud.WithDiscoverNodesOnStart(cfg.Data.ElasticSearch.GetDiscoverNodesOnStart()))
	}
	if cfg.Data.ElasticSearch.DiscoverNodesInterval != nil {
		options = append(options, elasticsearchCrud.WithDiscoverNodesInterval(cfg.Data.ElasticSearch.GetDiscoverNodesInterval().AsDuration()))
	}
	if cfg.Data.ElasticSearch.GetDisableRetry() {
		options = append(options, elasticsearchCrud.WithDisableRetry(cfg.Data.ElasticSearch.GetDisableRetry()))
	}
	if cfg.Data.ElasticSearch.MaxRetries != nil {
		options = append(options, elasticsearchCrud.WithMaxRetries(int(cfg.Data.ElasticSearch.GetMaxRetries())))
	}
	if cfg.Data.ElasticSearch.GetCompressRequestBody() {
		options = append(options, elasticsearchCrud.WithCompressRequestBody(cfg.Data.ElasticSearch.GetCompressRequestBody()))
	}
	if cfg.Data.ElasticSearch.GetCompressRequestBodyLevel() != 0 {
		options = append(options, elasticsearchCrud.WithCompressRequestBodyLevel(int(cfg.Data.ElasticSearch.GetCompressRequestBodyLevel())))
	}
	if cfg.Data.ElasticSearch.GetPoolCompressor() {
		options = append(options, elasticsearchCrud.WithPoolCompressor(cfg.Data.ElasticSearch.GetPoolCompressor()))
	}
	if cfg.Data.ElasticSearch.CloudId != nil {
		options = append(options, elasticsearchCrud.WithCloudID(cfg.Data.ElasticSearch.GetCloudId()))
	}
	if cfg.Data.ElasticSearch.GetApiKey() != "" {
		options = append(options, elasticsearchCrud.WithAPIKey(cfg.Data.ElasticSearch.GetApiKey()))
	}
	if cfg.Data.ElasticSearch.GetServiceToken() != "" {
		options = append(options, elasticsearchCrud.WithServiceToken(cfg.Data.ElasticSearch.GetServiceToken()))
	}
	if cfg.Data.ElasticSearch.GetCertificateFingerprint() != "" {
		options = append(options, elasticsearchCrud.WithCertificateFingerprint(cfg.Data.ElasticSearch.GetCertificateFingerprint()))
	}

	return elasticsearchCrud.NewClient(options...)
}

func loadServerTlsConfig(cfg *conf.TLS) (*tls.Config, error) {
	if cfg == nil {
		return nil, nil
	}

	var tlsCfg *tls.Config
	var err error

	if cfg.File != nil {
		if tlsCfg, err = tlsUtils.LoadServerTlsConfigFile(
			cfg.File.GetKeyPath(),
			cfg.File.GetCertPath(),
			cfg.File.GetCaPath(),
			cfg.InsecureSkipVerify,
		); err != nil {
			return nil, err
		}
	} else if cfg.Config != nil {
		if tlsCfg, err = tlsUtils.LoadServerTlsConfigString(
			cfg.Config.GetKeyPem(),
			cfg.Config.GetCertPem(),
			cfg.Config.GetCaPem(),
			cfg.InsecureSkipVerify,
		); err != nil {
			return nil, err
		}
	}

	return tlsCfg, err
}
