package script_engine

import (
	"context"
	"io"
	"time"
)

// EngineInterface Define the interface for script engines
type EngineInterface interface {
	// 生命周期管理

	Init(ctx context.Context) error
	Destroy() error
	IsInitialized() bool

	// 脚本加载

	LoadString(ctx context.Context, source string) error
	LoadFile(ctx context.Context, filePath string) error
	LoadReader(ctx context.Context, reader io.Reader, name string) error

	// 脚本执行

	Execute(ctx context.Context) (interface{}, error)
	ExecuteString(ctx context.Context, source string) (interface{}, error)
	ExecuteFile(ctx context.Context, filePath string) (interface{}, error)

	// 全局变量注册

	RegisterGlobal(name string, value interface{}) error
	GetGlobal(name string) (interface{}, error)

	// 函数调用

	RegisterFunction(name string, fn interface{}) error
	CallFunction(ctx context.Context, name string, args ...interface{}) (interface{}, error)

	// 模块管理

	RegisterModule(name string, module interface{}) error

	// 错误处理

	GetLastError() error
	ClearError()
}

// CallResult 函数调用结果
type CallResult struct {
	Values []interface{}
	Error  error
}

// ExecuteOptions 执行选项
type ExecuteOptions struct {
	Timeout  time.Duration
	Globals  map[string]interface{}
	MaxStack int
}
