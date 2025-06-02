package etcd

import (
	"github.com/go-kratos/kratos/v2/log"

	etcdClient "go.etcd.io/etcd/client/v3"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

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
