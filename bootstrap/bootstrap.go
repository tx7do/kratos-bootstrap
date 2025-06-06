package bootstrap

import (
	"fmt"
	"runtime"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	kratosRegistry "github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/tx7do/kratos-bootstrap/logger"
	"github.com/tx7do/kratos-bootstrap/tracer"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/utils"
)

var (
	Service = utils.NewServiceInfo(
		"",
		"1.0.0",
		"",
	)
)

// NewApp 创建应用程序
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

// DoBootstrap 执行引导
func DoBootstrap(serviceInfo *utils.ServiceInfo) (*conf.Bootstrap, log.Logger, kratosRegistry.Registrar) {
	// inject command flags
	Flags := NewCommandFlags()
	Flags.Init()

	var err error

	// load configs
	if err = LoadBootstrapConfig(Flags.Conf); err != nil {
		panic(fmt.Sprintf("load config failed: %v", err))
	}

	// init logger
	ll := logger.NewLoggerProvider(GetBootstrapConfig().Logger, serviceInfo)

	// init registrar
	reg := NewRegistry(GetBootstrapConfig().Registry)

	// init tracer
	if err = tracer.NewTracerProvider(GetBootstrapConfig().Trace, serviceInfo); err != nil {
		panic(fmt.Sprintf("init tracer failed: %v", err))
	}

	return GetBootstrapConfig(), ll, reg
}

type InitApp func(logger log.Logger, registrar kratosRegistry.Registrar, bootstrap *conf.Bootstrap) (*kratos.App, func(), error)

// Bootstrap 应用引导启动
func Bootstrap(initApp InitApp, serviceName, version *string) {
	if serviceName != nil && len(*serviceName) != 0 {
		Service.Name = *serviceName
	}
	if version != nil && len(*version) != 0 {
		Service.Version = *version
	}

	// bootstrap
	cfg, ll, reg := DoBootstrap(Service)

	// init app
	app, cleanup, err := initApp(ll, reg, cfg)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// run the app.
	if err = app.Run(); err != nil {
		buf := make([]byte, 1024)
		n := runtime.Stack(buf, false)
		panic(fmt.Sprintf("Panic: %v\nStack trace:\n%s", err, string(buf[:n])))
	}
}
