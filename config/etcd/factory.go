package etcd

import (
	"time"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

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

	dialOpts := []grpc.DialOption{
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  500 * time.Millisecond,
				Multiplier: 1.6,
				Jitter:     0.2,
				MaxDelay:   120 * time.Second,
			},
			MinConnectTimeout: 5 * time.Second,
		}),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		}),
		// 若使用 TLS，请替换为相应的 credentials.NewClientTLSFromCert(...)
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	cfg := etcdClient.Config{
		Endpoints:   c.Etcd.Endpoints,
		DialTimeout: c.Etcd.Timeout.AsDuration(),
		DialOptions: dialOpts,
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
