package bootstrap

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/zipkin"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	traceSdk "go.opentelemetry.io/otel/sdk/trace"
	semConv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/tx7do/kratos-bootstrap/gen/api/go/conf/v1"
)

// NewTracerExporter 创建一个导出器，支持：jaeger和zipkin
func NewTracerExporter(exporterName, endpoint string, insecure bool) (traceSdk.SpanExporter, error) {
	if exporterName == "" {
		exporterName = "jaeger"
	}

	ctx := context.Background()

	switch exporterName {
	case "zipkin":
		return NewZipkinExporter(ctx, endpoint)
	case "jaeger":
		fallthrough
	case "otlp-http":
		return NewOtlpHttpExporter(ctx, endpoint, insecure)
	case "otlp-grpc":
		return NewOtlpGrpcExporter(ctx, endpoint, insecure)
	case "stdout":
		return stdouttrace.New()
	default:
		return nil, errors.New("exporter type not support")
	}
}

// NewTracerProvider 创建一个链路追踪器
func NewTracerProvider(cfg *conf.Tracer, serviceInfo *ServiceInfo) error {
	if cfg == nil {
		return errors.New("tracer config is nil")
	}

	if cfg.Sampler == 0 {
		cfg.Sampler = 1.0
	}

	if cfg.Env == "" {
		cfg.Env = "dev"
	}

	opts := []traceSdk.TracerProviderOption{
		traceSdk.WithSampler(traceSdk.ParentBased(traceSdk.TraceIDRatioBased(cfg.Sampler))),
		traceSdk.WithResource(resource.NewSchemaless(
			semConv.ServiceNameKey.String(serviceInfo.Name),
			semConv.ServiceVersionKey.String(serviceInfo.Version),
			semConv.ServiceInstanceIDKey.String(serviceInfo.Id),
			attribute.String("env", cfg.Env),
		)),
	}

	if len(cfg.Endpoint) > 0 {
		exp, err := NewTracerExporter(cfg.Batcher, cfg.Endpoint)
		if err != nil {
			panic(err)
		}

		opts = append(opts, traceSdk.WithBatcher(exp))
	}

	tp := traceSdk.NewTracerProvider(opts...)
	if tp == nil {
		return errors.New("create tracer provider failed")
	}

	otel.SetTracerProvider(tp)

	return nil
}

// NewZipkinExporter 创建一个zipkin导出器
func NewZipkinExporter(_ context.Context, endpoint string) (traceSdk.SpanExporter, error) {
	return zipkin.New(endpoint)
}

//// NewJaegerExporter 创建一个jaeger导出器
//func NewJaegerExporter(_ context.Context, endpoint string) (traceSdk.SpanExporter, error) {
//	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
//}

// NewOtlpHttpExporter 创建一个OTLP HTTP导出器
func NewOtlpHttpExporter(ctx context.Context, endpoint string, insecure bool, options ...otlptracehttp.Option) (traceSdk.SpanExporter, error) {
	var opts []otlptracehttp.Option
	opts = append(opts, otlptracehttp.WithEndpoint(endpoint))

	if insecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	opts = append(opts, options...)

	return otlptrace.New(
		ctx,
		otlptracehttp.NewClient(opts...),
	)
}

// NewOtlpGrpcExporter 创建一个OTLP GRPC导出器
func NewOtlpGrpcExporter(ctx context.Context, endpoint string, insecure bool, options ...otlptracegrpc.Option) (traceSdk.SpanExporter, error) {
	var opts []otlptracegrpc.Option
	opts = append(opts, otlptracegrpc.WithEndpoint(endpoint))

	if insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	opts = append(opts, options...)

	return otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(opts...),
	)
}
