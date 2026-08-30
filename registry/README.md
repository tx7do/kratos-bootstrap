# registry 包说明

## 概述

`registry` 包提供线程安全的工厂注册与创建机制，支持两类组件：

- `Registrar`（服务注册器）
- `Discovery`（服务发现）

通过按名称注册工厂函数，运行时按字符串选择并创建对应实例，便于扩展与解耦实现与配置。

## 特性

- 按名称注册工厂并防止重复注册
- 提供 `MustRegister*` 便于在 `init` 中自动注册（注册失败会 panic）
- 线程安全（使用 `sync.RWMutex` 保护）
- 支持列举已注册工厂名称（有序）

## API 概览

下列签名基于包内部定义，`conf` 为项目配置类型，返回类型基于 `github.com/go-kratos/kratos/v2/registry`。

- `type RegistrarFactory func(cfg *conf.Registry) (registry.Registrar, error)`
- `RegisterRegistrarFactory(name string, f RegistrarFactory) error`
- `MustRegisterRegistrarFactory(name string, f RegistrarFactory)`
- `GetRegistrarFactory(name string) (RegistrarFactory, bool)`
- `NewRegistrar(name string, cfg *conf.Registry) (registry.Registrar, error)`
- `ListRegistrarFactories() []string`

- `type DiscoveryFactory func(cfg *conf.Registry) (registry.Discovery, error)`
- `RegisterDiscoveryFactory(name string, f DiscoveryFactory) error`
- `MustRegisterDiscoveryFactory(name string, f DiscoveryFactory)`
- `GetDiscoveryFactory(name string) (DiscoveryFactory, bool)`
- `NewDiscovery(name string, cfg *conf.Registry) (registry.Discovery, error)`
- `ListDiscoveryFactories() []string`

注意：`NewRegistrar` / `NewDiscovery` 需要传入工厂名称（`name`）和配置对象（`cfg`）。

## 使用示例

```go
package example

import (
	kRegistry "github.com/go-kratos/kratos/v2/registry"
	bRegistry "github.com/tx7do/kratos-bootstrap/registry"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	//_ "github.com/tx7do/kratos-bootstrap/registry/consul"
	//_ "github.com/tx7do/kratos-bootstrap/registry/etcd"
	//_ "github.com/tx7do/kratos-bootstrap/registry/eureka"
	//_ "github.com/tx7do/kratos-bootstrap/registry/kubernetes"
	//_ "github.com/tx7do/kratos-bootstrap/registry/nacos"
	//_ "github.com/tx7do/kratos-bootstrap/registry/servicecomb"
	//_ "github.com/tx7do/kratos-bootstrap/registry/zookeeper"
)

// NewRegistry 创建一个注册客户端
func NewRegistry(cfg *conf.Registry) (kRegistry.Registrar, error) {
	return bRegistry.NewRegistrar(cfg)
}

// NewDiscovery 创建一个发现客户端
func NewDiscovery(cfg *conf.Registry) (kRegistry.Discovery, error) {
	return bRegistry.NewDiscovery(cfg)
}
```

## 配置参数参考

各注册中心支持通过 `conf.Registry` 配置注入的参数如下（以 YAML 配置为例）：

### Consul

```yaml
registry:
  type: consul
  consul:
    scheme: http                # 网络样式: http、https
    address: 127.0.0.1:8500
    health_check: true
    datacenter: dc1             # 数据中心
    token: xxx                  # ACL Token
    token_file: /path/token     # ACL Token 文件路径
    namespace: ns               # 命名空间(企业版)
    partition: default          # 管理分区(企业版)
    path_prefix: /consul        # API网关路径前缀
    basic_auth:                 # HTTP Basic 认证
      username: user
      password: pass
    wait_time: 10s              # Watch 阻塞等待的最长时间
    tls:                        # TLS 配置(详见 conf.TLS)
      insecure_skip_verify: false
      file:
        cert_path: /path/cert.pem
        key_path: /path/key.pem
        ca_path: /path/ca.pem
    heartbeat: true             # 心跳检查开关,默认开启
    health_check_interval: 10   # 健康检查间隔(秒),默认10
    deregister_critical_service_after: 600  # 不健康多久后注销服务(秒),默认600
    timeout: 10s                # 服务发现超时时间,默认10s
```

