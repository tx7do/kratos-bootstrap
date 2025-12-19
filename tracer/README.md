# 链路追踪 (tracer)

本包为项目提供统一的链路追踪（Tracing）封装，基于 OpenTelemetry SDK 实现。目标是：

- 为大多数使用者提供开箱即用的、通过配置启用的 tracing 功能；
- 对外暴露小而稳定的 API（大多数用户不需要直接引用 OpenTelemetry 包）；
- 支持 IoC / 插件式扩展：第三方可在运行时注册自定义 exporter；
- 支持优雅的生命周期管理（创建与关闭/Flush）。

注意：tracer 的具体上报与高级能力依赖 OpenTelemetry（otel）。本包对外 API 力求保持简单、稳定；只有实现自定义 exporter 的包才需要显式 import otel。

---

## 设计要点

- Exporter 注册表（Registry）
  - 提供 `RegisterExporter(name string, factory ExporterFactory)`，第三方可在自己的 `init()` 中注册自定义 exporter。
  - 内置的 exporter（`zipkin`、`otlp-http`、`otlp-grpc`、`stdout`）在包初始化时注册。
  - Factory 签名为：`func(ctx context.Context, cfg *conf.Tracer) (traceSdk.SpanExporter, error)`，factory 可读取配置字段（endpoint、insecure、headers、tls 等）。

- 配置与语义
  - `exporter`：指定后端类型（例如 `otlp-grpc`, `otlp-http`, `zipkin`, `stdout`）。
  - `endpoint`：后端地址（例如 "`localhost:4317`" 或 "`http://host:9411/api/v2/spans`"）。
  - `batcher_options`：批处理（BatchSpanProcessor）相关的策略（队列大小、批量大小、调度延迟等）。`batcher` 字段历史上被用于指定 exporter 名称，建议改用 `exporter`，短期保留 `batcher` 的兼容回退处理。

- Context 传播与错误处理
  - 所有 exporter 的创建和 shutdown 均支持 `context.Context`，不会在库内部使用 `context.Background()`。
  - 不在库内部 `panic`，遇到可恢复错误会返回 `error` 给调用者。

- 生命周期管理
  - 提供 `NewTracerProviderWithShutdown(ctx, cfg, appInfo)` 返回 `tp` 与 `shutdown` 函数，建议在 `main` 的退出路径调用 `shutdown(ctx)` 以优雅停机（flush traces）。
  - 也提供兼容的 `NewTracerProvider`（仅返回 error），以及 `ShutdownTracerProvider(ctx)`（关闭全局 provider）供过渡使用。

---

## 配置示例（YAML）

推荐的 tracer 配置（示例）：

```yaml
tracer:
  exporter: "otlp-grpc"        # 后端类型：otlp-grpc / otlp-http / zipkin / stdout
  endpoint: "localhost:4317"   # 后端地址
  sampler: 1.0                 # 采样率，0.0-1.0，默认 1.0
  env: "dev"                   # 运行环境
  insecure: true               # 是否使用 insecure（用于 OTLP）
  batcher_options:
    enabled: true
    max_queue_size: 2048
    max_export_batch_size: 512
    schedule_delay_millis: 5000
    export_timeout_millis: 30000
```

---

## 快速开始（最小示例）

下面示例展示如何在 `main` 中创建 tracer provider 并优雅关闭：

```go
package main

import (
    "context"
    "log"
    "time"

    conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
    "github.com/tx7do/kratos-bootstrap/tracer"
)

func main() {
    ctx := context.Background()

    cfg := &conf.Tracer{
        Exporter: "otlp-grpc",
        Endpoint: "localhost:4317",
        Sampler:  1.0,
        Env:      "prod",
        Insecure: true,
    }

    appInfo := &conf.AppInfo{
        AppId:      "my-service",
        InstanceId: "instance-1",
        Name:       "my-service",
        Version:    "v1.0.0",
    }

    tp, shutdown, err := tracer.NewTracerProviderWithShutdown(ctx, cfg, appInfo)
    if err != nil {
        log.Fatalf("failed to create tracer provider: %v", err)
    }
    // 在这里可以使用 tp 创建 tracer（或直接使用全局 otel.GetTracerProvider）

    // 在程序退出时调用 shutdown
    defer func() {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        if err := shutdown(ctx); err != nil {
            log.Printf("tracer shutdown error: %v", err)
        }
    }()

    // ... 应用主逻辑
}
```

如果你想使用兼容的旧函数：

