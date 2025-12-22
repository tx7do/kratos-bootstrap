package script_engine

import (
	scriptEngine "github.com/tx7do/go-scripts"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewAutoGrowScriptEnginePool 创建自动增长的脚本引擎池
func NewAutoGrowScriptEnginePool(cfg *conf.Script) (*scriptEngine.AutoGrowEnginePool, error) {
	if cfg == nil {
		return nil, nil
	}

	typ := convertTypeScriptEngineType(cfg.GetEngine())
	if typ == "" {
		return nil, nil
	}

	return scriptEngine.NewAutoGrowEnginePool(int(cfg.GetPool().GetInitial().Value), int(cfg.GetPool().GetMax().Value), typ)
}

// NewScriptEnginePool 创建固定大小的脚本引擎池
func NewScriptEnginePool(cfg *conf.Script) (*scriptEngine.EnginePool, error) {
	if cfg == nil {
		return nil, nil
	}

	typ := convertTypeScriptEngineType(cfg.GetEngine())
	if typ == "" {
		return nil, nil
	}

	return scriptEngine.NewEnginePool(int(cfg.GetPool().GetMax().Value), typ)
}
