package js

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestJavascriptEngine(t *testing.T) {
	// 创建引擎
	engine := NewJavascriptEngine()
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

	// 注册函数
	engine.RegisterFunction("log", func(msg string) {
		fmt.Println("JS Log:", msg)
	})

	// 执行脚本
	result, err := engine.ExecuteString(ctx, `
    function add(a, b) {
        log('Adding ' + a + ' and ' + b);
        return a + b;
    }
    add(10, 20);
`)
	fmt.Println(result) // 输出: 30

	// 调用函数（带超时）
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err = engine.CallFunction(ctx, "add", 100, 200)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(result) // 输出: 300
}
