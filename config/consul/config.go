package consul

import (
	"strings"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"

	consulApi "github.com/hashicorp/consul/api"

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

// NewConfigSource 创建一个远程配置源 - Consul
func NewConfigSource(c *conf.RemoteConfig) config.Source {
	cfg := consulApi.DefaultConfig()
	cfg.Address = c.Consul.Address
	cfg.Scheme = c.Consul.Scheme

	cli, err := consulApi.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	src, err := New(cli,
		WithPath(getConfigKey(c.Consul.Key, true)),
	)
	if err != nil {
		log.Fatal(err)
	}

	return src
}
