package eureka

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	baseRegistry "github.com/tx7do/kratos-bootstrap/registry"
)

func init() {
	_ = baseRegistry.RegisterDiscoveryFactory(baseRegistry.Eureka, NewDiscovery)
	_ = baseRegistry.RegisterRegistrarFactory(baseRegistry.Eureka, NewRegistrar)
}

// NewRegistry 创建一个注册发现客户端 - Eureka
func NewRegistry(c *conf.Registry) (*Registry, error) {
	if c == nil || c.Eureka == nil {
		return nil, nil
	}

	var opts []Option
	if c.Eureka.HeartbeatInterval != nil {
		opts = append(opts, WithHeartbeat(c.Eureka.HeartbeatInterval.AsDuration()))
	}
	if c.Eureka.RefreshInterval != nil {
		opts = append(opts, WithRefresh(c.Eureka.RefreshInterval.AsDuration()))
	}
	if c.Eureka.Path != "" {
		opts = append(opts, WithEurekaPath(c.Eureka.Path))
	}

	var err error
	var reg *Registry
	if reg, err = New(c.Eureka.Endpoints, opts...); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return reg, nil
}

func NewDiscovery(c *conf.Registry) (registry.Discovery, error) {
	return NewRegistry(c)
}

func NewRegistrar(c *conf.Registry) (registry.Registrar, error) {
	return NewRegistry(c)
}
