<p align="center">
  <h1 align="center">kratos-bootstrap</h1>
  <p align="center">
    基于 go-kratos 的微服务应用引导框架
  </p>
  <p align="center">
    <em>一站式基础设施引导，让微服务开发从此告别重复搭建</em>
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

## 项目亮点

- **统一引导入口**：封装配置加载、日志初始化、服务注册、链路追踪等全流程，一行代码启动微服务
- **九大数据库引擎**：原生支持 ClickHouse、Doris、Elasticsearch、OpenSearch、MongoDB、InfluxDB、Cassandra、Ent、GORM 等主流数据库
- **七大注册中心**：集成 Consul、Etcd、Nacos、Zookeeper、Eureka、Polaris、ServiceComb 等主流服务发现组件
- **六大远程配置中心**：支持 Apollo、Consul、Etcd、Nacos、Polaris、Kubernetes ConfigMap 远程配置管理
- **五大传输层**：内置 Kafka、Asynq、MCP、MQTT、SSE 等消息与流式传输组件
- **三大 AI 框架**：集成 go-openai、LangChainGo、Eino，支持云端模型与本地模型（Ollama）无缝切换
- **多语言脚本引擎**：支持 JavaScript（ES6+）与 Lua 脚本动态扩展业务逻辑
- **全链路可观测**：基于 OpenTelemetry 的分布式链路追踪，支持 Jaeger 等多种 Exporter
- **六种日志框架**：支持 Zap、Zerolog、Logrus、Fluent、阿里云 SLS、腾讯云 CLS 日志输出
- **配置契约优先**：全部配置通过 Protobuf 定义，类型安全、强约束、自动生成

---

## 技术栈

| 层级 | 技术 | 说明 |
| --- | --- | --- |
| 语言 | Go 1.23+ | 高性能编译型语言 |
| 框架 | go-kratos v2 | B站开源微服务框架 |
| 配置定义 | Protobuf + buf.build | 接口契约优先，类型安全 |
| 链路追踪 | OpenTelemetry | 分布式可观测标准 |
| 对象存储 | MinIO / S3 | S3 兼容对象存储 |
| 缓存 | Redis | go-redis v9，支持链路追踪与指标采集 |
| AI 框架 | go-openai / LangChainGo / Eino | 多框架 LLM 接入，支持云端与本地模型 |
| 消息队列 | Kafka | 高吞吐事件流处理 |
| 异步任务 | Asynq | 基于 Redis 的异步任务队列 |
| AI 工具协议 | MCP (Model Context Protocol) | AI Agent 工具调用标准协议 |
| 实时推送 | SSE | 服务端事件推送 |
| 物联网 | MQTT | 轻量级消息传输协议 |
| 脚本引擎 | goja（JS）/ gopher-lua | 多语言动态脚本执行 |

---

## 项目结构

