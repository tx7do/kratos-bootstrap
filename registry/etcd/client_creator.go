package etcd

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	etcdClient "go.etcd.io/etcd/client/v3"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	r "github.com/tx7do/kratos-bootstrap/registry"
)

func init() {
	r.RegisterRegistrarCreator(string(r.Etcd), func(c *conf.Registry) registry.Registrar {
		return NewRegistry(c)
	})
	r.RegisterDiscoveryCreator(string(r.Etcd), func(c *conf.Registry) registry.Discovery {
		return NewRegistry(c)
	})
}

// NewRegistry 创建一个注册发现客户端 - Etcd
func NewRegistry(c *conf.Registry) *Registry {
	if c == nil || c.Etcd == nil {
		return nil
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

	return reg
}
