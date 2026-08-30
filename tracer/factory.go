package tracer

import (
	"context"
	"crypto/tls"
	"errors"
	"sync"

	traceSdk "go.opentelemetry.io/otel/sdk/trace"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"google.golang.org/grpc/credentials"

	tlsUtils "github.com/tx7do/go-utils/tls"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// ExporterFactory is a creator function that returns a SpanExporter given a context and tracer config.
// Implementations may read additional fields from cfg (headers, tls, auth, etc.).
type ExporterFactory func(ctx context.Context, cfg *conf.Tracer) (traceSdk.SpanExporter, error)

var (
	// registry holds named exporter factories (concurrent-safe)
	exporterMu       sync.RWMutex
	exporterRegistry = map[string]ExporterFactory{}
)

func init() {
	// register built-in exporters
	RegisterExporter(string(Zipkin), func(ctx context.Context, cfg *conf.Tracer) (traceSdk.SpanExporter, error) {
		return NewZipkinExporter(ctx, cfg.GetEndpoint())
	})
	RegisterExporter(string(OtlpHttp), func(ctx context.Context, cfg *conf.Tracer) (traceSdk.SpanExporter, error) {
		return NewOtlpHttpExporter(ctx, cfg.GetEndpoint(), false, buildOtlpHttpOptions(cfg)...)
	})
	RegisterExporter(string(OtlpGrpc), func(ctx context.Context, cfg *conf.Tracer) (traceSdk.SpanExporter, error) {
		return NewOtlpGrpcExporter(ctx, cfg.GetEndpoint(), false, buildOtlpGrpcOptions(cfg)...)
	})
	RegisterExporter(string(Std), func(ctx context.Context, cfg *conf.Tracer) (traceSdk.SpanExporter, error) {
		return NewStdoutExporter(ctx)
	})

	// legacy/unsupported entries can be mapped to explicit errors
	RegisterExporter(string(Jaeger), func(ctx context.Context, cfg *conf.Tracer) (traceSdk.SpanExporter, error) {
		return nil, errors.New("tracer: jaeger exporter is not supported in this build; use otlp-http or otlp-grpc instead")
	})
}

// RegisterExporter registers an exporter factory under the given name.
func RegisterExporter(name string, f ExporterFactory) {
	exporterMu.Lock()
	defer exporterMu.Unlock()
	exporterRegistry[name] = f
}

// GetExporterFactory returns a registered factory by name.
func GetExporterFactory(name string) (ExporterFactory, bool) {
	exporterMu.RLock()
	defer exporterMu.RUnlock()
	f, ok := exporterRegistry[name]
	return f, ok
}

// ListExporterNames returns the currently registered exporter names.
func ListExporterNames() []string {
	exporterMu.RLock()
	defer exporterMu.RUnlock()
	names := make([]string, 0, len(exporterRegistry))
	for k := range exporterRegistry {
		names = append(names, k)
	}
	return names
}

// buildOtlpHttpOptions 从 tracer 配置构造OTLP/HTTP导出器选项(请求头/TLS/超时).
func buildOtlpHttpOptions(cfg *conf.Tracer) []otlptracehttp.Option {
	var opts []otlptracehttp.Option

	if h := cfg.GetHeaders(); len(h) > 0 {
		headers := make(map[string]string, len(h))
		for k, v := range h {
			headers[k] = v
		}
		opts = append(opts, otlptracehttp.WithHeaders(headers))
	}

	if tlsCfg := loadExporterTlsConfig(cfg); tlsCfg != nil {
		opts = append(opts, otlptracehttp.WithTLSClientConfig(tlsCfg))
	} else if cfg.GetInsecure() {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	if cfg.GetTimeout().AsDuration() > 0 {
		opts = append(opts, otlptracehttp.WithTimeout(cfg.GetTimeout().AsDuration()))
	}

	return opts
}

// buildOtlpGrpcOptions 从 tracer 配置构造OTLP/gRPC导出器选项(请求头/TLS/超时).
func buildOtlpGrpcOptions(cfg *conf.Tracer) []otlptracegrpc.Option {
	var opts []otlptracegrpc.Option

	if h := cfg.GetHeaders(); len(h) > 0 {
		headers := make(map[string]string, len(h))
		for k, v := range h {
			headers[k] = v
		}
		opts = append(opts, otlptracegrpc.WithHeaders(headers))
	}

	if tlsCfg := loadExporterTlsConfig(cfg); tlsCfg != nil {
		opts = append(opts, otlptracegrpc.WithTLSCredentials(credentials.NewTLS(tlsCfg)))
	} else if cfg.GetInsecure() {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	if cfg.GetTimeout().AsDuration() > 0 {
		opts = append(opts, otlptracegrpc.WithTimeout(cfg.GetTimeout().AsDuration()))
	}

	return opts
}

func loadExporterTlsConfig(cfg *conf.Tracer) *tls.Config {
	tlsConf := cfg.GetTls()
	if tlsConf == nil {
		return nil
	}

	var tlsCfg *tls.Config
	var err error

	if tlsConf.File != nil {
		if tlsCfg, err = tlsUtils.LoadClientTlsConfigFile(
			tlsConf.File.GetKeyPath(),
			tlsConf.File.GetCertPath(),
			tlsConf.File.GetCaPath(),
		); err != nil {
			return nil
		}
	} else if tlsConf.Config != nil {
		if tlsCfg, err = tlsUtils.LoadClientTlsConfigString(
			tlsConf.Config.GetKeyPem(),
			tlsConf.Config.GetCertPem(),
			tlsConf.Config.GetCaPem(),
		); err != nil {
			return nil
		}
	}

	if tlsCfg != nil && tlsConf.GetInsecureSkipVerify() {
		tlsCfg.InsecureSkipVerify = true
	}

	return tlsCfg
}