```
kratos-bootstrap/
├── api/                                # Protobuf API 定义与生成代码
│   ├── protos/conf/v1/                 # .proto 源文件（配置结构定义）
│   └── gen/go/conf/v1/                 # buf 生成的 Go 代码
├── bootstrap/                          # 应用引导核心
│   ├── bootstrap.go                    # 引导入口（配置加载/日志/注册/追踪）
│   ├── cli.go                          # CLI 命令行框架（Cobra）
│   ├── context.go                      # 引导上下文
│   └── daemon.go                       # 守护进程支持
├── ai/                                 # AI 大模型集成
│   ├── model/                          # go-openai 原生客户端
│   ├── langchaingo/                    # LangChainGo 框架封装
│   └── eino/                           # 字节跳动 Eino 框架封装
├── config/                             # 远程配置中心
│   ├── apollo/                         # Apollo 配置中心
│   ├── consul/                         # Consul KV 配置
│   ├── etcd/                           # Etcd 配置
│   ├── nacos/                          # Nacos 配置
│   ├── polaris/                        # Polaris 配置
│   └── kubernetes/                     # Kubernetes ConfigMap
├── registry/                           # 服务注册与发现
│   ├── consul/                         # Consul 注册中心
│   ├── etcd/                           # Etcd 注册中心
│   ├── nacos/                          # Nacos 注册中心
│   ├── zookeeper/                      # Zookeeper 注册中心
│   ├── eureka/                         # Eureka 注册中心
│   ├── polaris/                        # Polaris 注册中心
│   ├── servicecomb/                    # ServiceComb 注册中心
│   └── kubernetes/                     # Kubernetes 注册中心
├── database/                           # 数据库客户端
│   ├── clickhouse/                     # ClickHouse（OLAP）
│   ├── doris/                          # Apache Doris（OLAP）
│   ├── elasticsearch/                  # Elasticsearch（搜索引擎）
│   ├── opensearch/                     # OpenSearch（搜索引擎）
│   ├── mongodb/                        # MongoDB（文档数据库）
│   ├── influxdb/                       # InfluxDB（时序数据库）
│   ├── cassandra/                      # Cassandra（宽列数据库）
│   ├── ent/                            # Ent ORM
│   └── gorm/                           # GORM ORM
├── transport/                          # 传输层组件
│   ├── kafka/                          # Kafka 消息队列
│   ├── asynq/                          # Asynq 异步任务
│   ├── mcp/                            # MCP AI 工具协议
│   ├── mqtt/                           # MQTT 物联网协议
│   └── sse/                            # SSE 服务端推送
├── logger/                             # 日志框架集成
│   ├── zap/                            # Zap 高性能日志
│   ├── zerolog/                        # Zerolog 零分配日志
│   ├── logrus/                         # Logrus 结构化日志
│   ├── fluent/                         # Fluentd 日志收集
│   ├── aliyun/                         # 阿里云 SLS 日志服务
│   └── tencent/                        # 腾讯云 CLS 日志服务
├── cache/                              # 缓存
│   └── redis/                          # Redis 客户端
├── tracer/                             # 链路追踪
│   └── exporter.go                     # OpenTelemetry Exporter 工厂
├── rpc/                                # RPC 通信
│   ├── grpc.go                         # gRPC 客户端/服务端
│   ├── rest.go                         # REST (HTTP) 客户端/服务端
│   └── middleware/                     # RPC 中间件
│       ├── validate/                   # Protobuf 参数校验
│       └── requestid/                  # 请求 ID 中间件
├── oss/                                # 对象存储
│   ├── minio/                          # MinIO 客户端
│   └── s3/                             # S3 兼容客户端
└── script_engine/                      # 脚本引擎
    └── script_engine.go                # JavaScript / Lua 脚本执行引擎
```

---

## 核心功能

### 应用引导

| 功能 | 说明 |
| --- | --- |
| 统一启动入口 | 通过 `Bootstrap` 函数封装配置加载、日志初始化、服务注册、链路追踪全流程 |
| CLI 框架 | 基于 Cobra 的命令行框架，支持子命令定制与 Flag 注入 |
| 守护进程 | 原生守护进程模式，支持后台运行与 PID 管理 |
| 应用元信息 | 统一管理应用名称、版本号、实例 ID、项目空间等元数据 |
| 优雅退出 | 内置信号捕获与优雅关停机制，确保服务安全下线 |

### 服务注册与发现

| 注册中心 | 说明 |
| --- | --- |
| Consul | HashiCorp 服务发现与 KV 存储 |
| Etcd | 高可用分布式 KV，Kubernetes 底层存储 |
| Nacos | 阿里云微服务注册与配置中心 |
| Zookeeper | Apache 分布式协调服务 |
| Eureka | Netflix 服务发现组件 |
| Polaris | 腾讯云北极星服务治理平台 |
| ServiceComb | 华为云微服务引擎 |
| Kubernetes | 原生 Kubernetes 服务发现 |

### 远程配置中心

| 配置中心 | 说明 |
| --- | --- |
| Apollo | 携程分布式配置管理中心 |
| Consul | Consul KV 配置存储 |
| Etcd | Etcd 分布式配置 |
| Nacos | Nacos 配置管理 |
| Polaris | Polaris 配置管理 |
| Kubernetes | Kubernetes ConfigMap |

### 数据库支持

