package logger

import (
	"fmt"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

type FactoryFunc func(cfg *conf.Logger) (log.Logger, error)

var (
	factoryMu sync.RWMutex
	factories = make(map[Type]FactoryFunc)
)

// Register 注册远程配置客户端构造函数
func Register(typ Type, f FactoryFunc) error {
	factoryMu.Lock()
	defer factoryMu.Unlock()
	if _, ok := factories[typ]; ok {
		return fmt.Errorf("logger factory %s already registered", typ)
	}
	factories[typ] = f
	return nil
}

// GetFactory returns a registered FactoryFunc for a given Type and whether it existed.
// Safe for concurrent use.
func GetFactory(typ Type) (FactoryFunc, bool) {
	factoryMu.RLock()
	defer factoryMu.RUnlock()
	f, ok := factories[typ]
	return f, ok
}

// ListFactories returns a slice of currently registered Types.
func ListFactories() []Type {
	factoryMu.RLock()
	defer factoryMu.RUnlock()
	res := make([]Type, 0, len(factories))
	for k := range factories {
		res = append(res, k)
	}
	return res
}

// Unregister removes a registered factory by Type. It returns true if a factory was removed.
// Use with caution in concurrent environments (primarily intended for tests).
func Unregister(typ Type) bool {
	factoryMu.Lock()
	defer factoryMu.Unlock()
	if _, ok := factories[typ]; ok {
		delete(factories, typ)
		return true
	}
	return false
}
