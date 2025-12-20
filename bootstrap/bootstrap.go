package bootstrap

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2"
	kratosLog "github.com/go-kratos/kratos/v2/log"
	kratosRegistry "github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"

	bConfig "github.com/tx7do/kratos-bootstrap/config"
	bLogger "github.com/tx7do/kratos-bootstrap/logger"
	bRegistry "github.com/tx7do/kratos-bootstrap/registry"
	"github.com/tx7do/kratos-bootstrap/tracer"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

var (
	// AppInfo 应用信息
	AppInfo = NewAppInfo(
		"",
		"1.0.0",
		"",
	)
)

// NewApp 创建应用程序
func NewApp(ll kratosLog.Logger, rr kratosRegistry.Registrar, srv ...transport.Server) *kratos.App {
	if AppInfo.InstanceId != "" {
		SetInstanceId(AppInfo, AppInfo.GetAppId(), AppInfo.GetName())
	}

	var opts []kratos.Option
	if ll != nil {
		opts = append(opts, kratos.Logger(ll))
	}
	if rr != nil {
		opts = append(opts, kratos.Registrar(rr))
	}
	if len(srv) > 0 {
		opts = append(opts, kratos.Server(srv...))
	}

	if AppInfo.Metadata != nil {
		opts = append(opts, kratos.Metadata(AppInfo.Metadata))
	}
	if AppInfo.Name != "" {
		opts = append(opts, kratos.Name(AppInfo.Name))
	}
	if AppInfo.Version != "" {
		opts = append(opts, kratos.Version(AppInfo.Version))
	}
	if AppInfo.InstanceId != "" {
		opts = append(opts, kratos.ID(AppInfo.InstanceId))
	}

	return kratos.New(opts...)
}

// Bootstrap 应用引导启动
func Bootstrap(initApp InitAppFunc, appInfo *conf.AppInfo) error {
	if appInfo != nil {
		if appInfo.Name != "" {
			AppInfo.Name = appInfo.Name
		}
		if appInfo.Version != "" {
			AppInfo.Version = appInfo.Version
		}
		if appInfo.InstanceId != "" {
			AppInfo.InstanceId = appInfo.InstanceId
		}
		if appInfo.Metadata != nil {
			AppInfo.Metadata = appInfo.Metadata
		}
	}

	// bootstrap
	bctx, err := initBootstrap(context.Background(), AppInfo)
	if err != nil {
		return err
	}

	// init app
	app, cleanup, err := initApp(bctx)
	if err != nil {
		return err
	}
	defer cleanup()

	// run the app.
	if err = app.Run(); err != nil {
		return err
	}

	return nil
}

// BootstrapWithNameVersion 使用服务名称和版本引导启动应用
func BootstrapWithNameVersion(initApp InitAppFunc, serviceName, version *string) error {
	ai := &conf.AppInfo{}
	if serviceName != nil {
		ai.Name = *serviceName
	}
	if version != nil {
		ai.Version = *version
	}
	return Bootstrap(initApp, ai)
}

// initBootstrap 初始化引导程序
func initBootstrap(ctx context.Context, appInfo *conf.AppInfo) (*Context, error) {
	// inject command flags
	Flags := NewCommandFlags()
	Flags.Init()

	var err error
	var bctx Context

	// load configs
	if err = bConfig.LoadBootstrapConfig(Flags.Conf); err != nil {
		return &bctx, err
	}

	// get bootstrap config
	bctx.Config = bConfig.GetBootstrapConfig()
	if bctx.Config == nil {
		return &bctx, fmt.Errorf("bootstrap config is nil")
	}

	// init logger
	bctx.Logger = bLogger.NewLoggerProvider(bctx.Config.Logger, appInfo)
	if bctx.Logger == nil {
		return &bctx, fmt.Errorf("init logger failed")
	}

	// init registrar
	bctx.Registrar, err = bRegistry.NewRegistrar(bctx.Config.Registry)
	if err != nil {
		return &bctx, err
	}

	// init tracer
	if err = tracer.NewTracerProvider(ctx, bctx.Config.Trace, appInfo); err != nil {
		return &bctx, err
	}

	return &bctx, nil
}
