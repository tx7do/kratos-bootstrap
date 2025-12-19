package config

import (
	"os"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/config"
	fileKratos "github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

const remoteConfigSourceConfigFile = "config.yaml"

// NewFileConfigSource 创建一个本地文件配置源
func NewFileConfigSource(filePath string) config.Source {
	return fileKratos.NewSource(filePath)
}

// NewConfigProvider 创建一个配置
func NewConfigProvider(configPath string) (config.Config, error) {
	err, rc := LoadRemoteConfigSourceConfigs(configPath)
	if err != nil {
		log.Error("LoadRemoteConfigSourceConfigs: ", err.Error())
		return nil, err
	}

	if rc != nil {
		rcs, err := NewProvider(rc)
		if err != nil {
			log.Error("NewProvider: ", err.Error())
			return nil, err
		}
		return config.New(
			config.WithSource(
				NewFileConfigSource(configPath),
				rcs,
			),
		), nil
	}

	return config.New(
		config.WithSource(
			NewFileConfigSource(configPath),
		),
	), nil
}

// LoadBootstrapConfig 加载程序引导配置
func LoadBootstrapConfig(configPath string) error {
	cfg, err := NewConfigProvider(configPath)
	if err != nil {
		return err
	}

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
	configPath = filepath.Join(configPath, remoteConfigSourceConfigFile)
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

// pathExists 判断路径是否存在
func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
