package zookeeper

import (
	"time"

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

	in := c.Zookeeper

	timeout := in.GetTimeout().AsDuration()
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	conn, _, err := zk.Connect(in.Endpoints, timeout)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}

	var opts []Option
	if in.GetNamespace() != "" {
		opts = append(opts, WithRootPath(in.GetNamespace()))
	}
	if in.GetUsername() != "" && in.GetPassword() != "" {
		opts = append(opts, WithDigestACL(in.GetUsername(), in.GetPassword()))
	}

	reg := New(conn, opts...)

	return reg, nil
}

func NewDiscovery(c *conf.Registry) (registry.Discovery, error) {
	return NewRegistry(c)
}

func NewRegistrar(c *conf.Registry) (registry.Registrar, error) {
	return NewRegistry(c)
}
