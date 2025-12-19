package servicecomb

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	servicecombClient "github.com/go-chassis/sc-client"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	baseRegistry "github.com/tx7do/kratos-bootstrap/registry"
)

func init() {
	_ = baseRegistry.RegisterDiscoveryFactory(baseRegistry.Servicecomb, NewDiscovery)
	_ = baseRegistry.RegisterRegistrarFactory(baseRegistry.Servicecomb, NewRegistrar)
}

// NewRegistry 创建一个注册发现客户端 - Servicecomb
func NewRegistry(c *conf.Registry) (*Registry, error) {
	if c == nil || c.Servicecomb == nil {
		return nil, nil
	}

	cfg := servicecombClient.Options{
		Endpoints: c.Servicecomb.Endpoints,
	}

	var cli *servicecombClient.Client
	var err error
	if cli, err = servicecombClient.NewClient(cfg); err != nil {
		log.Fatal(err)
		return nil, err
	}

	reg := New(cli)

	return reg, nil
}

func NewDiscovery(c *conf.Registry) (registry.Discovery, error) {
	return NewRegistry(c)
}

func NewRegistrar(c *conf.Registry) (registry.Registrar, error) {
	return NewRegistry(c)
}
