package bootstrap

import (
	"github.com/go-kratos/kratos/v2"
	kratosLog "github.com/go-kratos/kratos/v2/log"
	kratosRegistry "github.com/go-kratos/kratos/v2/registry"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// InitAppFunc 应用初始化函数类型
type InitAppFunc func(ctx *Context) (app *kratos.App, cleanup func(), err error)

// Context 引导上下文
type Context struct {
	Config    *conf.Bootstrap          // 引导配置
	Logger    kratosLog.Logger         // 日志记录器
	Registrar kratosRegistry.Registrar // 服务注册器
}
