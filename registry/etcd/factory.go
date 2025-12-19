package etcd

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	etcdClient "go.etcd.io/etcd/client/v3"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	baseRegistry "github.com/tx7do/kratos-bootstrap/registry"
)

func init() {
	_ = baseRegistry.RegisterDiscoveryFactory(baseRegistry.Etcd, NewDiscovery)
	_ = baseRegistry.RegisterRegistrarFactory(baseRegistry.Etcd, NewRegistrar)
}

// NewRegistry 创建一个注册发现客户端 - Etcd
func NewRegistry(c *conf.Registry) (*Registry, error) {
	if c == nil || c.Etcd == nil {
		return nil, nil
	}

	cfg := etcdClient.Config{
		Endpoints: c.Etcd.Endpoints,
	}

	var err error
	var cli *etcdClient.Client
	if cli, err = etcdClient.New(cfg); err != nil {
		log.Fatal(err)
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
