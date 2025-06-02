package config

import (
	"github.com/go-kratos/kratos/v2/config"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"

	"github.com/tx7do/kratos-bootstrap/remoteconfig/apollo"
	"github.com/tx7do/kratos-bootstrap/remoteconfig/consul"
	"github.com/tx7do/kratos-bootstrap/remoteconfig/etcd"
	"github.com/tx7do/kratos-bootstrap/remoteconfig/kubernetes"
	"github.com/tx7do/kratos-bootstrap/remoteconfig/nacos"
	"github.com/tx7do/kratos-bootstrap/remoteconfig/polaris"
)

const remoteConfigSourceConfigFile = "remote.yaml"

type Type string

const (
	LocalFile  Type = "file"
	Nacos      Type = "nacos"
	Consul     Type = "consul"
	Etcd       Type = "etcd"
	Apollo     Type = "apollo"
	Kubernetes Type = "kubernetes"
	Polaris    Type = "polaris"
)

// NewRemoteConfigSource 创建一个远程配置源
func NewRemoteConfigSource(c *conf.RemoteConfig) config.Source {
	switch Type(c.Type) {
	default:
		fallthrough
	case LocalFile:
		return nil
	case Nacos:
		return nacos.NewConfigSource(c)
	case Consul:
		return consul.NewConfigSource(c)
	case Etcd:
		return etcd.NewConfigSource(c)
	case Apollo:
		return apollo.NewConfigSource(c)
	case Kubernetes:
		return kubernetes.NewConfigSource(c)
	case Polaris:
		return polaris.NewConfigSource(c)
	}
}
