package zookeeper

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	"github.com/go-zookeeper/zk"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	baseRegistry "github.com/tx7do/kratos-bootstrap/registry"
)

func init() {
	_ = baseRegistry.RegisterDiscoveryFactory(baseRegistry.ZooKeeper, NewDiscovery)
	_ = baseRegistry.RegisterRegistrarFactory(baseRegistry.ZooKeeper, NewRegistrar)
}

// NewRegistry 创建一个注册发现客户端 - ZooKeeper
func NewRegistry(c *conf.Registry) (*Registry, error) {
	if c == nil || c.Zookeeper == nil {
		return nil, nil
	}

	conn, _, err := zk.Connect(c.Zookeeper.Endpoints, c.Zookeeper.Timeout.AsDuration())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	reg := New(conn)

	return reg, nil
}

func NewDiscovery(c *conf.Registry) (registry.Discovery, error) {
	return NewRegistry(c)
}

func NewRegistrar(c *conf.Registry) (registry.Registrar, error) {
	return NewRegistry(c)
}
