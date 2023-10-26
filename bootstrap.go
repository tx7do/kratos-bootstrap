package bootstrap

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	"github.com/tx7do/kratos-bootstrap/gen/api/go/conf/v1"
)

// Bootstrap 应用引导启动
func Bootstrap(serviceInfo *ServiceInfo) (*conf.Bootstrap, log.Logger, registry.Registrar) {
	// inject command flags
	Flags := NewCommandFlags()
	Flags.Init()

	var err error

	// load configs
	if err = LoadBootstrapConfig(Flags.Conf); err == nil {
		panic("load config failed")
	}

	// init logger
	ll := NewLoggerProvider(commonConfig.Logger, serviceInfo)

	// init registrar
	reg := NewRegistry(commonConfig.Registry)

	// init tracer
	if err = NewTracerProvider(commonConfig.Trace, serviceInfo); err != nil {
		panic(err)
	}

	return commonConfig, ll, reg
}
