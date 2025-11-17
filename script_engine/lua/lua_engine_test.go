package lua

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestLuaEngine(t *testing.T) {
	// 创建引擎
	engine := NewLuaEngine()
	defer engine.Destroy()

	// 初始化
	ctx := context.Background()
	if err := engine.Init(ctx); err != nil {
		t.Fatal(err)
	}

	// 注册全局变量
	engine.RegisterGlobal("config", map[string]interface{}{
		"host": "localhost",
		"port": 8080,
	})

	// 执行脚本
	result, err := engine.ExecuteString(ctx, `
    function add(a, b)
        return a + b
    end
`)

	// 调用函数（带超时）
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err = engine.CallFunction(ctx, "add", 10, 20)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(result) // 输出: 30
}
