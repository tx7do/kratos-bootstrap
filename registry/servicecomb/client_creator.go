package servicecomb

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	servicecombClient "github.com/go-chassis/sc-client"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	r "github.com/tx7do/kratos-bootstrap/registry"
)

func init() {
	r.RegisterRegistrarCreator(string(r.Servicecomb), func(c *conf.Registry) registry.Registrar {
		return NewRegistry(c)
	})
	r.RegisterDiscoveryCreator(string(r.Servicecomb), func(c *conf.Registry) registry.Discovery {
		return NewRegistry(c)
	})
}

// NewRegistry 创建一个注册发现客户端 - Servicecomb
func NewRegistry(c *conf.Registry) *Registry {
	if c == nil || c.Servicecomb == nil {
		return nil
	}

	cfg := servicecombClient.Options{
		Endpoints: c.Servicecomb.Endpoints,
	}

	var cli *servicecombClient.Client
	var err error
	if cli, err = servicecombClient.NewClient(cfg); err != nil {
		log.Fatal(err)
	}

	reg := New(cli)

	return reg
}
