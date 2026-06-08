<p align="center">
  <h1 align="center">kratos-bootstrap</h1>
  <p align="center">
    A Microservice Application Bootstrap Framework Built on go-kratos
  </p>
  <p align="center">
    <em>One-stop infrastructure bootstrapping — say goodbye to repetitive setup in microservice development</em>
  </p>
</p>

<p align="center">
  <a href="README.md">中文</a> · <a href="README_en.md">English</a> · <a href="README_ja.md">日本語</a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=Go" alt="Go Version" />
  <img src="https://img.shields.io/badge/Kratos-v2-00ADD8?style=flat-square" alt="Kratos" />
  <img src="https://img.shields.io/badge/License-MIT-green?style=flat-square" alt="License" />
  <img src="https://img.shields.io/badge/PRs-Welcome-brightgreen?style=flat-square" alt="PRs Welcome" />
</p>

---

## Highlights

- **Unified Bootstrap Entry**: Encapsulates config loading, logger initialization, service registration, and distributed tracing — launch your microservice in a single call
- **9 Database Engines**: Native support for ClickHouse, Doris, Elasticsearch, OpenSearch, MongoDB, InfluxDB, Cassandra, Ent, and GORM
- **7 Service Registries**: Integrated with Consul, Etcd, Nacos, Zookeeper, Eureka, Polaris, and ServiceComb
- **6 Remote Config Centers**: Apollo, Consul KV, Etcd, Nacos, Polaris, and Kubernetes ConfigMap
- **5 Transport Components**: Built-in Kafka, Asynq, MCP, MQTT, and SSE messaging & streaming
- **3 AI Frameworks**: go-openai, LangChainGo, and ByteDance Eino with seamless cloud/local model switching (Ollama)
- **Multi-language Script Engine**: JavaScript (ES6+) and Lua scripting for dynamic business logic
- **Full-chain Observability**: OpenTelemetry-based distributed tracing with pluggable exporters (Jaeger, etc.)
- **6 Logging Frameworks**: Zap, Zerolog, Logrus, Fluentd, Alibaba Cloud SLS, and Tencent Cloud CLS
- **Config-first with Protobuf**: All configurations defined via Protobuf — type-safe, strongly constrained, auto-generated

---

## Tech Stack

| Layer | Technology | Description |
| --- | --- | --- |
| Language | Go 1.23+ | High-performance compiled language |
| Framework | go-kratos v2 | Bilibili's open-source microservice framework |
| Config Definition | Protobuf + buf.build | Contract-first, type-safe configuration |
| Tracing | OpenTelemetry | Distributed observability standard |
| Object Storage | MinIO / S3 | S3-compatible object storage |
| Cache | Redis | go-redis v9 with tracing & metrics |
| AI Frameworks | go-openai / LangChainGo / Eino | Multi-framework LLM integration (cloud & local) |
| Message Queue | Kafka | High-throughput event streaming |
| Async Tasks | Asynq | Redis-based asynchronous task queue |
| AI Tool Protocol | MCP (Model Context Protocol) | Standard protocol for AI Agent tool invocation |
| Real-time Push | SSE | Server-Sent Events |
| IoT | MQTT | Lightweight messaging protocol |
| Script Engine | goja (JS) / gopher-lua | Multi-language dynamic scripting |

---

## Project Structure

