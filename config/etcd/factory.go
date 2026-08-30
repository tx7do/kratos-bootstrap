package etcd

import (
	"crypto/tls"

	"time"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"google.golang.org/grpc"

	etcdClient "go.etcd.io/etcd/client/v3"

	tlsUtils "github.com/tx7do/go-utils/tls"

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

	in := c.Etcd

	keepAliveTime := 30 * time.Second
	if in.GetDialKeepAliveTime().AsDuration() > 0 {
		keepAliveTime = in.GetDialKeepAliveTime().AsDuration()
	}
	keepAliveTimeout := 10 * time.Second
	if in.GetDialKeepAliveTimeout().AsDuration() > 0 {
		keepAliveTimeout = in.GetDialKeepAliveTimeout().AsDuration()
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
			Time:                keepAliveTime,
			Timeout:             keepAliveTimeout,
			PermitWithoutStream: true,
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	cfg := etcdClient.Config{
		Endpoints:   in.Endpoints,
		DialTimeout: in.GetTimeout().AsDuration(),
		Username:    in.GetUsername(),
		Password:    in.GetPassword(),
		DialOptions: dialOpts,
	}

	if tlsConf := in.GetTls(); tlsConf != nil {
		tlsClientCfg, err := loadClientTlsConfig(tlsConf)
		if err != nil {
			log.Errorf("%v", err)
			return nil, err
		}
		cfg.TLS = tlsClientCfg
		if tlsClientCfg != nil {
			// TLS模式下替换掉默认的insecure凭证
			cfg.DialOptions = dialOpts[:len(dialOpts)-1]
		}
	}

	cli, err := etcdClient.New(cfg)
	if err != nil {
		return nil, err
	}

	src, err := New(cli, WithPath(getConfigKey(c.Etcd.Key, true)))
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}

	return src, nil
}

func loadClientTlsConfig(cfg *conf.TLS) (*tls.Config, error) {
	if cfg == nil {
		return nil, nil
	}

	var tlsCfg *tls.Config
	var err error

	if cfg.File != nil {
		if tlsCfg, err = tlsUtils.LoadClientTlsConfigFile(
			cfg.File.GetKeyPath(),
			cfg.File.GetCertPath(),
			cfg.File.GetCaPath(),
		); err != nil {
			return nil, err
		}
	} else if cfg.Config != nil {
		if tlsCfg, err = tlsUtils.LoadClientTlsConfigString(
			cfg.Config.GetKeyPem(),
			cfg.Config.GetCertPem(),
			cfg.Config.GetCaPem(),
		); err != nil {
			return nil, err
		}
	}

	if tlsCfg != nil && cfg.GetInsecureSkipVerify() {
		tlsCfg.InsecureSkipVerify = true
	}

	return tlsCfg, nil
}
