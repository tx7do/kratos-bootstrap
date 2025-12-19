package registry

import (
	"fmt"
	"sort"
	"sync"

	"github.com/go-kratos/kratos/v2/registry"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// DiscoveryFactory .
type DiscoveryFactory func(cfg *conf.Registry) (registry.Discovery, error)

var (
	muDF               sync.RWMutex
	discoveryFactories = make(map[Type]DiscoveryFactory)
)

// NewDiscovery 使用已注册的工厂创建 discovery，找不到时返回错误
func NewDiscovery(cfg *conf.Registry) (registry.Discovery, error) {
	if cfg == nil {
		return nil, nil
	}

	if cfg.GetType() == "" {
		return nil, nil
	}

	name := Type(cfg.GetType())

	f, ok := GetDiscoveryFactory(name)
	if !ok {
		return nil, fmt.Errorf("registry: discovery factory %q not found", name)
	}
	return f(cfg)
}

// RegisterDiscoveryFactory 注册一个名为 name 的 DiscoveryFactory，已存在时返回错误
func RegisterDiscoveryFactory(name Type, f DiscoveryFactory) error {
	if name == "" {
		return fmt.Errorf("registry: discovery name is empty")
	}
	if f == nil {
		return fmt.Errorf("registry: discovery factory is nil")
	}
	muDF.Lock()
	defer muDF.Unlock()
	if _, exists := discoveryFactories[name]; exists {
		return fmt.Errorf("registry: discovery %q already registered", name)
	}
	discoveryFactories[name] = f
	return nil
}

// MustRegisterDiscoveryFactory 在注册失败时 panic（便于 init 中使用）
func MustRegisterDiscoveryFactory(name Type, f DiscoveryFactory) {
	if err := RegisterDiscoveryFactory(name, f); err != nil {
		panic(err)
	}
}

// GetDiscoveryFactory 返回名为 name 的 DiscoveryFactory 及是否存在
func GetDiscoveryFactory(name Type) (DiscoveryFactory, bool) {
	muDF.RLock()
	defer muDF.RUnlock()
	f, ok := discoveryFactories[name]
	return f, ok
}

// ListDiscoveryFactories 返回已注册的 discovery 名称有序列表
func ListDiscoveryFactories() []string {
	muDF.RLock()
	defer muDF.RUnlock()
	names := make([]string, 0, len(discoveryFactories))
	for k := range discoveryFactories {
		names = append(names, string(k))
	}
	sort.Strings(names)
	return names
}
