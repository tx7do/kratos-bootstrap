package config

import conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"

var configList []interface{}

var commonConfig = &conf.Bootstrap{}

func GetBootstrapConfig() *conf.Bootstrap {
	return commonConfig
}

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
