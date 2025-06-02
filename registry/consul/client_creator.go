package consul

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	consulClient "github.com/hashicorp/consul/api"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	r "github.com/tx7do/kratos-bootstrap/registry"
)

func init() {
	r.RegisterRegistrarCreator(string(r.Consul), func(c *conf.Registry) registry.Registrar {
		return NewRegistry(c)
	})
	r.RegisterDiscoveryCreator(string(r.Consul), func(c *conf.Registry) registry.Discovery {
		return NewRegistry(c)
	})
}

// NewRegistry 创建一个注册发现客户端 - Consul
func NewRegistry(c *conf.Registry) *Registry {
	if c == nil || c.Consul == nil {
		return nil
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

	return reg
}
