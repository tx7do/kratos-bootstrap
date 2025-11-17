package js

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

	jsString := `
console.log('ddddd')
role.Name = '角色'
menu.Name = '目录'
`

	err := exe.LoadString(jsString)
	assert.Nil(t, err)

	type Role struct {
		Name string
	}
	type Menu struct {
		Name string
	}

	var role Role
	var menu Menu

	err = exe.Register("role", &role)
	assert.Nil(t, err)

	err = exe.Register("menu", &menu)
	assert.Nil(t, err)

	err = exe.ExecuteString(jsString)
	assert.Nil(t, err)

	fmt.Println(role)
	fmt.Println(menu)
}

func sayHello(msg string) {
	fmt.Printf("golang say Hello, %s.\n", msg)
}

func TestVirtualMachine_ExecuteFile(t *testing.T) {
	exe := NewVirtualMachine()
	defer exe.Destroy()

	err := exe.LoadFile("./script/test.js")
	assert.Nil(t, err)

	u := &user{
		Name: "Tim",
	}

	err = exe.Register("u", u)
	assert.Nil(t, err)

	err = exe.Register("sayHello", sayHello)
	assert.Nil(t, err)

	err = exe.Execute()
	assert.Nil(t, err)

	fmt.Println("Js set your token to:", u.Token())

	var fn func(string)
	err = exe.GetFunction("printMessage", &fn)
	fn("hello")
	fn("world")
	assert.Nil(t, err)
}

func TestVirtualMachine_DirectExecuteFile(t *testing.T) {
	exe := NewVirtualMachine()
	defer exe.Destroy()

	err := exe.ExecuteFile("./script/test_require.js")
	assert.Nil(t, err)
	//fmt.Println(err.Error())
}

func TestVirtualMachine_Require(t *testing.T) {
	exe := NewVirtualMachine()
	defer exe.Destroy()

	const SCRIPT = `
	var m = require("./script/m.js");
	m.test();
	m.sayHi('CI');
	`

	err := exe.ExecuteString(SCRIPT)
	assert.Nil(t, err)
	err = exe.ExecuteFile("./script/test_require.js")
	assert.Nil(t, err)
	//fmt.Println(err.Error())
}
