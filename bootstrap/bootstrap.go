package bootstrap

import (
	"context"
	"fmt"
	"runtime"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	kratosRegistry "github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"

	bConfig "github.com/tx7do/kratos-bootstrap/config"
	bLogger "github.com/tx7do/kratos-bootstrap/logger"
	bRegistry "github.com/tx7do/kratos-bootstrap/registry"
	"github.com/tx7do/kratos-bootstrap/tracer"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

var (
	AppInfo = NewAppInfo(
		"",
		"1.0.0",
		"",
	)
)

// NewApp 创建应用程序
func NewApp(ll log.Logger, rr kratosRegistry.Registrar, srv ...transport.Server) *kratos.App {
	return kratos.New(
		kratos.Context(context.Background()),
		kratos.ID(AppInfo.GetInstanceId()),
		kratos.Name(AppInfo.Name),
		kratos.Version(AppInfo.Version),
		kratos.Metadata(AppInfo.Metadata),
		kratos.Logger(ll),
		kratos.Server(
			srv...,
		),
		kratos.Registrar(rr),
	)
}

// DoBootstrap 执行引导
func DoBootstrap(appInfo *conf.AppInfo) (*conf.Bootstrap, log.Logger, kratosRegistry.Registrar) {
	// inject command flags
	Flags := NewCommandFlags()
	Flags.Init()

	var err error

	// load configs
	if err = bConfig.LoadBootstrapConfig(Flags.Conf); err != nil {
		panic(fmt.Sprintf("load config failed: %v", err))
	}

	bootstrapCfg := bConfig.GetBootstrapConfig()
	if bootstrapCfg == nil {
		panic("bootstrap config is nil")
	}

	// init logger
	ll := bLogger.NewLoggerProvider(bootstrapCfg.Logger, appInfo)

	// init registrar
	reg, err := bRegistry.NewRegistrar(bootstrapCfg.Registry)
	if err != nil {
		panic(fmt.Sprintf("init registrar failed: %v", err))
		return nil, nil, nil
	}

	// init tracer
	if err = tracer.NewTracerProvider(context.Background(), bootstrapCfg.Trace, appInfo); err != nil {
		panic(fmt.Sprintf("init tracer failed: %v", err))
	}

	return bootstrapCfg, ll, reg
}

type InitApp func(logger log.Logger, registrar kratosRegistry.Registrar, bootstrap *conf.Bootstrap) (*kratos.App, func(), error)

// Bootstrap 应用引导启动
func Bootstrap(initApp InitApp, serviceName, version *string) {
	if serviceName != nil && len(*serviceName) != 0 {
		AppInfo.Name = *serviceName
	}
	if version != nil && len(*version) != 0 {
		AppInfo.Version = *version
	}

	// bootstrap
	cfg, ll, reg := DoBootstrap(AppInfo)

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
