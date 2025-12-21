package bootstrap

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2"
	kratosLog "github.com/go-kratos/kratos/v2/log"
	kratosRegistry "github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/spf13/cobra"
	bConfig "github.com/tx7do/kratos-bootstrap/config"
	bLogger "github.com/tx7do/kratos-bootstrap/logger"
	bRegistry "github.com/tx7do/kratos-bootstrap/registry"
	"github.com/tx7do/kratos-bootstrap/tracer"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewApp 创建应用程序
func NewApp(ll kratosLog.Logger, rr kratosRegistry.Registrar, srv ...transport.Server) *kratos.App {
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

	if appInfo.Metadata != nil {
		opts = append(opts, kratos.Metadata(appInfo.Metadata))
	}
	if appInfo.AppId != "" {
		opts = append(opts, kratos.Name(appInfo.AppId))
	}
	if appInfo.Version != "" {
		opts = append(opts, kratos.Version(appInfo.Version))
	}
	if appInfo.InstanceId != "" {
		opts = append(opts, kratos.ID(appInfo.InstanceId))
	}

	return kratos.New(opts...)
}

// RunApp 运行应用程序并允许在执行前对 root 命令做定制。
// opts 可用于注册子命令、对 root 添加 flag 或其他修改。
func RunApp(initApp InitAppFunc, ai *conf.AppInfo, opts ...func(root *cobra.Command)) error {
	// 注入命令行参数
	root := NewRootCmd(flags, func(cmd *cobra.Command, args []string) error {
		return Bootstrap(initApp, ai)
	})

	// 允许调用方定制 root（如添加子命令、注册额外 flag 等）
	for _, opt := range opts {
		if opt != nil {
			opt(root)
		}
	}

	// 如果 flags 实现了 Register，就在 Execute 前注册到命令上，确保 cobra 能解析这些 flag
	if rb, ok := interface{}(flags).(interface{ Register(cmd *cobra.Command) }); ok {
		rb.Register(root)
	}

	if err := root.Execute(); err != nil {
		return err
	}
	return nil
}

// Bootstrap 应用引导启动
func Bootstrap(initApp InitAppFunc, ai *conf.AppInfo) error {
	// 设置应用信息
	copyAppInfo(ai)

	// 根据注册中心类型规范化 AppId
	if bConfig.GetBootstrapConfig().Registry != nil && bConfig.GetBootstrapConfig().Registry.GetType() != "" {
		appInfo.AppId = bRegistry.NormalizeForRegistry(appInfo.AppId, bConfig.GetBootstrapConfig().Registry.GetType())
	}

	// 打印应用信息
	printAppInfo()

	// bootstrap
	bctx, err := initBootstrap(context.Background(), appInfo)
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

// initBootstrap 初始化引导程序
func initBootstrap(ctx context.Context, ai *conf.AppInfo) (*Context, error) {
	var err error
	var bctx Context

	// load configs
	if err = bConfig.LoadBootstrapConfig(flags.Conf); err != nil {
		return &bctx, err
	}

	// get bootstrap config
	bctx.Config = bConfig.GetBootstrapConfig()
	if bctx.Config == nil {
		return &bctx, fmt.Errorf("bootstrap config is nil")
	}

	// init logger
	bctx.Logger = bLogger.NewLoggerProvider(bctx.Config.Logger, ai)
	if bctx.Logger == nil {
		return &bctx, fmt.Errorf("init logger failed")
	}

	// init registrar
	bctx.Registrar, err = bRegistry.NewRegistrar(bctx.Config.Registry)
	if err != nil {
		return &bctx, err
	}

	// init tracer
	if err = tracer.NewTracerProvider(ctx, bctx.Config.Trace, ai); err != nil {
		return &bctx, err
	}

	return &bctx, nil
}
