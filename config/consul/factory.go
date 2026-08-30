package consul

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"

	consulApi "github.com/hashicorp/consul/api"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	bConfig "github.com/tx7do/kratos-bootstrap/config"
)

func init() {
	bConfig.MustRegisterFactory(bConfig.TypeConsul, NewConfigSource)
}

// NewConfigSource 创建一个远程配置源 - Consul
func NewConfigSource(c *conf.RemoteConfig) (config.Source, error) {
	if c == nil || c.Consul == nil {
		return nil, nil
	}

	in := c.Consul

	cfg := consulApi.DefaultConfig()
	cfg.Address = in.GetAddress()
	cfg.Scheme = in.GetScheme()
	cfg.Datacenter = in.GetDatacenter()
	cfg.Token = in.GetToken()
	cfg.TokenFile = in.GetTokenFile()
	cfg.Namespace = in.GetNamespace()
	cfg.Partition = in.GetPartition()
	cfg.WaitTime = in.GetWaitTime().AsDuration()

	if tlsConf := in.GetTls(); tlsConf != nil {
		if in.GetScheme() == "" {
			cfg.Scheme = "https"
		}
		tlsCfg := &consulApi.TLSConfig{
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

	cli, err := consulApi.NewClient(cfg)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}

	src, err := New(cli,
		WithPath(getConfigKey(c.Consul.Key, true)),
	)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}

	return src, nil
}
