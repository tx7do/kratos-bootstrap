package js

import (
	"errors"
	"os"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/process"
	"github.com/dop251/goja_nodejs/require"
)

type virtualMachine struct {
	vm       *goja.Runtime
	program  *goja.Program
	registry *require.Registry
}

func NewVirtualMachine() *virtualMachine {
	exec := &virtualMachine{
		vm:       goja.New(),
		registry: new(require.Registry),
	}
	exec.init()
	return exec
}

func (e *virtualMachine) init() {
	_ = e.registry.Enable(e.vm)
	console.Enable(e.vm)
	process.Enable(e.vm)
}

// Destroy 销毁虚拟机，为了性能考虑，现在只是将之还给虚拟机池。
func (e *virtualMachine) Destroy() {
}

// LoadString 加载字符串，并编译成字节码
func (e *virtualMachine) LoadString(source string) error {
	program, err := goja.Compile("", source, true)
	if err != nil {
		return err
	}

	e.program = program

	return nil
}

// LoadFile 加载文件，并编译成字节码
func (e *virtualMachine) LoadFile(filePath string) error {
	code, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	program, err := goja.Compile("", string(code), false)
	if err != nil {
		return err
	}

	e.program = program

	return nil
}

// Execute 执行已编译的lua代码
func (e *virtualMachine) Execute() error {
	if e.program == nil {
		return errors.New("no js")
	}
	_, err := e.vm.RunProgram(e.program)
	return err
}

// ExecuteString 直接执行字符串
func (e *virtualMachine) ExecuteString(source string) error {
	_, err := e.vm.RunString(source)
	return err
}

// ExecuteFile 直接执行lua文件
func (e *virtualMachine) ExecuteFile(filePath string) error {
	code, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	_, err = e.vm.RunScript("", string(code))
	return err
}

// Register 注册方法或者变量到js
func (e *virtualMachine) Register(name string, value interface{}) error {
	err := e.vm.Set(name, value)
	return err
}

// GetFunction 获取js中的方法
func (e *virtualMachine) GetFunction(name string, fn interface{}) error {
	err := e.vm.ExportTo(e.vm.Get(name), fn)
	return err
}
