package bootstrap

import (
	"sync"

	kratosLog "github.com/go-kratos/kratos/v2/log"
	kratosRegistry "github.com/go-kratos/kratos/v2/registry"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// Context 引导上下文
type Context struct {
	Config    *conf.Bootstrap          // 引导配置
	Logger    kratosLog.Logger         // 日志记录器
	Registrar kratosRegistry.Registrar // 服务注册器

	CustomConfig sync.Map // 自定义配置项
	Values       sync.Map // 自定义值存储
}

func (c *Context) NewLoggerHelper(moduleName string) *kratosLog.Helper {
	return kratosLog.NewHelper(kratosLog.With(c.Logger, "module", moduleName))
}

// GetConfig 返回当前的 *conf.Bootstrap（并发安全）
func (c *Context) GetConfig() *conf.Bootstrap {
	return c.Config
}

// SetCustomConfig 存入自定义配置
func (c *Context) SetCustomConfig(key string, val any) {
	c.CustomConfig.Store(key, val)
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
