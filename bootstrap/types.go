package bootstrap

import (
	"github.com/go-kratos/kratos/v2"
)

// InitAppFunc 应用初始化函数类型
type InitAppFunc func(ctx *Context) (app *kratos.App, cleanup func(), err error)
