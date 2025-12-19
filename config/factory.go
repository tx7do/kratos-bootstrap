package config

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/go-kratos/kratos/v2/config"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// Factory 根据传入的配置创建一个 Provider 实例。
type Factory func(cfg *conf.RemoteConfig) (config.Source, error)

var (
	mu        sync.RWMutex
	factories = make(map[Type]Factory)
)

// NewProvider 使用指定 name 的 Factory 和 cfg 创建 Provider 实例。
func NewProvider(cfg *conf.RemoteConfig) (config.Source, error) {
	if cfg == nil {
		return nil, errors.New("config: config is nil")
	}

	name := Type(cfg.GetType())
	if name == "" {
		return nil, errors.New("config: factory name is empty")
	}

	f, ok := GetFactory(name)
	if !ok {
		return nil, fmt.Errorf("config: factory %q not found", name)
	}

	p, err := f(cfg)
	if err != nil {
		return nil, fmt.Errorf("config: factory %q create failed: %w", name, err)
	}
	if p == nil {
		return nil, fmt.Errorf("config: factory %q returned nil provider", name)
	}
	return p, nil
}

// RegisterFactory 注册一个名为 name 的 Factory。
// name 不能为空且不可重复；f 不能为 nil。
func RegisterFactory(name Type, f Factory) error {
	if name == "" {
		return errors.New("config: factory name is empty")
	}
	if f == nil {
		return errors.New("config: factory is nil")
	}

	mu.Lock()
	defer mu.Unlock()

	if _, ok := factories[name]; ok {
		return fmt.Errorf("config: factory %q already registered", name)
	}
	factories[name] = f
	return nil
}

// MustRegisterFactory 等同于 RegisterFactory，但在出错时 panic（适用于 init）。
func MustRegisterFactory(name Type, f Factory) {
	if err := RegisterFactory(name, f); err != nil {
		panic(err)
	}
}

// GetFactory 返回已注册的 Factory 及其存在标识。
func GetFactory(name Type) (Factory, bool) {
	mu.RLock()
	defer mu.RUnlock()
	f, ok := factories[name]
	return f, ok
}

// ListFactories 返回已注册工厂名称的字母序列表。
func ListFactories() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(factories))
	for n := range factories {
		names = append(names, string(n))
	}
	sort.Strings(names)
	return names
}