### Etcd

```yaml
registry:
  type: etcd
  etcd:
    endpoints: [127.0.0.1:2379]
    username: root
    password: pass
    tls: {}                     # TLS 配置
    dial_timeout: 5s
    auto_sync_interval: 1m      # 集群成员自动同步间隔
    dial_keep_alive_time: 30s
    dial_keep_alive_timeout: 10s
    reject_old_cluster: false
    permit_without_stream: false
    namespace: /microservices   # 键前缀
    register_ttl: 15s           # 注册租约 TTL
    max_retry: 5                # 心跳重试次数
```

### Nacos

```yaml
registry:
  type: nacos
  nacos:
    address: 127.0.0.1
    port: 8848
    namespace_id: ""
    region_id: ""
    app_name: ""
    app_key: ""
    access_key: ""              # 阿里云 AccessKey
    secret_key: ""              # 阿里云 SecretKey
    username: nacos
    password: nacos
    timeout: 10s
    beat_interval: 5s
    update_thread_num: 20
    not_load_cache_at_start: false
    update_cache_when_empty: false
    open_kms: false
    log_level: info
    log_dir: /tmp/nacos/log
    cache_dir: /tmp/nacos/cache
    context_path: /nacos
    scheme: http                # 网络样式: http、https
    grpc_port: 9848             # gRPC 长连接端口,默认 port+1000
    endpoint: ""                # 地址服务器 endpoint
    disable_use_snap_shot: false
    append_to_stdout: false
    async_update_service: false # 异步订阅更新服务列表
    group: DEFAULT_GROUP        # 服务分组名
    cluster: DEFAULT            # 集群名
    weight: 100                 # 实例权重
    kind: grpc                  # 默认服务协议类型
```

### ZooKeeper

```yaml
registry:
  type: zookeeper
  zookeeper:
    endpoints: [127.0.0.1:2181]
    timeout: 5s
    namespace: /microservices   # 根路径
    username: user              # Digest ACL 用户名
    password: pass              # Digest ACL 密码
```

### Eureka

```yaml
registry:
  type: eureka
  eureka:
    endpoints: [http://127.0.0.1:18761]
    heartbeat_interval: 10s
    refresh_interval: 30s
    path: eureka/v2
    max_retry: 3                # 请求失败重试次数,默认 endpoints 数量
```

### Polaris

```yaml
registry:
  type: polaris
  polaris:
    address: 127.0.0.1          # 服务端地址(设置后 SDK 直连该地址)
    port: 8091
    instance_count: 1
    namespace: default
    service: ""
    token: ""
    config_file: ""             # SDK 配置文件路径(优先级最高)
    weight: 100
    priority: 0
    healthy: true
    isolate: false
    heartbeat: true
    ttl: 5                      # 心跳上报 TTL(秒),开启心跳时默认5
    protocol: grpc
    timeout: 1s
    retry_count: 0
```

### Servicecomb

```yaml
registry:
  type: servicecomb
  servicecomb:
    endpoints: [127.0.0.1:30100]
    enable_ssl: false
    timeout: 10s                # 默认10s
    enable_auth: false
    username: ""
    password: ""
    token_expiration: 0s
    verbose: false
    compressed: false
    tls: {}                     # TLS 配置(启用后自动开启 SSL)
```

### Kubernetes

```yaml
registry:
  type: kubernetes
  kubernetes:
    namespace: default          # 默认读取 ServiceAccount 的 namespace
    kubeconfig: ""              # 默认依次尝试 InCluster 配置和 ~/.kube/config
    in_cluster: false           # 强制使用集群内配置
```
