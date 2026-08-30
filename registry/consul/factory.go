package consul

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	consulClient "github.com/hashicorp/consul/api"

	baseRegistry "github.com/tx7do/kratos-bootstrap/registry"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func init() {
	_ = baseRegistry.RegisterDiscoveryFactory(baseRegistry.Consul, NewDiscovery)
	_ = baseRegistry.RegisterRegistrarFactory(baseRegistry.Consul, NewRegistrar)
}

// NewRegistry 创建一个注册发现客户端 - Consul
func NewRegistry(c *conf.Registry) (*Registry, error) {
	if c == nil || c.Consul == nil {
		return nil, nil
	}

	in := c.Consul

	cfg := consulClient.DefaultConfig()
	cfg.Address = in.GetAddress()
	cfg.Scheme = in.GetScheme()
	cfg.Datacenter = in.GetDatacenter()
	cfg.Token = in.GetToken()
	cfg.TokenFile = in.GetTokenFile()
	cfg.Namespace = in.GetNamespace()
	cfg.Partition = in.GetPartition()
	cfg.PathPrefix = in.GetPathPrefix()
	cfg.WaitTime = in.GetWaitTime().AsDuration()

	if auth := in.GetBasicAuth(); auth != nil {
		cfg.HttpAuth = &consulClient.HttpBasicAuth{
			Username: auth.GetUsername(),
			Password: auth.GetPassword(),
		}
	}

	if tlsConf := in.GetTls(); tlsConf != nil {
		if in.GetScheme() == "" {
			cfg.Scheme = "https"
		}
		tlsCfg := &consulClient.TLSConfig{
			InsecureSkipVerify: tlsConf.GetInsecureSkipVerify(),
		}
		if file := tlsConf.GetFile(); file != nil {
			tlsCfg.CertFile = file.GetCertPath()
			tlsCfg.KeyFile = file.GetKeyPath()
			tlsCfg.CAFile = file.GetCaPath()
		}
		if pem := tlsConf.GetConfig(); pem != nil {
			tlsCfg.CertPEM = pem.GetCertPem()
			tlsCfg.KeyPEM = pem.GetKeyPem()
			tlsCfg.CAPem = pem.GetCaPem()
		}
		cfg.TLSConfig = *tlsCfg
	}

	var cli *consulClient.Client
	var err error
	if cli, err = consulClient.NewClient(cfg); err != nil {
		log.Errorf("%v", err)
		return nil, err
	}

	var opts []Option
	opts = append(opts, WithHealthCheck(in.GetHealthCheck()))

	if in.Heartbeat != nil {
		opts = append(opts, WithHeartbeat(in.GetHeartbeat()))
	}
	if in.HealthCheckInterval != nil {
		opts = append(opts, WithHealthCheckInterval(int(in.GetHealthCheckInterval())))
	}
	if in.DeregisterCriticalServiceAfter != nil {
		opts = append(opts, WithDeregisterCriticalServiceAfter(int(in.GetDeregisterCriticalServiceAfter())))
	}
	if in.Timeout != nil {
		opts = append(opts, WithTimeout(in.GetTimeout().AsDuration()))
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
