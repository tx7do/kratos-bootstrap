# Script Engine

This directory contains the implementation of a script engine that supports multiple scripting languages. The engine allows users to write and execute scripts in different languages seamlessly.

## Support Script Languages

- Lua
- JavaScript

## Example Usage

```go
import (
    "context"
    "fmt"
    
    _ "github.com/tx7do/go-scripts/javascript"
    conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
    "github.com/tx7do/kratos-bootstrap/script_engine"
    "google.golang.org/protobuf/types/known/wrapperspb"
)

	var cfg *conf.Script
	cfg = &conf.Script{
		Engine: conf.Script_JAVASCRIPT,
		Pool: &conf.Script_Pool{
			Initial: &wrapperspb.Int32Value{Value: 2},
			Max:     &wrapperspb.Int32Value{Value: 10},
		},
	}

	// initialize script engine pool
	enginePool, err := script_engine.NewAutoGrowScriptEnginePool(cfg)
	if err != nil {
		// handle initialization error
	}
	defer enginePool.Close()

	// acquire a script engine instance from the pool
	eng, err := enginePool.Acquire()
	if err != nil {
		// handle error acquiring engine
		return
	}

	// register to the script engine for JavaScript to call
	err = eng.RegisterFunction("updateUserStatus", updateUserStatus)
	if err != nil {
		// handle registration error
	}

	// execute script, define permission check function
	ret, err := eng.ExecuteString(context.Background(), `
function checkPermission(userId, permission) {
  // simulate permission check logic
  const userPermissions = ["order:view", "order:edit"];
  return userPermissions.includes(permission);
}
`)
	if err != nil {
		// handle execution error
	}

	result, err := eng.CallFunction(
		context.Background(),
		"checkPermission",
		10086,        // user ID
		"order:edit", // permission identifier
	)
	if err != nil {
		// handle call error
	}
	// print permission check result
	fmt.Printf("Permission check: %v\n", result)
```