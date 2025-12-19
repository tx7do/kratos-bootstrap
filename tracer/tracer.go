package tracer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	traceSdk "go.opentelemetry.io/otel/sdk/trace"
	semConv "go.opentelemetry.io/otel/semconv/v1.4.0"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

var (
	// tpInstance holds the currently active global tracer provider (interface type, nil-able)
	tpMu       sync.Mutex
	tpInstance *traceSdk.TracerProvider
)

// NewTracerExporter 构建 exporter：优先使用注册表中的 factory。
// exporterName 不能为空，cfg 传入给 factory 用于读取 endpoint/insecure/headers 等信息。
func NewTracerExporter(ctx context.Context, cfg *conf.Tracer) (traceSdk.SpanExporter, error) {
	if cfg == nil {
		return nil, errors.New("tracer cfg is nil")
	}
	if cfg.GetExporter() == "" {
		return nil, errors.New("exporter name is empty")
	}

	if f, ok := GetExporterFactory(cfg.GetExporter()); ok {
		return f(ctx, cfg)
	}
	return nil, fmt.Errorf("unknown exporter %q; available: %v", cfg.GetExporter(), ListExporterNames())
}

// ShutdownTracerProvider gracefully shuts down the active global tracer provider (if set).
// Safe to call multiple times; returns error if shutdown fails.
func ShutdownTracerProvider(ctx context.Context) error {
	tpMu.Lock()
	defer tpMu.Unlock()
	if tpInstance == nil {
		return nil
	}
	if err := tpInstance.Shutdown(ctx); err != nil {
		return err
	}
	tpInstance = nil
	return nil
}

// NewTracerProvider 创建 tracer provider 并设置为全局 provider。
// 注：为了更好的资源管理，可以使用 NewTracerProviderWithShutdown 获得 shutdown 函数。
func NewTracerProvider(ctx context.Context, cfg *conf.Tracer, appInfo *conf.AppInfo) error {
	_, _, err := NewTracerProviderWithShutdown(ctx, cfg, appInfo)
	return err
}

// NewTracerProviderWithShutdown 返回 (tp, shutdownFunc, error)，推荐在 main 中使用并在退出时调用 shutdownFunc(ctx)
func NewTracerProviderWithShutdown(ctx context.Context, cfg *conf.Tracer, appInfo *conf.AppInfo) (*traceSdk.TracerProvider, func(context.Context) error, error) {
	if cfg == nil || appInfo == nil {
		return nil, func(context.Context) error { return nil }, nil
	}

	// do not mutate caller cfg; use local defaults
	sampler := cfg.GetSampler()
	if sampler == 0 {
		sampler = 1.0
	}
	env := cfg.GetEnv()
	if env == "" {
		env = "dev"
	}

	opts := []traceSdk.TracerProviderOption{
		traceSdk.WithSampler(traceSdk.ParentBased(traceSdk.TraceIDRatioBased(sampler))),
		traceSdk.WithResource(resource.NewSchemaless(
			semConv.ServiceNameKey.String(appInfo.GetName()),
			semConv.ServiceVersionKey.String(appInfo.GetVersion()),
			semConv.ServiceInstanceIDKey.String(appInfo.GetAppId()),
			attribute.String("env", env),
		)),
	}

	// NOTE: cfg.GetBatcher() historically used as exporter name in this project.
	// Consider renaming conf.Tracer fields to `Exporter`/`ExporterName` in the future.
	if len(cfg.GetEndpoint()) > 0 {
		exp, err := NewTracerExporter(ctx, cfg)
		if err != nil {
			return nil, nil, err
		}
		opts = append(opts, traceSdk.WithBatcher(exp))
	}

	tp := traceSdk.NewTracerProvider(opts...)

	// defensive check (NewTracerProvider does not return nil in normal cases)
	if tp == nil {
		return nil, nil, errors.New("create tracer provider failed")
	}

	// set global provider and keep reference for shutdown
	tpMu.Lock()
	// shutdown previous provider if present
	if tpInstance != nil {
		_ = tpInstance.Shutdown(ctx)
	}
	tpInstance = tp
	tpMu.Unlock()

	otel.SetTracerProvider(tp)

	shutdown := func(c context.Context) error {
		// shutdown global provider if it's still the same one
		return ShutdownTracerProvider(c)
	}

	return tp, shutdown, nil
}