| 数据库 | 类型 | 说明 |
| --- | --- | --- |
| ClickHouse | OLAP | 列式存储，极致分析性能 |
| Apache Doris | OLAP | 高性能实时分析引擎 |
| Elasticsearch | 搜索引擎 | 全文检索与日志分析 |
| OpenSearch | 搜索引擎 | Elasticsearch 开源分支 |
| MongoDB | 文档数据库 | 灵活的文档模型 |
| InfluxDB | 时序数据库 | 时序数据存储与查询 |
| Cassandra | 宽列数据库 | 高可用分布式存储 |
| Ent | ORM | Go 实体框架 |
| GORM | ORM | Go 最流行的 ORM |

### AI 大模型

| 模块 | 框架 | 说明 |
| --- | --- | --- |
| model | go-openai | OpenAI 兼容 API 原生客户端 |
| langchaingo | LangChainGo | Chain / Agent / Embedding / VectorStore / Memory |
| eino | 字节跳动 Eino | Chain / Tool / Prompt / Compose |

所有 AI 模块均支持：
- **云端模型**：OpenAI、通义千问、DeepSeek 等兼容 OpenAI API 的模型服务
- **本地模型**：Ollama 本地部署，零依赖云端

### 传输层

| 传输组件 | 说明 |
| --- | --- |
| Kafka | 高吞吐消息队列，支持 SASL/SCRAM 认证与 TLS 加密 |
| Asynq | 基于 Redis 的异步任务队列，支持单机/集群/Sentinel 模式 |
| MCP | Model Context Protocol，AI Agent 工具调用标准协议（HTTP/SSE/Stdio） |
| MQTT | 轻量级物联网消息传输协议 |
| SSE | Server-Sent Events 实时服务端推送 |

### 日志框架

| 日志框架 | 说明 |
| --- | --- |
| Zap | Uber 开源高性能日志库 |
| Zerolog | 零分配高性能 JSON 日志 |
| Logrus | 结构化日志库 |
| Fluent | Fluentd 日志收集 |
| 阿里云 SLS | 阿里云日志服务 |
| 腾讯云 CLS | 腾讯云日志服务 |

### 链路追踪

基于 OpenTelemetry 标准实现，支持：
- 采样率配置（TraceIDRatioBased 采样）
- W3C TraceContext 与 Baggage 传播
- 可插拔 Exporter 工厂模式
- 服务元信息自动注入（服务名、版本、实例 ID、环境）

### RPC 通信

| 功能 | 说明 |
| --- | --- |
| gRPC 服务端/客户端 | 原生 gRPC 支持，内置 TLS、超时、中间件链 |
| REST (HTTP) 服务端 | HTTP/RESTful API 服务，支持 CORS、pprof 调试 |
| 中间件链 | Recovery / Tracing / Validate / RateLimit / Metadata |
| 限流 | BBR（Google BBR 算法）自适应限流 |
| 白名单 | 基于方法名的中间件白名单机制 |

### 脚本引擎

| 语言 | 引擎 | 说明 |
| --- | --- | --- |
| JavaScript | goja | 支持 ES6+ 语法，适用于规则引擎与动态计算 |
| Lua | gopher-lua | 轻量级嵌入式脚本语言 |

支持脚本池管理（固定大小 / 自动伸缩），预加载脚本与入口脚本。

### 对象存储

| 存储 | 说明 |
| --- | --- |
| MinIO | 自托管 S3 兼容对象存储 |
| S3 | Amazon S3 / 兼容存储 |

---

## 设计原则

- **配置驱动**：全部组件通过 Protobuf 配置定义初始化，类型安全，避免硬编码
- **非侵入封装**：封装层仅做配置转译与实例创建，不侵入框架原生 API，保留框架原生使用方式
- **工厂模式**：注册中心、日志、追踪等组件均采用工厂注册模式，按需加载，松耦合
- **模块化设计**：每个功能模块独立 Go Module，按需引入，不产生冗余依赖
- **防御性编程**：全面的空指针检查与错误处理，确保组件缺失时不影响主流程

---

## 适用场景

- 微服务架构下的标准化应用引导与基础设施管理
- 需要集成多种数据库与中间件的复杂业务系统
- AI 应用开发，需对接多种大模型与 Agent 框架
- 物联网场景，需要 MQTT 消息传输与边缘计算
- 实时数据处理，需要 Kafka / Asynq 异步任务调度
- 需要动态脚本扩展业务规则的场景（如规则引擎、动态计算）

---

## 许可证

项目基于 MIT 许可证开源，允许自由使用、修改和分发，需保留原版权信息。
