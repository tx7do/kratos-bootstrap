package script_engine

import (
	scriptEngine "github.com/tx7do/go-scripts"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// convertEngineType 转换 conf.Script_EngineType 到 scriptEngine.Type。
// 返回空字符串表示未指定或不支持的引擎类型。
func convertEngineType(t conf.Script_EngineType) scriptEngine.Type {
	switch t {
	case conf.Script_LUA:
		return scriptEngine.LuaType
	case conf.Script_JAVASCRIPT:
		return scriptEngine.JavaScriptType
	case conf.Script_GPYTHON:
		return scriptEngine.GPythonType
	case conf.Script_YAEGI:
		return scriptEngine.YaegiType
	case conf.Script_WAZERO:
		return scriptEngine.WazeroType
	case conf.Script_CEL:
		return scriptEngine.CELType
	case conf.Script_EXPR:
		return scriptEngine.ExprType
	case conf.Script_STARLARK:
		return scriptEngine.StarlarkType
	case conf.Script_TCL:
		return scriptEngine.TclType
	default:
		return ""
	}
}

// isFullEngine 报告给定引擎类型是否支持完整能力接口（ScriptLoader + ScriptExecutor +
// GlobalAccessor + FunctionRegistrar + ModuleRegistrar + ScriptWatcher）。
//
// 完整引擎：Lua、JavaScript
// 轻量引擎（仅 ScriptEngine + ScriptExecutor）：CEL、Expr、Starlark 等
func isFullEngine(typ scriptEngine.Type) bool {
	switch typ {
	case scriptEngine.LuaType,
		scriptEngine.JavaScriptType:
		return true
	default:
		return false
	}
}
