package lua

import (
	"fmt"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func testWorker() {
	L := luaPool.Borrow()
	defer luaPool.Return(L)
	if err := L.DoString(`print("hello")`); err != nil {
		panic(err)
	}
}

func TestStatePool(t *testing.T) {
	defer luaPool.Shutdown()
	go testWorker()
	go testWorker()
}

func TestLuaTableMap(t *testing.T) {
	// 将Lua表映射到Go结构
	type Role struct {
		Name string
	}

	type Person struct {
		Name      string
		Age       int
		WorkPlace string
		Role      []*Role
	}

	L := luaPool.Borrow()
	if err := L.DoString(`
person = {
  name = "Michel",
  age  = "31", -- weakly input
  work_place = "San Jose",
  role = {
    {
      name = "Administrator"
    },
    {
      name = "Operator"
    }
  }
}
`); err != nil {
		panic(err)
	}

	var person Person
	if err := gluamapper.Map(L.GetGlobal("person").(*lua.LTable), &person); err != nil {
		panic(err)
	}
	fmt.Printf("%v+", person)
}