```
kratos-bootstrap/
├── api/                                # Protobuf API definitions & generated code
│   ├── protos/conf/v1/                 # .proto source files (config schemas)
│   └── gen/go/conf/v1/                 # buf-generated Go code
├── bootstrap/                          # Application bootstrap core
│   ├── bootstrap.go                    # Bootstrap entry (config/logger/registry/tracer)
│   ├── cli.go                          # CLI framework (Cobra)
│   ├── context.go                      # Bootstrap context
│   └── daemon.go                       # Daemon process support
├── ai/                                 # LLM integration
│   ├── model/                          # go-openai native client
│   ├── langchaingo/                    # LangChainGo framework wrapper
│   └── eino/                           # ByteDance Eino framework wrapper
├── config/                             # Remote config centers
│   ├── apollo/                         # Apollo
│   ├── consul/                         # Consul KV
│   ├── etcd/                           # Etcd
│   ├── nacos/                          # Nacos
│   ├── polaris/                        # Polaris
│   └── kubernetes/                     # Kubernetes ConfigMap
├── registry/                           # Service registration & discovery
│   ├── consul/                         # Consul
│   ├── etcd/                           # Etcd
│   ├── nacos/                          # Nacos
│   ├── zookeeper/                      # Zookeeper
│   ├── eureka/                         # Eureka
│   ├── polaris/                        # Polaris
│   ├── servicecomb/                    # ServiceComb
│   └── kubernetes/                     # Kubernetes
├── database/                           # Database clients
│   ├── clickhouse/                     # ClickHouse (OLAP)
│   ├── doris/                          # Apache Doris (OLAP)
│   ├── elasticsearch/                  # Elasticsearch
│   ├── opensearch/                     # OpenSearch
│   ├── mongodb/                        # MongoDB
│   ├── influxdb/                       # InfluxDB (time-series)
│   ├── cassandra/                      # Cassandra
│   ├── ent/                            # Ent ORM
│   └── gorm/                           # GORM ORM
├── transport/                          # Transport layer
│   ├── kafka/                          # Kafka
│   ├── asynq/                          # Asynq async tasks
│   ├── mcp/                            # MCP AI tool protocol
│   ├── mqtt/                           # MQTT IoT protocol
│   └── sse/                            # SSE server push
├── logger/                             # Logging frameworks
│   ├── zap/                            # Zap
│   ├── zerolog/                        # Zerolog
│   ├── logrus/                         # Logrus
│   ├── fluent/                         # Fluentd
│   ├── aliyun/                         # Alibaba Cloud SLS
│   └── tencent/                        # Tencent Cloud CLS
├── cache/                              # Cache
│   └── redis/                          # Redis client
├── tracer/                             # Distributed tracing
│   └── exporter.go                     # OpenTelemetry exporter factory
├── rpc/                                # RPC communication
│   ├── grpc.go                         # gRPC client/server
│   ├── rest.go                         # REST (HTTP) client/server
│   └── middleware/                     # RPC middleware
│       ├── validate/                   # Protobuf validation
│       └── requestid/                  # Request ID middleware
├── oss/                                # Object storage
│   ├── minio/                          # MinIO client
│   └── s3/                             # S3-compatible client
└── script_engine/                      # Script engine
    └── script_engine.go                # JavaScript / Lua engine
```

---

## Core Features

### Application Bootstrap

| Feature | Description |
| --- | --- |
| Unified Entry | Single `Bootstrap` function handles config loading, logger init, service registration, and tracing |
| CLI Framework | Cobra-based CLI with subcommand customization and flag injection |
| Daemon Mode | Native daemon support with background execution and PID management |
| App Metadata | Unified management of app name, version, instance ID, and namespace |
| Graceful Shutdown | Built-in signal capture and graceful shutdown for safe service termination |

### Service Registration & Discovery

| Registry | Description |
| --- | --- |
| Consul | HashiCorp service discovery & KV store |
| Etcd | Highly available distributed KV (Kubernetes backbone) |
| Nacos | Alibaba Cloud microservice registration & config center |
| Zookeeper | Apache distributed coordination service |
| Eureka | Netflix service discovery |
| Polaris | Tencent Cloud Polaris service governance |
| ServiceComb | Huawei Cloud microservice engine |
| Kubernetes | Native Kubernetes service discovery |

### Remote Config Centers

| Config Center | Description |
| --- | --- |
| Apollo | Ctrip's distributed configuration management |
| Consul | Consul KV config storage |
| Etcd | Etcd distributed configuration |
| Nacos | Nacos configuration management |
| Polaris | Polaris configuration management |
| Kubernetes | Kubernetes ConfigMap |

### Database Support

