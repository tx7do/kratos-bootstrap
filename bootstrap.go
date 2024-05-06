package bootstrap

import (
	"fmt"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	kratosRegistry "github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/tx7do/kratos-bootstrap/config"
	"github.com/tx7do/kratos-bootstrap/logger"
	"github.com/tx7do/kratos-bootstrap/registry"
	"github.com/tx7do/kratos-bootstrap/tracer"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

var (
	Service = config.NewServiceInfo(
		"",
		"1.0.0",
		"",
	)
)

func NewApp(ll log.Logger, rr kratosRegistry.Registrar, srv ...transport.Server) *kratos.App {
	return kratos.New(
		kratos.ID(Service.GetInstanceId()),
		kratos.Name(Service.Name),
		kratos.Version(Service.Version),
		kratos.Metadata(Service.Metadata),
		kratos.Logger(ll),
		kratos.Server(
			srv...,
		),
		kratos.Registrar(rr),
	)
}

// doBootstrap 应用引导启动
func doBootstrap(serviceInfo *config.ServiceInfo) (*conf.Bootstrap, log.Logger, kratosRegistry.Registrar) {
	// inject command flags
	Flags := config.NewCommandFlags()
	Flags.Init()

	var err error

	// load configs
	if err = config.LoadBootstrapConfig(Flags.Conf); err != nil {
		panic(fmt.Sprintf("load config failed: %v", err))
	}

	// init logger
	ll := logger.NewLoggerProvider(config.GetBootstrapConfig().Logger, serviceInfo)

	// init registrar
	reg := registry.NewRegistry(config.GetBootstrapConfig().Registry)

	// init tracer
	if err = tracer.NewTracerProvider(config.GetBootstrapConfig().Trace, serviceInfo); err != nil {
		panic(fmt.Sprintf("init tracer failed: %v", err))
	}

	return config.GetBootstrapConfig(), ll, reg
}

type InitApp func(logger log.Logger, registrar kratosRegistry.Registrar, bootstrap *conf.Bootstrap) (*kratos.App, func(), error)

func Bootstrap(initApp InitApp) {
	// bootstrap
	cfg, ll, reg := doBootstrap(Service)

	app, cleanup, err := initApp(ll, reg, cfg)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err = app.Run(); err != nil {
		panic(err)
	}
}
