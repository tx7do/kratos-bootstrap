package config

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"

	fileKratos "github.com/go-kratos/kratos/v2/config/file"

	"github.com/tx7do/kratos-bootstrap/config/apollo"
	"github.com/tx7do/kratos-bootstrap/config/consul"
	"github.com/tx7do/kratos-bootstrap/config/etcd"
	"github.com/tx7do/kratos-bootstrap/config/kubernetes"
	"github.com/tx7do/kratos-bootstrap/config/nacos"
	"github.com/tx7do/kratos-bootstrap/config/polaris"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

const remoteConfigSourceConfigFile = "remote.yaml"

type Type string

const (
	LocalFileType  Type = "file"
	NacosType      Type = "nacos"
	ConsulType     Type = "consul"
	EtcdType       Type = "etcd"
	ApolloType     Type = "apollo"
	KubernetesType Type = "kubernetes"
	PolarisType    Type = "polaris"
)

// NewRemoteConfigSource 创建一个远程配置源
func NewRemoteConfigSource(c *conf.RemoteConfig) config.Source {
	switch Type(c.Type) {
	default:
		fallthrough
	case LocalFileType:
		return nil
	case NacosType:
		return nacos.NewConfigSource(c)
	case ConsulType:
		return consul.NewConfigSource(c)
	case EtcdType:
		return etcd.NewConfigSource(c)
	case ApolloType:
		return apollo.NewConfigSource(c)
	case KubernetesType:
		return kubernetes.NewConfigSource(c)
	case PolarisType:
		return polaris.NewConfigSource(c)
	}
}

// NewFileConfigSource 创建一个本地文件配置源
func NewFileConfigSource(filePath string) config.Source {
	return fileKratos.NewSource(filePath)
}

// NewConfigProvider 创建一个配置
func NewConfigProvider(configPath string) config.Config {
	err, rc := LoadRemoteConfigSourceConfigs(configPath)
	if err != nil {
		log.Error("LoadRemoteConfigSourceConfigs: ", err.Error())
	}
	if rc != nil {
		return config.New(
			config.WithSource(
				NewFileConfigSource(configPath),
				NewRemoteConfigSource(rc),
			),
		)
	} else {
		return config.New(
			config.WithSource(
				NewFileConfigSource(configPath),
			),
		)
	}
}

// LoadBootstrapConfig 加载程序引导配置
func LoadBootstrapConfig(configPath string) error {
	cfg := NewConfigProvider(configPath)

	var err error

	if err = cfg.Load(); err != nil {
		return err
	}

	initBootstrapConfig()

	if err = scanConfigs(cfg); err != nil {
		return err
	}

	return nil
}

func scanConfigs(cfg config.Config) error {
	initBootstrapConfig()

	for _, c := range configList {
		if err := cfg.Scan(c); err != nil {
			return err
		}
	}
	return nil
}

// LoadRemoteConfigSourceConfigs 加载远程配置源的本地配置
func LoadRemoteConfigSourceConfigs(configPath string) (error, *conf.RemoteConfig) {
	configPath = configPath + "/" + remoteConfigSourceConfigFile
	if !pathExists(configPath) {
		return nil, nil
	}

	cfg := config.New(
		config.WithSource(
			NewFileConfigSource(configPath),
		),
	)
	defer func(cfg config.Config) {
		if err := cfg.Close(); err != nil {
			panic(err)
		}
	}(cfg)

	var err error

	if err = cfg.Load(); err != nil {
		return err, nil
	}

	if err = scanConfigs(cfg); err != nil {
		return err, nil
	}

	return nil, GetBootstrapConfig().Config
}
