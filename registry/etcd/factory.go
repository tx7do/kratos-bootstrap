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

	in := c.Etcd

	cfg := etcdClient.Config{
		Endpoints: in.Endpoints,

		Username:             in.GetUsername(),
		Password:             in.GetPassword(),
		AutoSyncInterval:     in.GetAutoSyncInterval().AsDuration(),
		DialTimeout:          in.GetDialTimeout().AsDuration(),
		DialKeepAliveTime:    in.GetDialKeepAliveTime().AsDuration(),
		DialKeepAliveTimeout: in.GetDialKeepAliveTimeout().AsDuration(),
		RejectOldCluster:     in.GetRejectOldCluster(),
		PermitWithoutStream:  in.GetPermitWithoutStream(),
	}

	if tlsConf := in.GetTls(); tlsConf != nil {
		tlsCfg, err := baseRegistry.LoadClientTlsConfig(tlsConf)
		if err != nil {
			log.Errorf("%v", err)
			return nil, err
		}
		cfg.TLS = tlsCfg
	}

	var err error
	var cli *etcdClient.Client
	if cli, err = etcdClient.New(cfg); err != nil {
		log.Errorf("%v", err)
		return nil, err
	}

	var opts []Option
	if in.GetNamespace() != "" {
		opts = append(opts, Namespace(in.GetNamespace()))
	}
	if in.GetRegisterTtl().AsDuration() > 0 {
		opts = append(opts, RegisterTTL(in.GetRegisterTtl().AsDuration()))
	}
	if in.GetMaxRetry() > 0 {
		opts = append(opts, MaxRetry(int(in.GetMaxRetry())))
	}

	reg := New(cli, opts...)

	return reg, nil
}

func NewDiscovery(c *conf.Registry) (registry.Discovery, error) {
	return NewRegistry(c)
}

func NewRegistrar(c *conf.Registry) (registry.Registrar, error) {
	return NewRegistry(c)
}
