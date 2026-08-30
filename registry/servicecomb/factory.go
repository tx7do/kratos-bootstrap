package servicecomb

import (
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	servicecombRbac "github.com/go-chassis/cari/rbac"
	servicecombClient "github.com/go-chassis/sc-client"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	baseRegistry "github.com/tx7do/kratos-bootstrap/registry"
)

func init() {
	_ = baseRegistry.RegisterDiscoveryFactory(baseRegistry.Servicecomb, NewDiscovery)
	_ = baseRegistry.RegisterRegistrarFactory(baseRegistry.Servicecomb, NewRegistrar)
}

// NewRegistry 创建一个注册发现客户端 - Servicecomb
func NewRegistry(c *conf.Registry) (*Registry, error) {
	if c == nil || c.Servicecomb == nil {
		return nil, nil
	}

	in := c.Servicecomb

	cfg := servicecombClient.Options{
		Endpoints:  in.Endpoints,
		EnableSSL:  in.GetEnableSsl(),
		Timeout:    in.GetTimeout().AsDuration(),
		Compressed: in.GetCompressed(),
		Verbose:    in.GetVerbose(),
	}

	if in.GetTimeout().AsDuration() <= 0 {
		cfg.Timeout = 10 * time.Second
	}

	if in.GetEnableAuth() {
		cfg.EnableAuth = true
		cfg.AuthUser = &servicecombRbac.AuthUser{
			Username: in.GetUsername(),
			Password: in.GetPassword(),
		}
		if in.GetTokenExpiration().AsDuration() > 0 {
			cfg.TokenExpiration = in.GetTokenExpiration().AsDuration()
		}
	}

	if tlsConf := in.GetTls(); tlsConf != nil {
		tlsCfg, err := baseRegistry.LoadClientTlsConfig(tlsConf)
		if err != nil {
			log.Errorf("%v", err)
			return nil, err
		}
		cfg.TLSConfig = tlsCfg
		cfg.EnableSSL = true
	}

	var cli *servicecombClient.Client
	var err error
	if cli, err = servicecombClient.NewClient(cfg); err != nil {
		log.Errorf("%v", err)
		return nil, err
	}

	reg := New(cli)

	return reg, nil
}

func NewDiscovery(c *conf.Registry) (registry.Discovery, error) {
	return NewRegistry(c)
}

func NewRegistrar(c *conf.Registry) (registry.Registrar, error) {
	return NewRegistry(c)
}
