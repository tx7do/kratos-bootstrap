package lua

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	Lua "github.com/yuin/gopher-lua"
)

// LuaEngine Lua 脚本引擎实现
type LuaEngine struct {
	vm          *virtualMachine
	initialized bool
	lastError   error
}

// NewLuaEngine 创建 Lua 引擎实例
func NewLuaEngine() *LuaEngine {
	return &LuaEngine{
		initialized: false,
	}
}

// Init 初始化引擎
func (e *LuaEngine) Init(ctx context.Context) error {
	if e.initialized {
		return fmt.Errorf("engine already initialized")
	}

	e.vm = NewVirtualMachine()
	e.initialized = true
	e.lastError = nil

	return nil
}

// Destroy 销毁引擎
func (e *LuaEngine) Destroy() error {
	if !e.initialized {
		return fmt.Errorf("engine not initialized")
	}

	e.vm.Destroy()
	e.vm = nil
	e.initialized = false

	return nil
}

// IsInitialized 检查是否已初始化
func (e *LuaEngine) IsInitialized() bool {
	return e.initialized
}

// LoadString 加载字符串脚本
func (e *LuaEngine) LoadString(ctx context.Context, source string) error {
	if !e.initialized {
		return fmt.Errorf("engine not initialized")
	}

	if err := e.vm.LoadString(source); err != nil {
		e.lastError = err
		return err
	}

	return nil
}

// LoadFile 加载脚本文件
func (e *LuaEngine) LoadFile(ctx context.Context, filePath string) error {
	if !e.initialized {
		return fmt.Errorf("engine not initialized")
	}

	if err := e.vm.LoadFile(filePath); err != nil {
		e.lastError = err
		return err
	}

	return nil
}

// LoadReader 从 Reader 加载脚本
func (e *LuaEngine) LoadReader(ctx context.Context, reader io.Reader, name string) error {
	if !e.initialized {
		return fmt.Errorf("engine not initialized")
	}

	source, err := ioutil.ReadAll(reader)
	if err != nil {
		e.lastError = err
		return err
	}

	return e.LoadString(ctx, string(source))
}

// Execute 执行已加载的脚本
func (e *LuaEngine) Execute(ctx context.Context) (interface{}, error) {
	if !e.initialized {
		return nil, fmt.Errorf("engine not initialized")
	}

	// 使用 channel 处理超时
	done := make(chan error, 1)

	go func() {
		done <- e.vm.Execute()
	}()

	select {
	case <-ctx.Done():
		e.lastError = ctx.Err()
		return nil, ctx.Err()
	case err := <-done:
		if err != nil {
			e.lastError = err
			return nil, err
		}
		return nil, nil
	}
}

// ExecuteString 执行字符串脚本
func (e *LuaEngine) ExecuteString(ctx context.Context, source string) (interface{}, error) {
	if !e.initialized {
		return nil, fmt.Errorf("engine not initialized")
	}

	done := make(chan error, 1)

	go func() {
		done <- e.vm.ExecuteString(source)
	}()

	select {
	case <-ctx.Done():
		e.lastError = ctx.Err()
		return nil, ctx.Err()
	case err := <-done:
		if err != nil {
			e.lastError = err
			return nil, err
		}
		return nil, nil
	}
}

// ExecuteFile 执行脚本文件
func (e *LuaEngine) ExecuteFile(ctx context.Context, filePath string) (interface{}, error) {
	if !e.initialized {
		return nil, fmt.Errorf("engine not initialized")
	}

	done := make(chan error, 1)

	go func() {
		done <- e.vm.ExecuteFile(filePath)
	}()

	select {
	case <-ctx.Done():
		e.lastError = ctx.Err()
		return nil, ctx.Err()
	case err := <-done:
		if err != nil {
			e.lastError = err
			return nil, err
		}
		return nil, nil
	}
}

// RegisterGlobal 注册全局变量
func (e *LuaEngine) RegisterGlobal(name string, value interface{}) error {
	if !e.initialized {
		return fmt.Errorf("engine not initialized")
	}

	e.vm.BindStruct(name, value)
	return nil
}

// GetGlobal 获取全局变量
func (e *LuaEngine) GetGlobal(name string) (interface{}, error) {
	if !e.initialized {
		return nil, fmt.Errorf("engine not initialized")
	}

	lv := e.vm.L.GetGlobal(name)
	return e.vm.convertFromLValue(lv), nil
}

// RegisterFunction 注册全局函数
func (e *LuaEngine) RegisterFunction(name string, fn interface{}) error {
	if !e.initialized {
		return fmt.Errorf("engine not initialized")
	}

	// 类型断言检查是否为 Lua.LGFunction
	if lf, ok := fn.(Lua.LGFunction); ok {
		e.vm.RegisterFunction(name, lf)
		return nil
	}

	return fmt.Errorf("function must be of type Lua.LGFunction")
}

// CallFunction 调用 Lua 函数
func (e *LuaEngine) CallFunction(ctx context.Context, name string, args ...interface{}) (interface{}, error) {
	if !e.initialized {
		return nil, fmt.Errorf("engine not initialized")
	}

	type result struct {
		value interface{}
		err   error
	}

	done := make(chan result, 1)

	go func() {
		// 转换参数
		var lArgs []Lua.LValue
		for _, arg := range args {
			lArgs = append(lArgs, e.vm.convertToLValue(arg))
		}

		// 调用函数
		err := e.vm.L.CallByParam(Lua.P{
			Fn:      e.vm.L.GetGlobal(name),
			NRet:    1,
			Protect: true,
		}, lArgs...)

		if err != nil {
			done <- result{nil, err}
			return
		}

		// 获取返回值
		ret := e.vm.L.Get(-1)
		e.vm.L.Pop(1)

		done <- result{e.vm.convertFromLValue(ret), nil}
	}()

	select {
	case <-ctx.Done():
		e.lastError = ctx.Err()
		return nil, ctx.Err()
	case res := <-done:
		if res.err != nil {
			e.lastError = res.err
		}
		return res.value, res.err
	}
}

// RegisterModule 注册模块
func (e *LuaEngine) RegisterModule(name string, module interface{}) error {
	if !e.initialized {
		return fmt.Errorf("engine not initialized")
	}

	if mod, ok := module.(Lua.LGFunction); ok {
		e.vm.RegisterModule(name, mod)
		return nil
	}

	return fmt.Errorf("module must be of type Lua.LGFunction")
}

// GetLastError 获取最后一个错误
func (e *LuaEngine) GetLastError() error {
	return e.lastError
}

// ClearError 清除错误
func (e *LuaEngine) ClearError() {
	e.lastError = nil
}

// GetState 获取 Lua 状态机（扩展方法）
func (e *LuaEngine) GetState() *Lua.LState {
	if e.vm == nil {
		return nil
	}
	return e.vm.L
}
