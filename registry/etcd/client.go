package etcd

import (
	"github.com/go-kratos/kratos/v2/log"

	etcdKratos "github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	etcdClient "go.etcd.io/etcd/client/v3"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewRegistry 创建一个注册发现客户端 - Etcd
func NewRegistry(c *conf.Registry) *etcdKratos.Registry {
	cfg := etcdClient.Config{
		Endpoints: c.Etcd.Endpoints,
	}

	var err error
	var cli *etcdClient.Client
	if cli, err = etcdClient.New(cfg); err != nil {
		log.Fatal(err)
	}

	reg := etcdKratos.New(cli)

	return reg
}
