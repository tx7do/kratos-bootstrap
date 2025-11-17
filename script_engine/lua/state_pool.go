package lua

import (
	"sync"

	Lua "github.com/yuin/gopher-lua"
)

func init() {
	luaPool = newStatePool()
}

var luaPool = newStatePool()

type lLuaStateArray []*Lua.LState

type lStatePool struct {
	m     sync.Mutex
	saved lLuaStateArray
}

func newStatePool() *lStatePool {
	return &lStatePool{
		saved: make(lLuaStateArray, 0, 10),
	}
}

func (pl *lStatePool) createLuaState() *Lua.LState {
	vm := Lua.NewState(Lua.Options{
		CallStackSize:       4096,
		RegistrySize:        4096,
		SkipOpenLibs:        true,
		IncludeGoStackTrace: true,
	})
	return vm
}

func (pl *lStatePool) Borrow() *Lua.LState {
	pl.m.Lock()
	defer pl.m.Unlock()
	n := len(pl.saved)
	if n == 0 {
		return pl.createLuaState()
	}
	x := pl.saved[n-1]
	pl.saved = pl.saved[0 : n-1]
	return x
}

func (pl *lStatePool) Return(L *Lua.LState) {
	pl.m.Lock()
	defer pl.m.Unlock()
	pl.saved = append(pl.saved, L)
}

func (pl *lStatePool) Shutdown() {
	for _, L := range pl.saved {
		L.Close()
	}
}
