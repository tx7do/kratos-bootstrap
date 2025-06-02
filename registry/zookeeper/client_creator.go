package zookeeper

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	"github.com/go-zookeeper/zk"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	r "github.com/tx7do/kratos-bootstrap/registry"
)

func init() {
	r.RegisterRegistrarCreator(string(r.ZooKeeper), func(c *conf.Registry) registry.Registrar {
		return NewRegistry(c)
	})
	r.RegisterDiscoveryCreator(string(r.ZooKeeper), func(c *conf.Registry) registry.Discovery {
		return NewRegistry(c)
	})
}

// NewRegistry 创建一个注册发现客户端 - ZooKeeper
func NewRegistry(c *conf.Registry) *Registry {
	if c == nil || c.Zookeeper == nil {
		return nil
	}

	conn, _, err := zk.Connect(c.Zookeeper.Endpoints, c.Zookeeper.Timeout.AsDuration())
	if err != nil {
		log.Fatal(err)
	}

	reg := New(conn)
	if err != nil {
		log.Fatal(err)
	}

	return reg
}
