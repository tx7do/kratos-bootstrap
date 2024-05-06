package etcd

import (
	"strings"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"

	"google.golang.org/grpc"

	etcdKratos "github.com/go-kratos/kratos/contrib/config/etcd/v2"
	etcdClient "go.etcd.io/etcd/client/v3"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// getConfigKey 获取合法的配置名
func getConfigKey(configKey string, useBackslash bool) string {
	if useBackslash {
		return strings.Replace(configKey, `.`, `/`, -1)
	} else {
		return configKey
	}
}

// NewConfigSource 创建一个远程配置源 - Etcd
func NewConfigSource(c *conf.RemoteConfig) config.Source {
	cfg := etcdClient.Config{
		Endpoints:   c.Etcd.Endpoints,
		DialTimeout: c.Etcd.Timeout.AsDuration(),
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	}

	cli, err := etcdClient.New(cfg)
	if err != nil {
		panic(err)
	}

	source, err := etcdKratos.New(cli, etcdKratos.WithPath(getConfigKey(c.Etcd.Key, true)))
	if err != nil {
		log.Fatal(err)
	}

	return source
}
