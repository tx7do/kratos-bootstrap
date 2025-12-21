package bootstrap

import (
	"context"
	"sync"

	"google.golang.org/protobuf/proto"

	kratosLog "github.com/go-kratos/kratos/v2/log"
	kratosRegistry "github.com/go-kratos/kratos/v2/registry"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	bConfig "github.com/tx7do/kratos-bootstrap/config"
)

// Context 引导上下文
type Context struct {
	Config    *conf.Bootstrap          // 引导配置
	Logger    kratosLog.Logger         // 日志记录器
	Registrar kratosRegistry.Registrar // 服务注册器

	CustomConfig sync.Map // 自定义配置项
	Values       sync.Map // 自定义值存储

	RootCtx context.Context    // 应用级根上下文（可用于优雅关闭）
	Cancel  context.CancelFunc // 取消函数
}

// NewContext 创建带 cancel 的应用级 Context（传 nil 使用 Background）
func NewContext(parent context.Context) *Context {
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)
	return &Context{
		RootCtx: ctx,
		Cancel:  cancel,
	}
}

// Context 返回应用级根 context（保证非 nil）
func (c *Context) Context() context.Context {
	if c == nil || c.RootCtx == nil {
		return context.Background()
	}
	return c.RootCtx
}

// CancelContext 触发取消（幂等）
func (c *Context) CancelContext() {
	if c == nil {
		return
	}
	if c.Cancel != nil {
		c.Cancel()
	}
}

func (c *Context) NewLoggerHelper(moduleName string) *kratosLog.Helper {
	return kratosLog.NewHelper(kratosLog.With(c.Logger, "module", moduleName))
}

// GetConfig 返回当前的 *conf.Bootstrap（并发安全）
func (c *Context) GetConfig() *conf.Bootstrap {
	return c.Config
}

// RegisterCustomConfig 注册自定义配置
func (c *Context) RegisterCustomConfig(key string, cfg proto.Message) {
	if key == "" || cfg == nil {
		return
	}

	if _, ok := c.CustomConfig.Load(key); ok {
		return
	}

	c.CustomConfig.Store(key, cfg)

	bConfig.RegisterConfig(cfg)
}

// SetCustomConfig 存入自定义配置
func (c *Context) SetCustomConfig(key string, cfg proto.Message) {
	if key == "" || cfg == nil {
		return
	}

	c.CustomConfig.Store(key, cfg)
}

// GetCustomConfig 获取自定义配置（原始类型）
func (c *Context) GetCustomConfig(key string) (any, bool) {
	return c.CustomConfig.Load(key)
}

// DeleteCustomConfig 删除自定义配置
func (c *Context) DeleteCustomConfig(key string) {
	c.CustomConfig.Delete(key)
}

// RangeCustomConfig 遍历自定义配置，回调返回 false 可停止遍历
func (c *Context) RangeCustomConfig(fn func(key string, val any) bool) {
	c.CustomConfig.Range(func(k, v any) bool {
		ks, _ := k.(string)
		return fn(ks, v)
	})
}

// SetValue 将任意值放入通用存储
func (c *Context) SetValue(key string, val interface{}) {
	c.Values.Store(key, val)
}

// GetValue 从通用存储读取值
func (c *Context) GetValue(key string) (interface{}, bool) {
	return c.Values.Load(key)
}
