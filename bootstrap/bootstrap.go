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

// InitApp 应用初始化函数类型
type InitApp func(logger kratosLog.Logger, registrar kratosRegistry.Registrar, bootstrap *conf.Bootstrap) (app *kratos.App, cleanup func(), err error)

// Context 引导上下文
type Context struct {
	Config    *conf.Bootstrap          // 引导配置
	Logger    kratosLog.Logger         // 日志记录器
	Registrar kratosRegistry.Registrar // 服务注册器
}

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
func Bootstrap(initApp InitApp, serviceName, version *string) error {
	if serviceName != nil && len(*serviceName) != 0 {
		AppInfo.Name = *serviceName
	}
	if version != nil && len(*version) != 0 {
		AppInfo.Version = *version
	}

	// bootstrap
	bctx, err := initBootstrap(context.Background(), AppInfo)
	if err != nil {
		return err
	}

	// init app
	app, cleanup, err := initApp(bctx.Logger, bctx.Registrar, bctx.Config)
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