| Database | Type | Description |
| --- | --- | --- |
| ClickHouse | OLAP | Columnar storage, extreme analytical performance |
| Apache Doris | OLAP | High-performance real-time analytics engine |
| Elasticsearch | Search Engine | Full-text search & log analytics |
| OpenSearch | Search Engine | Open-source Elasticsearch fork |
| MongoDB | Document DB | Flexible document model |
| InfluxDB | Time-series DB | Time-series data storage & querying |
| Cassandra | Wide-column DB | Highly available distributed storage |
| Ent | ORM | Go entity framework |
| GORM | ORM | Most popular Go ORM |

### AI / LLM Integration

| Module | Framework | Description |
| --- | --- | --- |
| model | go-openai | OpenAI-compatible API native client |
| langchaingo | LangChainGo | Chain / Agent / Embedding / VectorStore / Memory |
| eino | ByteDance Eino | Chain / Tool / Prompt / Compose |

All AI modules support:
- **Cloud Models**: OpenAI, Qwen, DeepSeek, and any OpenAI API-compatible model service
- **Local Models**: Ollama local deployment, zero cloud dependency

### Transport Layer

| Component | Description |
| --- | --- |
| Kafka | High-throughput message queue with SASL/SCRAM auth & TLS encryption |
| Asynq | Redis-based async task queue (standalone/cluster/sentinel) |
| MCP | Model Context Protocol — AI Agent tool invocation (HTTP/SSE/Stdio) |
| MQTT | Lightweight IoT messaging protocol |
| SSE | Server-Sent Events real-time push |

### Logging Frameworks

| Framework | Description |
| --- | --- |
| Zap | Uber's high-performance structured logger |
| Zerolog | Zero-allocation JSON logger |
| Logrus | Structured logger |
| Fluent | Fluentd log collection |
| Alibaba Cloud SLS | Alibaba Cloud Log Service |
| Tencent Cloud CLS | Tencent Cloud Log Service |

### Distributed Tracing

Built on the OpenTelemetry standard:
- Configurable sampling rates (TraceIDRatioBased sampling)
- W3C TraceContext & Baggage propagation
- Pluggable exporter factory pattern
- Automatic service metadata injection (service name, version, instance ID, environment)

### RPC Communication

| Feature | Description |
| --- | --- |
| gRPC Server/Client | Native gRPC with built-in TLS, timeout, and middleware chain |
| REST (HTTP) Server | HTTP/RESTful API service with CORS & pprof debugging |
| Middleware Chain | Recovery / Tracing / Validate / RateLimit / Metadata |
| Rate Limiting | BBR (Google BBR algorithm) adaptive rate limiting |
| Whitelist | Method-name-based middleware whitelist mechanism |

### Script Engine

| Language | Engine | Description |
| --- | --- | --- |
| JavaScript | goja | ES6+ support, ideal for rule engines & dynamic computation |
| Lua | gopher-lua | Lightweight embedded scripting |

Supports script pool management (fixed-size / auto-scaling), preload scripts, and entry scripts.

### Object Storage

| Storage | Description |
| --- | --- |
| MinIO | Self-hosted S3-compatible object storage |
| S3 | Amazon S3 / compatible storage |

---

## Design Principles

- **Config-Driven**: All components initialized via Protobuf-defined configurations — type-safe, no hardcoding
- **Non-intrusive Wrapping**: The wrapper layer only handles config translation and instance creation, preserving the framework's native API
- **Factory Pattern**: Registries, loggers, tracers, and other components use factory registration for on-demand loading and loose coupling
- **Modular Design**: Each feature is an independent Go module — import only what you need, no redundant dependencies
- **Defensive Programming**: Comprehensive nil checks and error handling ensure missing components don't break the main flow

---

## Use Cases

- Standardized microservice bootstrapping and infrastructure management
- Complex business systems requiring multiple databases and middleware
- AI application development with multi-LLM and Agent framework integration
- IoT scenarios requiring MQTT messaging and edge computing
- Real-time data processing with Kafka / Asynq async task scheduling
- Dynamic script-based business rule extension (rule engines, dynamic computation)

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
