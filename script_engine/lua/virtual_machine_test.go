package lua

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type user struct {
	Name  string
	token string
}

func (u *user) SetToken(t string) {
	u.token = t
}

func (u *user) Token() string {
	return u.token
}

func TestVirtualMachine_ExecuteString(t *testing.T) {
	exe := NewVirtualMachine()
	defer exe.Destroy()

	luaString := `
role = 
    {
      name = "Administrator"
    }
menu = 
    {
      name = "TestMenu"
    }
print("hello")`

	err := exe.LoadString(luaString)
	assert.Nil(t, err)

	type Role struct {
		Name string
	}
	type Menu struct {
		Name string
	}

	var role Role
	var menu Menu

	err = exe.Execute()
	assert.Nil(t, err)

	_ = exe.GetLuaTableToStruct("role", &role)
	_ = exe.GetLuaTableToStruct("menu", &menu)

	fmt.Println(role)
	fmt.Println(menu)
}

func TestVirtualMachine_ExecuteFile(t *testing.T) {
	exe := NewVirtualMachine()
	defer exe.Destroy()

	err := exe.LoadFile("./script/test.lua")
	assert.Nil(t, err)

	u := &user{
		Name: "Tim",
	}

	exe.BindStruct("u", u)

	err = exe.Execute()
	assert.Nil(t, err)

	fmt.Println("Lua set your token to:", u.Token())
}

func TestVirtualMachine_HttpModule(t *testing.T) {
	exe := NewVirtualMachine()
	defer exe.Destroy()

	err := exe.LoadFile("./script/test_http.lua")
	assert.Nil(t, err)

	err = exe.Execute()
	assert.Nil(t, err)
}

func TestVirtualMachine_LoadModule(t *testing.T) {
	exe := NewVirtualMachine()
	defer exe.Destroy()

	err := exe.LoadFile("./script/test_load_module.lua")
	assert.Nil(t, err)

	err = exe.Execute()
	assert.Nil(t, err)
}

func TestVirtualMachine_CryptoModule(t *testing.T) {
	exe := NewVirtualMachine()
	defer exe.Destroy()

	err := exe.LoadFile("./script/test_crypto.lua")
	assert.Nil(t, err)

	err = exe.Execute()
	assert.Nil(t, err)
}

func TestVirtualMachine_Debugger(t *testing.T) {
	exe := NewVirtualMachine()
	defer exe.Destroy()

	err := exe.LoadFile("./script/test_debugger.lua")
	assert.Nil(t, err)

	err = exe.Execute()
	assert.Nil(t, err)
}
