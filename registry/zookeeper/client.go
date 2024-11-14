package zookeeper

import (
	zookeeperKratos "github.com/go-kratos/kratos/contrib/registry/zookeeper/v2"
	"github.com/go-kratos/kratos/v2/log"

	"github.com/go-zookeeper/zk"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewRegistry 创建一个注册发现客户端 - ZooKeeper
func NewRegistry(c *conf.Registry) *zookeeperKratos.Registry {
	if c == nil || c.Zookeeper == nil {
		return nil
	}

	conn, _, err := zk.Connect(c.Zookeeper.Endpoints, c.Zookeeper.Timeout.AsDuration())
	if err != nil {
		log.Fatal(err)
	}

	reg := zookeeperKratos.New(conn)
	if err != nil {
		log.Fatal(err)
	}

	return reg
}
