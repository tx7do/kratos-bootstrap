package script_engine

import (
	"context"
	"fmt"

	scriptEngine "github.com/tx7do/go-scripts"
	"github.com/tx7do/go-scripts/source"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewScriptEngine 根据配置创建并初始化单个脚本引擎实例。
//
// 流程:
//  1. 转换引擎类型，校验配置
//  2. 检查 enabled 标志（未设置或为 true 时继续）
//  3. 通过工厂创建引擎实例
//  4. 调用 Init 初始化引擎
//  5. 创建并绑定 Source（如果配置了）
//  6. 预加载脚本（pre_load_scripts + entry）
//  7. 启动热加载（如果配置了 hot_reload）
//
// 返回 (nil, nil) 表示配置为 nil 或引擎未启用。
func NewScriptEngine(ctx context.Context, cfg *conf.Script) (scriptEngine.Engine, error) {
	if cfg == nil {
		return nil, nil
	}

	// 转换引擎类型
	typ := convertEngineType(cfg.GetEngine())
	if typ == "" {
		return nil, nil
	}

	// 检查 enabled
	opts := cfg.GetOptions()
	if !isEnabled(opts) {
		return nil, nil
	}

	// 创建引擎
	eng, err := scriptEngine.NewScriptEngine(typ)
	if err != nil {
		return nil, fmt.Errorf("script engine: create %s failed: %w", typ, err)
	}

	// 初始化
	if err = eng.Init(ctx); err != nil {
		_ = eng.Close()
		return nil, fmt.Errorf("script engine: init %s failed: %w", typ, err)
	}

	// 创建并绑定 Source
	src, err := createSource(cfg.GetSource())
	if err != nil {
		_ = eng.Close()
		return nil, err
	}
	if src != nil {
		eng.SetSource(src)
	}

	// 预加载脚本
	if err = loadScriptsToEngine(ctx, eng, opts, src); err != nil {
		_ = eng.Close()
		return nil, err
	}

	// 热加载
	startHotReload(ctx, eng, opts)

	return eng, nil
}

// NewScriptEnginePool 创建固定大小的脚本引擎池。
//
// 池中每个引擎实例都会执行与 NewScriptEngine 相同的初始化流程：
// Init → SetSource → PreLoad → HotReload。
//
// 注意：引擎池的 per-call 代理方法（如 RegisterFunction）只作用于被 Acquire 的单个实例。
// 如需池级别统一绑定，请在获取池后手动遍历 Acquire/Release。
//
// 返回 (nil, nil) 表示配置为 nil 或引擎未启用。
func NewScriptEnginePool(ctx context.Context, cfg *conf.Script) (*scriptEngine.EnginePool, error) {
	if cfg == nil {
		return nil, nil
	}

	typ := convertEngineType(cfg.GetEngine())
	if typ == "" {
		return nil, nil
	}

	opts := cfg.GetOptions()
	if !isEnabled(opts) {
		return nil, nil
	}

	// 获取池大小
	poolCfg := cfg.GetPool()
	size := getPoolSize(poolCfg)

	// 创建池
	pool, err := scriptEngine.NewEnginePool(size, typ)
	if err != nil {
		return nil, fmt.Errorf("script engine: create pool failed: %w", err)
	}

	// 为池中每个引擎实例配置 Source 和预加载脚本
	if err = configurePoolEngines(ctx, pool, cfg, size); err != nil {
		_ = pool.Close()
		return nil, err
	}

	return pool, nil
}

// NewAutoGrowScriptEnginePool 创建自动扩容的脚本引擎池。
//
// 初始创建 initial 个实例，可按需扩容到 max 个。
// 每个实例都会执行完整的初始化流程。
//
// 返回 (nil, nil) 表示配置为 nil 或引擎未启用。
func NewAutoGrowScriptEnginePool(ctx context.Context, cfg *conf.Script) (*scriptEngine.AutoGrowEnginePool, error) {
	if cfg == nil {
		return nil, nil
	}

	typ := convertEngineType(cfg.GetEngine())
	if typ == "" {
		return nil, nil
	}

	opts := cfg.GetOptions()
	if !isEnabled(opts) {
		return nil, nil
	}

	// 获取池大小
	poolCfg := cfg.GetPool()
	initial, max := getAutoGrowPoolSize(poolCfg)

	// 创建池
	pool, err := scriptEngine.NewAutoGrowEnginePool(initial, max, typ)
	if err != nil {
		return nil, fmt.Errorf("script engine: create auto-grow pool failed: %w", err)
	}

	// 为初始实例配置 Source 和预加载脚本
	if err = configureAutoGrowPoolEngines(ctx, pool, cfg, initial); err != nil {
		_ = pool.Close()
		return nil, err
	}

	return pool, nil
}

// configurePoolEngines 为固定池中的每个引擎实例配置 Source 和预加载脚本。
func configurePoolEngines(ctx context.Context, pool *scriptEngine.EnginePool, cfg *conf.Script, size int) error {
	opts := cfg.GetOptions()
	srcCfg := cfg.GetSource()

	// 预创建 Source（MEMORY 类型需每个实例独立，FILE 类型可共享）
	engines := make([]scriptEngine.Engine, 0, size)
	for i := 0; i < size; i++ {
		eng, err := pool.Acquire()
		if err != nil {
			break
		}
		engines = append(engines, eng)
	}

	for _, eng := range engines {
		// 为每个实例创建独立的 Source（MEMORY/EMBED 需要独立状态）
		src, err := createSource(srcCfg)
		if err != nil {
			return err
		}
		if src != nil {
			eng.SetSource(src)
		}

		if err := loadScriptsToEngine(ctx, eng, opts, src); err != nil {
			return err
		}

		startHotReload(ctx, eng, opts)
	}

	// 释放回池
	for _, eng := range engines {
		pool.Release(eng)
	}

	return nil
}

// configureAutoGrowPoolEngines 为自动扩容池的初始实例配置 Source 和预加载脚本。
func configureAutoGrowPoolEngines(ctx context.Context, pool *scriptEngine.AutoGrowEnginePool, cfg *conf.Script, initial int) error {
	opts := cfg.GetOptions()
	srcCfg := cfg.GetSource()

	engines := make([]scriptEngine.Engine, 0, initial)
	for i := 0; i < initial; i++ {
		eng, err := pool.Acquire()
		if err != nil {
			break
		}
		engines = append(engines, eng)
	}

	for _, eng := range engines {
		src, err := createSource(srcCfg)
		if err != nil {
			return err
		}
		if src != nil {
			eng.SetSource(src)
		}

		if err := loadScriptsToEngine(ctx, eng, opts, src); err != nil {
			return err
		}

		startHotReload(ctx, eng, opts)
	}

	for _, eng := range engines {
		pool.Release(eng)
	}

	return nil
}

// getPoolSize 从配置中解析固定池大小，默认为 4。
func getPoolSize(cfg *conf.Script_Pool) int {
	if cfg == nil {
		return 4
	}
	if m := cfg.GetMax(); m != nil && m.Value > 0 {
		return int(m.Value)
	}
	return 4
}

// getAutoGrowPoolSize 从配置中解析自动扩容池的 initial 和 max，默认 initial=2 max=8。
func getAutoGrowPoolSize(cfg *conf.Script_Pool) (initial, max int) {
	initial, max = 2, 8
	if cfg == nil {
		return
	}
	if i := cfg.GetInitial(); i != nil && i.Value > 0 {
		initial = int(i.Value)
	}
	if m := cfg.GetMax(); m != nil && m.Value > 0 {
		max = int(m.Value)
	}
	// 确保 initial <= max
	if initial > max {
		initial = max
	}
	return
}

// Compile-time assertions: ensure source.Reader is imported and used.
var _ source.Reader = (*source.FileSource)(nil)
