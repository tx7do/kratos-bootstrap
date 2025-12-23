package script_engine

import (
	"context"

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

	pool, err := scriptEngine.NewAutoGrowEnginePool(
		int(cfg.GetPool().GetInitial().Value),
		int(cfg.GetPool().GetMax().Value),
		typ,
	)
	if err != nil {
		return nil, err
	}

	return pool, nil
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

	pool, err := scriptEngine.NewEnginePool(int(cfg.GetPool().GetMax().Value), typ)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

// NewScriptEngine 创建脚本引擎
func NewScriptEngine(ctx context.Context, cfg *conf.Script) (scriptEngine.Engine, error) {
	if cfg == nil {
		return nil, nil
	}

	// convert type
	typ := convertTypeScriptEngineType(cfg.GetEngine())
	if typ == "" {
		return nil, nil
	}

	// create engine
	eng, err := scriptEngine.NewScriptEngine(typ)
	if err != nil {
		return nil, err
	}

	// initialize engine
	if err = eng.Init(ctx); err != nil {
		return nil, err
	}

	// preload scripts
	if err = PreLoadScripts(ctx, cfg, typ, eng); err != nil {
		return nil, err
	}

	return eng, nil
}

// PreLoadScripts 预加载脚本文件
func PreLoadScripts(ctx context.Context, cfg *conf.Script, typ scriptEngine.Type, eng scriptEngine.Engine) error {
	var entryScriptFilePath string
	var preloadScripts []string
	switch typ {
	case scriptEngine.JavaScriptType:
		jsCfg := cfg.GetJavascript()
		entryScriptFilePath = jsCfg.GetEntry().Value
		preloadScripts = jsCfg.GetPreLoadScripts()
	case scriptEngine.LuaType:
		luaCfg := cfg.GetLua()
		entryScriptFilePath = luaCfg.GetEntry().Value
		preloadScripts = luaCfg.GetPreLoadScripts()
	}

	// execute preload script files
	if len(preloadScripts) > 0 {
		if err := eng.LoadFiles(ctx, preloadScripts); err != nil {
			return err
		}
	}

	// execute entry script file
	if entryScriptFilePath != "" {
		if err := eng.LoadFile(ctx, entryScriptFilePath); err != nil {
			return err
		}
	}

	return nil
}
