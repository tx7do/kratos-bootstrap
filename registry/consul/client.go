package consul

import (
	"github.com/go-kratos/kratos/v2/log"

	consulKratos "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	consulClient "github.com/hashicorp/consul/api"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewRegistry 创建一个注册发现客户端 - Consul
func NewRegistry(c *conf.Registry) *consulKratos.Registry {
	cfg := consulClient.DefaultConfig()
	cfg.Address = c.Consul.GetAddress()
	cfg.Scheme = c.Consul.GetScheme()

	var cli *consulClient.Client
	var err error
	if cli, err = consulClient.NewClient(cfg); err != nil {
		log.Fatal(err)
	}

	reg := consulKratos.New(cli, consulKratos.WithHealthCheck(c.Consul.GetHealthCheck()))

	return reg
}