```go
// 与上面等效（兼容旧签名）；但不返回 shutdown，建议逐步切换到 WithShutdown
_ = tracer.NewTracerProvider(ctx, cfg, appInfo)
defer tracer.ShutdownTracerProvider(ctx)
```

---

## 注册自定义 exporter（插件 / IoC）

如果需要在使用者侧实现并注册自定义 exporter（例如内部企业后端、特殊编码、鉴权等），可在该包的 `init()` 中注册 factory：

```go
package myexporter

import (
    "context"
    traceSdk "go.opentelemetry.io/otel/sdk/trace"
    conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
    "github.com/tx7do/kratos-bootstrap/tracer"
)

func init() {
    tracer.RegisterExporter("my-exporter", func(ctx context.Context, cfg *conf.Tracer) (traceSdk.SpanExporter, error) {
        // 根据 cfg 构建你的 exporter（可能使用 otel SDK 或自定义实现）
        return NewMyExporter(ctx, cfg) // 返回 traceSdk.SpanExporter
    })
}
```

说明：
- 唯一需要注意的是，编写自定义 exporter 的包通常需要 import OpenTelemetry 的 SDK 包（这是合理的：只有实现者才需直接依赖 otel）。
- 普通使用者仅需在配置中指定 `exporter: "my-exporter"` 即可使用。

---

## 生命周期与优雅关闭

- 强烈建议在应用退出路径调用 `shutdown(ctx)`（或者调用 `ShutdownTracerProvider(ctx)`）以确保 batcher/queue 中的数据被 flush 到后端。  
- `NewTracerProviderWithShutdown` 会返回一个 `shutdown` 函数，该函数会调用底层 SDK 的 `Shutdown`。示例见上文 `Quick Start`。

---

## 迁移与兼容说明（`batcher` -> `exporter`）

- 历史原因，旧配置中 `batcher` 字段曾被用来指定 exporter 名称。为了语义清晰，推荐改用 `exporter` 字段。  
- 本包实现会优先读取 `exporter` 字段；若为空则回退读取 `batcher` 字段以保持向后兼容，并在日志中输出弃用告警。  
- 建议在下一个 major 版本中移除 `batcher` 字段。

迁移建议步骤：
1. 在配置中把 `batcher: "zipkin"` → `exporter: "zipkin"`。  
2. 配置 `batcher_options`（或保留默认）以控制批量策略。  
3. 在应用代码中把 `NewTracerProvider` 替换为 `NewTracerProviderWithShutdown` 并在退出时调用 shutdown。

---

## 测试与调试建议

- 单元测试：为 registry、NewTracerExporter、NewTracerProviderWithShutdown 编写测试，覆盖正常路径与错误路径（例如未知 exporter）。
- 并发测试：如果添加自定义注册逻辑，使用 `go test -race` 检查 Register/Get 并发安全性。
- 本地调试：可以使用 `stdout` exporter 进行本地打印调试（`exporter: "stdout"`）。
- 集成测试：可在 CI 中运行一个临时 OTLP/Zipkin 接收进程，验证 traces 能成功发送并被接收。

---

## 常见问题（FAQ）

Q: 我必须在项目中 import OpenTelemetry 吗？
A: 对于普通使用者（只通过配置使用 tracer），不需要直接 import otel；但是实现自定义 exporter 的包会需要 import otel SDK。tracer 包本身会在 go.mod 中间接带入 otel 依赖。

Q: 如果 exporter 创建失败会发生什么？
A: factory 创建 exporter 的错误会被返回给调用方（不会在库内部 panic）。调用方应决定是否重试或退出。

---

## API 概览（常用函数）

- `RegisterExporter(name string, f ExporterFactory)` — 注册 exporter factory。  
- `NewTracerExporter(ctx context.Context, exporterName string, cfg *conf.Tracer) (traceSdk.SpanExporter, error)` — 通过 registry 构建 exporter。  
- `NewTracerProviderWithShutdown(ctx context.Context, cfg *conf.Tracer, appInfo *conf.AppInfo) (traceSdk.TracerProvider, func(context.Context) error, error)` — 创建 provider 并返回 shutdown 函数。  
- `NewTracerProvider(ctx context.Context, cfg *conf.Tracer, appInfo *conf.AppInfo) error` — 兼容签名（仅返回 error）。  
- `ShutdownTracerProvider(ctx context.Context) error` — 关闭全局 provider（可重复调用）。  
- `ListExporterNames() []string` — 列出已注册的 exporter 名称。
