package script_engine

import (
	scriptEngine "github.com/tx7do/go-scripts"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func NewScriptEngine(cfg *conf.Script) (*scriptEngine.AutoGrowEnginePool, error) {
	if cfg == nil {
		return nil, nil
	}

	typ := convertTypeScriptEngineType(cfg.GetEngine())
	if typ == "" {
		return nil, nil
	}

	return scriptEngine.NewAutoGrowEnginePool(int(cfg.GetPool().GetInitial().Value), int(cfg.GetPool().GetMax().Value), typ)
}

// convertTypeScriptEngineType .
func convertTypeScriptEngineType(t conf.Script_EngineType) scriptEngine.Type {
	switch t {
	case conf.Script_LUA:
		return scriptEngine.LuaType
	case conf.Script_JAVASCRIPT:
		return scriptEngine.JavaScriptType
	default:
		return ""
	}
}
