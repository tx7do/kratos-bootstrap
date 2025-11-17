package js

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"sync"

	"github.com/dop251/goja"
)

// JavascriptEngine JavaScript 脚本引擎实现
type JavascriptEngine struct {
	runtime     *goja.Runtime
	program     *goja.Program
	initialized bool
	lastError   error
	mu          sync.RWMutex
}

// NewJavascriptEngine 创建 JavaScript 引擎实例
func NewJavascriptEngine() *JavascriptEngine {
	return &JavascriptEngine{
		initialized: false,
	}
}

// Init 初始化引擎
func (e *JavascriptEngine) Init(ctx context.Context) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.initialized {
		return fmt.Errorf("engine already initialized")
	}

	e.runtime = goja.New()
	e.initialized = true
	e.lastError = nil

	return nil
}

// Destroy 销毁引擎
func (e *JavascriptEngine) Destroy() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.initialized {
		return fmt.Errorf("engine not initialized")
	}

	e.runtime = nil
	e.program = nil
	e.initialized = false

	return nil
}

// IsInitialized 检查是否已初始化
func (e *JavascriptEngine) IsInitialized() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.initialized
}

// LoadString 加载字符串脚本
func (e *JavascriptEngine) LoadString(ctx context.Context, source string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.initialized {
		return fmt.Errorf("engine not initialized")
	}

	program, err := goja.Compile("", source, true)
	if err != nil {
		e.lastError = err
		return err
	}

	e.program = program
	return nil
}

// LoadFile 加载脚本文件
func (e *JavascriptEngine) LoadFile(ctx context.Context, filePath string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.initialized {
		return fmt.Errorf("engine not initialized")
	}

	source, err := ioutil.ReadFile(filePath)
	if err != nil {
		e.lastError = err
		return err
	}

	program, err := goja.Compile(filePath, string(source), true)
	if err != nil {
		e.lastError = err
		return err
	}

	e.program = program
	return nil
}

// LoadReader 从 Reader 加载脚本
func (e *JavascriptEngine) LoadReader(ctx context.Context, reader io.Reader, name string) error {
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
func (e *JavascriptEngine) Execute(ctx context.Context) (interface{}, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.initialized {
		return nil, fmt.Errorf("engine not initialized")
	}

	if e.program == nil {
		return nil, fmt.Errorf("no program loaded")
	}

	type result struct {
		value goja.Value
		err   error
	}

	done := make(chan result, 1)

	go func() {
		val, err := e.runtime.RunProgram(e.program)
		done <- result{val, err}
	}()

	select {
	case <-ctx.Done():
		e.runtime.Interrupt(ctx.Err())
		e.lastError = ctx.Err()
		return nil, ctx.Err()
	case res := <-done:
		if res.err != nil {
			e.lastError = res.err
			return nil, res.err
		}
		return res.value.Export(), nil
	}
}

// ExecuteString 执行字符串脚本
func (e *JavascriptEngine) ExecuteString(ctx context.Context, source string) (interface{}, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.initialized {
		return nil, fmt.Errorf("engine not initialized")
	}

	type result struct {
		value goja.Value
		err   error
	}

	done := make(chan result, 1)

	go func() {
		val, err := e.runtime.RunString(source)
		done <- result{val, err}
	}()

	select {
	case <-ctx.Done():
		e.runtime.Interrupt(ctx.Err())
		e.lastError = ctx.Err()
		return nil, ctx.Err()
	case res := <-done:
		if res.err != nil {
			e.lastError = res.err
			return nil, res.err
		}
		return res.value.Export(), nil
	}
}

// ExecuteFile 执行脚本文件
func (e *JavascriptEngine) ExecuteFile(ctx context.Context, filePath string) (interface{}, error) {
	if err := e.LoadFile(ctx, filePath); err != nil {
		return nil, err
	}
	return e.Execute(ctx)
}

// RegisterGlobal 注册全局变量
func (e *JavascriptEngine) RegisterGlobal(name string, value interface{}) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.initialized {
		return fmt.Errorf("engine not initialized")
	}

	e.runtime.Set(name, value)
	return nil
}

// GetGlobal 获取全局变量
func (e *JavascriptEngine) GetGlobal(name string) (interface{}, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if !e.initialized {
		return nil, fmt.Errorf("engine not initialized")
	}

	val := e.runtime.Get(name)
	if val == nil {
		return nil, fmt.Errorf("global variable %s not found", name)
	}

	return val.Export(), nil
}

// RegisterFunction 注册全局函数
func (e *JavascriptEngine) RegisterFunction(name string, fn interface{}) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.initialized {
		return fmt.Errorf("engine not initialized")
	}

	e.runtime.Set(name, fn)
	return nil
}

// CallFunction 调用 JavaScript 函数
func (e *JavascriptEngine) CallFunction(ctx context.Context, name string, args ...interface{}) (interface{}, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.initialized {
		return nil, fmt.Errorf("engine not initialized")
	}

	type result struct {
		value goja.Value
		err   error
	}

	done := make(chan result, 1)

	go func() {
		fn, ok := goja.AssertFunction(e.runtime.Get(name))
		if !ok {
			done <- result{nil, fmt.Errorf("function %s not found", name)}
			return
		}

		// 转换参数
		var gojaArgs []goja.Value
		for _, arg := range args {
			gojaArgs = append(gojaArgs, e.runtime.ToValue(arg))
		}

		val, err := fn(goja.Undefined(), gojaArgs...)
		done <- result{val, err}
	}()

	select {
	case <-ctx.Done():
		e.runtime.Interrupt(ctx.Err())
		e.lastError = ctx.Err()
		return nil, ctx.Err()
	case res := <-done:
		if res.err != nil {
			e.lastError = res.err
			return nil, res.err
		}
		return res.value.Export(), nil
	}
}

// RegisterModule 注册模块
func (e *JavascriptEngine) RegisterModule(name string, module interface{}) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.initialized {
		return fmt.Errorf("engine not initialized")
	}

	// 创建模块对象
	moduleObj := e.runtime.NewObject()

	// 如果 module 是 map，则设置属性
	if m, ok := module.(map[string]interface{}); ok {
		for k, v := range m {
			moduleObj.Set(k, v)
		}
	} else {
		// 否则直接设置整个对象
		e.runtime.Set(name, module)
		return nil
	}

	e.runtime.Set(name, moduleObj)
	return nil
}

// GetLastError 获取最后一个错误
func (e *JavascriptEngine) GetLastError() error {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.lastError
}

// ClearError 清除错误
func (e *JavascriptEngine) ClearError() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.lastError = nil
}

// GetRuntime 获取 Goja 运行时（扩展方法）
func (e *JavascriptEngine) GetRuntime() *goja.Runtime {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.runtime
}

// RunProgram 运行已编译的程序（扩展方法）
func (e *JavascriptEngine) RunProgram(ctx context.Context, program *goja.Program) (goja.Value, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.initialized {
		return nil, fmt.Errorf("engine not initialized")
	}

	type result struct {
		value goja.Value
		err   error
	}

	done := make(chan result, 1)

	go func() {
		val, err := e.runtime.RunProgram(program)
		done <- result{val, err}
	}()

	select {
	case <-ctx.Done():
		e.runtime.Interrupt(ctx.Err())
		return nil, ctx.Err()
	case res := <-done:
		return res.value, res.err
	}
}
