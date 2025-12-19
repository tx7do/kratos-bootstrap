package registry

import (
	"fmt"
	"sort"
	"sync"

	"github.com/go-kratos/kratos/v2/registry"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// RegistrarFactory .
type RegistrarFactory func(cfg *conf.Registry) (registry.Registrar, error)

var (
	muRF               sync.RWMutex
	registrarFactories = make(map[Type]RegistrarFactory)
)

// NewRegistrar 使用已注册的工厂创建 registrar，找不到时返回错误
func NewRegistrar(cfg *conf.Registry) (registry.Registrar, error) {
	if cfg == nil {
		return nil, nil
	}

	if cfg.GetType() == "" {
		return nil, nil
	}

	name := Type(cfg.GetType())

	f, ok := GetRegistrarFactory(name)
	if !ok {
		return nil, fmt.Errorf("registry: registrar factory %q not found", name)
	}
	return f(cfg)
}

// RegisterRegistrarFactory 注册一个名为 name 的 RegistrarFactory，已存在时返回错误
func RegisterRegistrarFactory(name Type, f RegistrarFactory) error {
	if name == "" {
		return fmt.Errorf("registry: registrar name is empty")
	}
	if f == nil {
		return fmt.Errorf("registry: registrar factory is nil")
	}
	muRF.Lock()
	defer muRF.Unlock()
	if _, exists := registrarFactories[name]; exists {
		return fmt.Errorf("registry: registrar %q already registered", name)
	}
	registrarFactories[name] = f
	return nil
}

// MustRegisterRegistrarFactory 在注册失败时 panic（便于 init 中使用）
func MustRegisterRegistrarFactory(name Type, f RegistrarFactory) {
	if err := RegisterRegistrarFactory(name, f); err != nil {
		panic(err)
	}
}

// GetRegistrarFactory 返回名为 name 的 RegistrarFactory 及是否存在
func GetRegistrarFactory(name Type) (RegistrarFactory, bool) {
	muRF.RLock()
	defer muRF.RUnlock()
	f, ok := registrarFactories[name]
	return f, ok
}

// ListRegistrarFactories 返回已注册的 registrar 名称有序列表
func ListRegistrarFactories() []string {
	muRF.RLock()
	defer muRF.RUnlock()
	names := make([]string, 0, len(registrarFactories))
	for k := range registrarFactories {
		names = append(names, string(k))
	}
	sort.Strings(names)
	return names
}
