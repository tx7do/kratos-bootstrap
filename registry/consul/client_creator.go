package consul

import (
	"github.com/go-kratos/kratos/v2/log"
	consulClient "github.com/hashicorp/consul/api"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func init() {
	// 注册 Consul 注册发现客户端
	//registry.RegisterDiscoveryCreator(registry.Consul, NewRegistry)
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
