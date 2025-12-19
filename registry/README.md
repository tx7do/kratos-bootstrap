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
