package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	kratosLog "github.com/go-kratos/kratos/v2/log"
	kratosRegistry "github.com/go-kratos/kratos/v2/registry"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	bConfig "github.com/tx7do/kratos-bootstrap/config"
)

// Context 引导上下文
type Context struct {
	config  *conf.Bootstrap // 引导配置
	appInfo *conf.AppInfo   // 应用信息

	logger    kratosLog.Logger         // 日志记录器
	registrar kratosRegistry.Registrar // 服务注册器

	customConfig sync.Map // 自定义配置项
	values       sync.Map // 自定义值存储

	rootCtx context.Context    // 应用级根上下文（可用于优雅关闭）
	cancel  context.CancelFunc // 取消函数
}

// NewContext 创建带 cancel 的应用级 Context（传 nil 使用 Background）
func NewContext(parent context.Context, ai *conf.AppInfo) *Context {
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)

	c := &Context{
		appInfo: &conf.AppInfo{},
	}
	// 初始化默认信息
	AdjustAppInfo(c.appInfo)

	c.copyAppInfo(ai)

	// 其余初始化例如 RootCtx/Cancel/Logger 可在这里设置
	_ = cancel // 保留 cancel 给调用者或另行设置
	_ = ctx
	return c
}

func NewContextWithParam(parent context.Context, ai *conf.AppInfo, cfg *conf.Bootstrap, log kratosLog.Logger) *Context {
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)

	c := &Context{
		appInfo: &conf.AppInfo{},
		config:  cfg,
		logger:  log,
	}
	// 初始化默认信息
	AdjustAppInfo(c.appInfo)

	c.copyAppInfo(ai)

	// 其余初始化例如 RootCtx/Cancel/Logger 可在这里设置
	_ = cancel // 保留 cancel 给调用者或另行设置
	_ = ctx
	return c
}

// Context 返回应用级根 context（保证非 nil）
func (c *Context) Context() context.Context {
	if c == nil || c.rootCtx == nil {
		return context.Background()
	}
	return c.rootCtx
}

// CancelContext 触发取消（幂等）
func (c *Context) CancelContext() {
	if c == nil {
		return
	}
	if c.cancel != nil {
		c.cancel()
	}
}

func (c *Context) NewLoggerHelper(moduleName string) *kratosLog.Helper {
	return kratosLog.NewHelper(kratosLog.With(c.logger, "module", moduleName))
}

func (c *Context) GetLogger() kratosLog.Logger {
	return c.logger
}

// GetConfig 返回当前的 *conf.bootstrap（并发安全）
func (c *Context) GetConfig() *conf.Bootstrap {
	if c.config == nil {
		return nil
	}
	if clone := proto.Clone(c.config); clone != nil {
		if b, ok := clone.(*conf.Bootstrap); ok {
			return b
		}
	}
	return nil
}

func (c *Context) GetAppInfo() *conf.AppInfo {
	if c.appInfo == nil {
		return nil
	}
	if clone := proto.Clone(c.appInfo); clone != nil {
		if a, ok := clone.(*conf.AppInfo); ok {
			return a
		}
	}
	return nil
}

// setAppInfo 用受控方式替换整个 appInfo（可选）
func (c *Context) setAppInfo(src *conf.AppInfo) {
	if c == nil || src == nil {
		return
	}
	AdjustAppInfo(src)

	c.appInfo = &conf.AppInfo{
		Name:       src.Name,
		Version:    src.Version,
		AppId:      src.AppId,
		Project:    src.Project,
		InstanceId: src.InstanceId,
		Hostname:   src.Hostname,
		StartTime:  src.StartTime,
		Metadata:   cloneMetadata(src.Metadata),
	}
}

// copyAppInfo 复制应用信息
func (c *Context) copyAppInfo(ai *conf.AppInfo) {
	if ai == nil {
		return
	}

	// 先修正输入，避免未初始化字段
	AdjustAppInfo(ai)

	if ai.Name != "" {
		c.appInfo.Name = ai.Name
	}
	if ai.Project != "" {
		c.appInfo.Project = ai.Project
	}
	if ai.AppId != "" {
		c.appInfo.AppId = ai.AppId
	}
	if ai.Version != "" {
		c.appInfo.Version = ai.Version
	}
	if ai.InstanceId != "" {
		c.appInfo.InstanceId = ai.InstanceId
	}
	if ai.Metadata != nil {
		c.appInfo.Metadata = ai.Metadata
	}
}

func (c *Context) PrintAppInfo() {
	ai := c.GetAppInfo()
	if ai == nil {
		return
	}
	ts := time.Now().Format(time.RFC3339)
	host, _ := os.Hostname()
	pid := os.Getpid()

	if os.Getenv("APPINFO_FORMAT") == "json" {
		out := map[string]interface{}{
			"timestamp":   ts,
			"host":        host,
			"pid":         pid,
			"name":        ai.Name,
			"version":     ai.Version,
			"app_id":      ai.AppId,
			"instance_id": ai.InstanceId,
			"metadata":    ai.Metadata,
		}
		if b, err := json.Marshal(out); err == nil {
			fmt.Println(string(b))
		} else {
			fmt.Printf("Application info marshal error: %v\n", err)
		}
		return
	}

	fmt.Printf("[%s] %s (pid:%d@%s)\n", ts, ai.Name, pid, host)
	fmt.Printf("  Version: %s\n", ai.Version)
	fmt.Printf("  AppId: %s\n", ai.AppId)
	fmt.Printf("  InstanceId: %s\n", ai.InstanceId)
	if len(ai.Metadata) > 0 {
		fmt.Println("  Metadata:")
		keys := make([]string, 0, len(ai.Metadata))
		for k := range ai.Metadata {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Printf("    %s=%s\n", k, ai.Metadata[k])
		}
	}
}

func (c *Context) GetRegistrar() kratosRegistry.Registrar {
	return c.registrar
}

// RegisterCustomConfig 注册自定义配置
func (c *Context) RegisterCustomConfig(key string, cfg proto.Message) {
	if key == "" || cfg == nil {
		return
	}

	if _, ok := c.customConfig.Load(key); ok {
		return
	}

	c.customConfig.Store(key, cfg)

	bConfig.RegisterConfig(cfg)
}

// SetCustomConfig 存入自定义配置
func (c *Context) SetCustomConfig(key string, cfg proto.Message) {
	if key == "" || cfg == nil {
		return
	}

	c.customConfig.Store(key, cfg)
}

// GetCustomConfig 获取自定义配置（原始类型）
func (c *Context) GetCustomConfig(key string) (any, bool) {
	return c.customConfig.Load(key)
}

// DeleteCustomConfig 删除自定义配置
func (c *Context) DeleteCustomConfig(key string) {
	c.customConfig.Delete(key)
}

// RangeCustomConfig 遍历自定义配置，回调返回 false 可停止遍历
func (c *Context) RangeCustomConfig(fn func(key string, val any) bool) {
	c.customConfig.Range(func(k, v any) bool {
		ks, _ := k.(string)
		return fn(ks, v)
	})
}

// SetValue 将任意值放入通用存储
func (c *Context) SetValue(key string, val interface{}) {
	c.values.Store(key, val)
}

// GetValue 从通用存储读取值
func (c *Context) GetValue(key string) (interface{}, bool) {
	return c.values.Load(key)
}
