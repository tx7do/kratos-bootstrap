package tracer

import (
	"context"
	"errors"
	"sync"

	traceSdk "go.opentelemetry.io/otel/sdk/trace"

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
		return NewOtlpHttpExporter(ctx, cfg.GetEndpoint(), cfg.GetInsecure())
	})
	RegisterExporter(string(OtlpGrpc), func(ctx context.Context, cfg *conf.Tracer) (traceSdk.SpanExporter, error) {
		return NewOtlpGrpcExporter(ctx, cfg.GetEndpoint(), cfg.GetInsecure())
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
