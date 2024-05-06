package bootstrap

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"

	"github.com/tx7do/kratos-bootstrap/config/apollo"
	"github.com/tx7do/kratos-bootstrap/config/consul"
	"github.com/tx7do/kratos-bootstrap/config/etcd"
	"github.com/tx7do/kratos-bootstrap/config/kubernetes"
	"github.com/tx7do/kratos-bootstrap/config/nacos"
	"github.com/tx7do/kratos-bootstrap/config/polaris"

	// file
	fileKratos "github.com/go-kratos/kratos/v2/config/file"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

var commonConfig = &conf.Bootstrap{}
var configList []interface{}

const remoteConfigSourceConfigFile = "remote.yaml"

// RegisterConfig 注册配置
func RegisterConfig(c interface{}) {
	initBootstrapConfig()
	configList = append(configList, c)
}

func initBootstrapConfig() {
	if len(configList) > 0 {
		return
	}

	configList = append(configList, commonConfig)

	if commonConfig.Server == nil {
		commonConfig.Server = &conf.Server{}
		configList = append(configList, commonConfig.Server)
	}

	if commonConfig.Client == nil {
		commonConfig.Client = &conf.Client{}
		configList = append(configList, commonConfig.Client)
	}

	if commonConfig.Data == nil {
		commonConfig.Data = &conf.Data{}
		configList = append(configList, commonConfig.Data)
	}

	if commonConfig.Trace == nil {
		commonConfig.Trace = &conf.Tracer{}
		configList = append(configList, commonConfig.Trace)
	}

	if commonConfig.Logger == nil {
		commonConfig.Logger = &conf.Logger{}
		configList = append(configList, commonConfig.Logger)
	}

	if commonConfig.Registry == nil {
		commonConfig.Registry = &conf.Registry{}
		configList = append(configList, commonConfig.Registry)
	}

	if commonConfig.Oss == nil {
		commonConfig.Oss = &conf.OSS{}
		configList = append(configList, commonConfig.Oss)
	}

	if commonConfig.Notify == nil {
		commonConfig.Notify = &conf.Notification{}
		configList = append(configList, commonConfig.Notify)
	}
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

	return nil, commonConfig.Config
}

type ConfigType string

const (
	ConfigTypeLocalFile  ConfigType = "file"
	ConfigTypeNacos      ConfigType = "nacos"
	ConfigTypeConsul     ConfigType = "consul"
	ConfigTypeEtcd       ConfigType = "etcd"
	ConfigTypeApollo     ConfigType = "apollo"
	ConfigTypeKubernetes ConfigType = "kubernetes"
	ConfigTypePolaris    ConfigType = "polaris"
)

// NewRemoteConfigSource 创建一个远程配置源
func NewRemoteConfigSource(c *conf.RemoteConfig) config.Source {
	switch ConfigType(c.Type) {
	default:
		fallthrough
	case ConfigTypeLocalFile:
		return nil
	case ConfigTypeNacos:
		return nacos.NewConfigSource(c)
	case ConfigTypeConsul:
		return consul.NewConfigSource(c)
	case ConfigTypeEtcd:
		return etcd.NewConfigSource(c)
	case ConfigTypeApollo:
		return apollo.NewConfigSource(c)
	case ConfigTypeKubernetes:
		return kubernetes.NewConfigSource(c)
	case ConfigTypePolaris:
		return polaris.NewConfigSource(c)
	}
}

// NewFileConfigSource 创建一个本地文件配置源
func NewFileConfigSource(filePath string) config.Source {
	return fileKratos.NewSource(filePath)
}
