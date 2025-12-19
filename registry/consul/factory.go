package consul

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	consulClient "github.com/hashicorp/consul/api"

	baseRegistry "github.com/tx7do/kratos-bootstrap/registry"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func init() {
	_ = baseRegistry.RegisterDiscoveryFactory(baseRegistry.Consul, NewDiscovery)
	_ = baseRegistry.RegisterRegistrarFactory(baseRegistry.Consul, NewRegistrar)
}

// NewRegistry 创建一个注册发现客户端 - Consul
func NewRegistry(c *conf.Registry) (*Registry, error) {
	if c == nil || c.Consul == nil {
		return nil, nil
	}

	cfg := consulClient.DefaultConfig()
	cfg.Address = c.Consul.GetAddress()
	cfg.Scheme = c.Consul.GetScheme()

	var cli *consulClient.Client
	var err error
	if cli, err = consulClient.NewClient(cfg); err != nil {
		log.Fatal(err)
	}

	reg := New(cli, WithHealthCheck(c.Consul.GetHealthCheck()))

	return reg, nil
}

func NewDiscovery(c *conf.Registry) (registry.Discovery, error) {
	return NewRegistry(c)
}

func NewRegistrar(c *conf.Registry) (registry.Registrar, error) {
	return NewRegistry(c)
}
