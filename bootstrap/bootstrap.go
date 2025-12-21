package bootstrap

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport"

	bConfig "github.com/tx7do/kratos-bootstrap/config"
	bLogger "github.com/tx7do/kratos-bootstrap/logger"
	bRegistry "github.com/tx7do/kratos-bootstrap/registry"
	"github.com/tx7do/kratos-bootstrap/tracer"
)

// NewApp 创建应用程序
func NewApp(ctx *Context, srv ...transport.Server) *kratos.App {
	var opts []kratos.Option
	if ctx.logger != nil {
		opts = append(opts, kratos.Logger(ctx.logger))
	}
	if ctx.registrar != nil {
		opts = append(opts, kratos.Registrar(ctx.registrar))
	}
	if len(srv) > 0 {
		opts = append(opts, kratos.Server(srv...))
	}

	if ctx.appInfo.Metadata != nil {
		opts = append(opts, kratos.Metadata(ctx.appInfo.Metadata))
	}
	if ctx.appInfo.AppId != "" {
		registerName := ctx.appInfo.Project + "/" + ctx.appInfo.AppId
		// 根据注册中心类型规范化 AppId
		if bConfig.GetBootstrapConfig().Registry != nil && bConfig.GetBootstrapConfig().Registry.GetType() != "" {
			registerName = bRegistry.NormalizeForRegistry(registerName, bConfig.GetBootstrapConfig().Registry.GetType())
		}
		opts = append(opts, kratos.Name(registerName))
	}
	if ctx.appInfo.Version != "" {
		opts = append(opts, kratos.Version(ctx.appInfo.Version))
	}
	if ctx.appInfo.InstanceId != "" {
		opts = append(opts, kratos.ID(ctx.appInfo.InstanceId))
	}

	return kratos.New(opts...)
}

// RunApp 运行应用程序并允许在执行前对 root 命令做定制。
// opts 可用于注册子命令、对 root 添加 flag 或其他修改。
func RunApp(ctx *Context, initApp InitAppFunc, opts ...func(root *cobra.Command)) error {
	if ctx == nil {
		return fmt.Errorf("bootstrap context is nil")
	}

	// 注入命令行参数
	root := NewRootCmd(flags, func(cmd *cobra.Command, args []string) error {
		return bootstrap(ctx, initApp)
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

// bootstrap 应用引导启动
func bootstrap(ctx *Context, initApp InitAppFunc) error {
	// 打印应用信息
	ctx.PrintAppInfo()

	var err error

	// load configs
	if err = bConfig.LoadBootstrapConfig(flags.Conf); err != nil {
		return err
	}

	// get bootstrap config
	ctx.config = bConfig.GetBootstrapConfig()
	if ctx.config == nil {
		return fmt.Errorf("bootstrap config is nil")
	}

	// init logger
	ctx.logger = bLogger.NewLoggerProvider(ctx.config.Logger, ctx.appInfo)
	if ctx.logger == nil {
		return fmt.Errorf("init logger failed")
	}

	// init registrar
	ctx.registrar, err = bRegistry.NewRegistrar(ctx.config.Registry)
	if err != nil {
		return err
	}

	// init tracer
	if err = tracer.NewTracerProvider(ctx.Context(), ctx.config.Trace, ctx.appInfo); err != nil {
		return err
	}

	// init app
	app, cleanup, err := initApp(ctx)
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
