package etcd

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"

	"google.golang.org/grpc"

	etcdClient "go.etcd.io/etcd/client/v3"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	bConfig "github.com/tx7do/kratos-bootstrap/config"
)

func init() {
	bConfig.MustRegisterFactory(bConfig.TypeEtcd, NewConfigSource)
}

// NewConfigSource 创建一个远程配置源 - Etcd
func NewConfigSource(c *conf.RemoteConfig) (config.Source, error) {
	if c == nil || c.Etcd == nil {
		return nil, nil
	}

	cfg := etcdClient.Config{
		Endpoints:   c.Etcd.Endpoints,
		DialTimeout: c.Etcd.Timeout.AsDuration(),
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	}

	cli, err := etcdClient.New(cfg)
	if err != nil {
		return nil, err
	}

	src, err := New(cli, WithPath(getConfigKey(c.Etcd.Key, true)))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return src, nil
}
