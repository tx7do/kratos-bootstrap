package script_engine

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	scriptEngine "github.com/tx7do/go-scripts"
	"github.com/tx7do/go-scripts/source"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

////////////////////////////////////////////////////////////////////////////////
// Embed FS Provider — 用于 go:embed 脚本加载
////////////////////////////////////////////////////////////////////////////////

// EmbedFSProvider 是用于提供 go:embed fs.FS 的回调函数类型。
// 由于 embed.FS 无法通过 protobuf 配置传递，需要宿主程序在代码中注册。
type EmbedFSProvider func() (map[string]string, error)

// embedFSProvider 全局 embed FS 提供者，由宿主程序通过 SetEmbedFSProvider 设置。
var embedFSProvider EmbedFSProvider

// SetEmbedFSProvider 注册全局 embed FS 提供者。
// 宿主程序应在初始化脚本引擎之前调用此函数。
//
// 示例:
//
//	//go:embed scripts/*
//	var embedFS embed.FS
//
//	script_engine.SetEmbedFSProvider(func() (map[string]string, error) {
//	    // 返回 prefix -> 脚本内容的映射
//	})
func SetEmbedFSProvider(p EmbedFSProvider) {
	embedFSProvider = p
}

////////////////////////////////////////////////////////////////////////////////
// Source Factory Registry — 可扩展的来源创建机制
////////////////////////////////////////////////////////////////////////////////

// SourceFactoryFunc 是自定义 Source 创建函数的签名。
// cfg 参数是完整的 Source 配置（type + paths + strategy + options），
// 工厂实现应根据 options 中的特有配置创建对应的 source.Reader。
//
// 内置类型（FILE / MEMORY / EMBED / MULTI）由本包自动处理，无需注册。
// 扩展类型（S3 / ETCD / CONSUL / REDIS / HTTP / GIT / DATABASE）需要宿主程序注册。
type SourceFactoryFunc func(cfg *conf.Script_Source) (source.Reader, error)

var (
	sourceFactoryMu sync.RWMutex
	sourceFactories = make(map[conf.Script_Source_Type]SourceFactoryFunc)
)

// RegisterSourceFactory 注册自定义 Source 创建工厂。
// 如果 typ 对应的工厂已存在，返回错误。
//
// 用于扩展来源类型（S3、etcd、Consul、Redis、HTTP、Git、Database 等），
// 当 go-scripts 发布对应子包或用户自行实现时调用。
//
// 示例 —— 注册 etcd Source 工厂（go-scripts/source/etcd 发布后）：
//
//	func init() {
//	    _ = script_engine.RegisterSourceFactory(conf.Script_Source_ETCD,
//	        func(cfg *conf.Script_Source) (source.Reader, error) {
//	            opts := cfg.GetOptions()
//	            endpoints := opts.GetFields()["endpoints"].GetListValue()
//	            // ... 解析配置并创建 etcd source
//	            return etcdSource.New(endpoints...)
//	        })
//	}
func RegisterSourceFactory(typ conf.Script_Source_Type, f SourceFactoryFunc) error {
	if f == nil {
		return fmt.Errorf("script engine: source factory function is nil")
	}
	sourceFactoryMu.Lock()
	defer sourceFactoryMu.Unlock()
	if _, ok := sourceFactories[typ]; ok {
		return fmt.Errorf("script engine: source factory for %s already registered", typ)
	}
	sourceFactories[typ] = f
	return nil
}

// MustRegisterSourceFactory 是 RegisterSourceFactory 的 panic 版本，用于 init() 中。
func MustRegisterSourceFactory(typ conf.Script_Source_Type, f SourceFactoryFunc) {
	if err := RegisterSourceFactory(typ, f); err != nil {
		panic(err)
	}
}

// getSourceFactory 返回已注册的工厂及其是否存在。
func getSourceFactory(typ conf.Script_Source_Type) (SourceFactoryFunc, bool) {
	sourceFactoryMu.RLock()
	defer sourceFactoryMu.RUnlock()
	f, ok := sourceFactories[typ]
	return f, ok
}

////////////////////////////////////////////////////////////////////////////////
// Source Creation — 核心创建逻辑
////////////////////////////////////////////////////////////////////////////////

// createSource 根据配置创建脚本来源 (source.Reader)。
// 如果 cfg 为 nil 或类型未指定，返回 nil（表示不绑定 Source）。
//
// 创建顺序：
//  1. 内置类型（FILE / MEMORY / EMBED / MULTI）在主包中直接处理
//  2. 扩展类型（S3 / ETCD / ...）统一通过工厂注册表创建
//     —— 每个扩展来源是独立 go.mod 子模块，通过 init() 注册工厂
//  3. 如果配置了 cache_ttl，自动包裹 CachedSource 缓存层
func createSource(cfg *conf.Script_Source) (source.Reader, error) {
	if cfg == nil {
		return nil, nil
	}

	var src source.Reader
	var err error

	switch cfg.GetType() {
	// ---- 内置来源（go-scripts v0.0.6 核心，无额外依赖）----
	case conf.Script_Source_FILE:
		src = source.NewFileSource()

	case conf.Script_Source_MEMORY:
		src = source.NewMemSource()

	case conf.Script_Source_EMBED:
		src, err = createEmbedSource()

	case conf.Script_Source_MULTI:
		src, err = createMultiSource(cfg)

	case conf.Script_Source_TYPE_UNSPECIFIED:
		// 未指定类型时默认不绑定 Source
		return nil, nil

	// ---- 扩展来源（独立 go.mod 子模块，通过工厂注册表）----
	// S3 / ETCD / CONSUL / REDIS / HTTP / GIT / DATABASE 及任何自定义类型
	// 用户 import 对应子模块后，其 init() 自动注册工厂
	default:
		src, err = createSourceFromFactory(cfg)
	}

	if err != nil {
		return nil, err
	}
	if src == nil {
		return nil, nil
	}

	// 可选包装：CachedSource（缓存层）
	src, err = wrapWithCache(src, cfg)
	if err != nil {
		return nil, err
	}

	return src, nil
}

// createSourceFromFactory 通过注册的工厂创建 Source。
// 如果没有注册对应类型的工厂，返回错误。
//
// 扩展来源是独立 go.mod 子模块，用户需 import 对应模块：
//
//	import _ "github.com/tx7do/kratos-bootstrap/script_engine/source/etcd"
//
// 该模块的 init() 会自动调用 MustRegisterSourceFactory 完成注册。
func createSourceFromFactory(cfg *conf.Script_Source) (source.Reader, error) {
	f, ok := getSourceFactory(cfg.GetType())
	if !ok {
		return nil, fmt.Errorf(
			"script engine: source type %s not registered — "+
				"import the corresponding source sub-module (e.g. _ \"github.com/tx7do/kratos-bootstrap/script_engine/source/etcd\")",
			cfg.GetType(),
		)
	}
	return f(cfg)
}

// createEmbedSource 从 EmbedFSProvider 创建内存 Source。
func createEmbedSource() (source.Reader, error) {
	if embedFSProvider == nil {
		return nil, fmt.Errorf("script engine: embed source requires SetEmbedFSProvider() to be called first")
	}
	memSrc := source.NewMemSource()
	scripts, err := embedFSProvider()
	if err != nil {
		return nil, fmt.Errorf("script engine: embed fs provider error: %w", err)
	}
	for key, code := range scripts {
		memSrc.Set(key, code)
	}
	return memSrc, nil
}

// createMultiSource 创建多源聚合 Source。
//
// 子源来源有两种方式：
//  1. 从 paths 创建：每个 path 创建一个 FileSource（简单场景）
//  2. 从 options.sources 创建：每个元素可指定不同来源类型（异构聚合）
//
// 优先使用 options.sources（如果存在），否则回退到 paths。
func createMultiSource(cfg *conf.Script_Source) (source.Reader, error) {
	var subSources []source.Reader

	// 方式 1：从 options.sources 解析异构子源
	if opts := cfg.GetOptions(); opts != nil {
		if sourcesField, ok := opts.AsMap()["sources"]; ok {
			parsed, err := parseMultiSubSources(sourcesField)
			if err != nil {
				return nil, fmt.Errorf("script engine: parse multi sources failed: %w", err)
			}
			subSources = parsed
		}
	}

	// 方式 2：回退到 paths（每个 path 一个 FileSource）
	if len(subSources) == 0 {
		paths := cfg.GetPaths()
		if len(paths) == 0 {
			return nil, fmt.Errorf("script engine: multi source requires options.sources or paths")
		}
		subSources = make([]source.Reader, 0, len(paths))
		for range paths {
			subSources = append(subSources, source.NewFileSource())
		}
	}

	if len(subSources) == 0 {
		return nil, fmt.Errorf("script engine: multi source has no sub-sources")
	}

	strategy := source.MultiStrategyFallback
	if cfg.GetStrategy() == conf.Script_Source_FIRST_OK {
		strategy = source.MultiStrategyFirstOK
	}

	return source.NewMultiSource(strategy, subSources...)
}

// parseMultiSubSources 从 options.sources 配置解析子源列表。
// sources 格式为 ListValue，每个元素是一个 Struct，包含：
//
//	{ "type": "FILE", "paths": ["./scripts"] }
//	{ "type": "S3", "bucket": "my-scripts", "prefix": "lua/" }
func parseMultiSubSources(raw any) ([]source.Reader, error) {
	list, ok := raw.([]any)
	if !ok {
		return nil, fmt.Errorf("sources must be a list")
	}

	subSources := make([]source.Reader, 0, len(list))
	for i, item := range list {
		subCfg, err := parseSubSourceConfig(item)
		if err != nil {
			return nil, fmt.Errorf("source[%d]: %w", i, err)
		}

		src, err := createSource(subCfg)
		if err != nil {
			return nil, fmt.Errorf("source[%d]: %w", i, err)
		}
		if src == nil {
			return nil, fmt.Errorf("source[%d]: created nil source", i)
		}
		subSources = append(subSources, src)
	}
	return subSources, nil
}

// parseSubSourceConfig 从 any 类型的配置解析为 Script_Source。
func parseSubSourceConfig(raw any) (*conf.Script_Source, error) {
	m, ok := raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("must be a map, got %T", raw)
	}

	subCfg := &conf.Script_Source{}

	// 解析 type
	if typeStr, ok := m["type"].(string); ok {
		subCfg.Type = parseSourceType(typeStr)
	}
	if subCfg.Type == conf.Script_Source_TYPE_UNSPECIFIED {
		subCfg.Type = conf.Script_Source_FILE // 默认 FILE
	}

	// 解析 paths
	if paths, ok := m["paths"].([]any); ok {
		subCfg.Paths = anySliceToStrings(paths)
	}

	// 解析 strategy
	if strategyStr, ok := m["strategy"].(string); ok {
		subCfg.Strategy = parseStrategy(strategyStr)
	}

	return subCfg, nil
}

// parseSourceType 将字符串解析为 Script_Source_Type。
func parseSourceType(s string) conf.Script_Source_Type {
	switch strings.ToUpper(s) {
	case "FILE":
		return conf.Script_Source_FILE
	case "MEMORY":
		return conf.Script_Source_MEMORY
	case "EMBED":
		return conf.Script_Source_EMBED
	case "MULTI":
		return conf.Script_Source_MULTI
	case "S3":
		return conf.Script_Source_S3
	case "ETCD":
		return conf.Script_Source_ETCD
	case "CONSUL":
		return conf.Script_Source_CONSUL
	case "REDIS":
		return conf.Script_Source_REDIS
	case "HTTP":
		return conf.Script_Source_HTTP
	case "GIT":
		return conf.Script_Source_GIT
	case "DATABASE":
		return conf.Script_Source_DATABASE
	default:
		return conf.Script_Source_TYPE_UNSPECIFIED
	}
}

// parseStrategy 将字符串解析为 Script_Source_MultiStrategy。
func parseStrategy(s string) conf.Script_Source_MultiStrategy {
	switch strings.ToUpper(s) {
	case "FALLBACK":
		return conf.Script_Source_FALLBACK
	case "FIRST_OK", "FIRSTOK":
		return conf.Script_Source_FIRST_OK
	default:
		return conf.Script_Source_MULTI_STRATEGY_UNSPECIFIED
	}
}

// anySliceToStrings 将 []any 转为 []string。
func anySliceToStrings(items []any) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		if s, ok := item.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

// wrapWithCache 可选地将 Source 包装为 CachedSource。
// 仅当 options 中配置了 cache_ttl 时生效。
//
// YAML 示例:
//
//	source:
//	  type: S3
//	  options:
//	    cache_ttl: "5m"  # 启用缓存，TTL 5 分钟
func wrapWithCache(src source.Reader, cfg *conf.Script_Source) (source.Reader, error) {
	opts := cfg.GetOptions()
	if opts == nil {
		return src, nil
	}
	m := opts.AsMap()
	ttlStr, ok := m["cache_ttl"].(string)
	if !ok || ttlStr == "" {
		return src, nil
	}

	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		return nil, fmt.Errorf("script engine: invalid cache_ttl %q: %w", ttlStr, err)
	}

	return source.NewCachedSource(src, source.WithTTL(ttl))
}

// resolveKey 根据搜索路径列表解析脚本 key。
// 如果 key 是绝对路径，直接返回；
// 否则依次尝试 paths 中的目录前缀，返回拼接后的路径。
// 如果 paths 为空或均不匹配，返回原始 key。
func resolveKey(paths []string, key string) string {
	if key == "" {
		return key
	}
	// 绝对路径直接返回
	if filepath.IsAbs(key) {
		return key
	}
	// 无搜索路径时返回原始 key
	if len(paths) == 0 {
		return key
	}
	// 尝试每个搜索路径
	for _, p := range paths {
		candidate := filepath.Join(p, key)
		// 标准化路径，去除冗余分隔符
		candidate = filepath.Clean(candidate)
		// 直接返回第一个拼接结果（Source.Load 会处理文件不存在的情况）
		return candidate
	}
	return key
}

// resolveKeys 批量解析脚本 key。
func resolveKeys(paths []string, keys []string) []string {
	if len(paths) == 0 {
		return keys
	}
	resolved := make([]string, len(keys))
	for i, k := range keys {
		resolved[i] = resolveKey(paths, k)
	}
	return resolved
}

// loadScriptsToEngine 将预加载脚本和入口脚本加载到引擎中。
//
// 加载流程:
//  1. 如果引擎支持 ScriptLoader 能力，通过 Source 加载 pre_load_scripts
//  2. 如果引擎支持 ScriptLoader 能力，通过 Source 加载 entry 脚本
//  3. 如果引擎不支持 ScriptLoader（轻量引擎），则跳过文件加载
func loadScriptsToEngine(ctx context.Context, eng any, opts *conf.Script_EngineOptions, src source.Reader) error {
	if opts == nil {
		return nil
	}

	loader := scriptEngine.AsLoader(eng)
	if loader == nil {
		// 轻量引擎（CEL/Expr 等）不支持文件加载，跳过
		return nil
	}

	paths := opts.GetPaths()

	// 加载预加载脚本列表
	preScripts := opts.GetPreLoadScripts()
	if len(preScripts) > 0 {
		keys := resolveKeys(paths, preScripts)
		if src != nil {
			if err := loader.LoadMulti(ctx, keys); err != nil {
				return fmt.Errorf("script engine: preload scripts failed: %w", err)
			}
		} else {
			// 无 Source 绑定时，尝试用 LoadString 逐个加载
			for i, key := range keys {
				code, err := readFileForLoad(key)
				if err != nil {
					return fmt.Errorf("script engine: read preload script %q: %w", preScripts[i], err)
				}
				if err := loader.LoadString(ctx, preScripts[i], code); err != nil {
					return fmt.Errorf("script engine: load preload script %q: %w", preScripts[i], err)
				}
			}
		}
	}

	// 加载入口脚本
	if entry := opts.GetEntry(); entry != nil && entry.Value != "" {
		entryKey := resolveKey(paths, entry.Value)
		if src != nil {
			if err := loader.Load(ctx, entryKey); err != nil {
				return fmt.Errorf("script engine: load entry script %q: %w", entry.Value, err)
			}
		} else {
			code, err := readFileForLoad(entryKey)
			if err != nil {
				return fmt.Errorf("script engine: read entry script %q: %w", entry.Value, err)
			}
			if err := loader.LoadString(ctx, entry.Value, code); err != nil {
				return fmt.Errorf("script engine: load entry script %q: %w", entry.Value, err)
			}
		}
	}

	return nil
}

// startHotReload 为入口脚本和预加载脚本启动热加载监听。
// 仅当引擎支持 ScriptWatcher 能力且 Source 支持 Watcher 时生效。
func startHotReload(ctx context.Context, eng any, opts *conf.Script_EngineOptions) {
	if opts == nil {
		return
	}
	if !isHotReloadEnabled(opts) {
		return
	}

	watcher := scriptEngine.AsWatcher(eng)
	if watcher == nil {
		return
	}

	paths := opts.GetPaths()

	// 监听入口脚本
	if entry := opts.GetEntry(); entry != nil && entry.Value != "" {
		entryKey := resolveKey(paths, entry.Value)
		_ = watcher.StartWatch(ctx, entryKey)
	}

	// 监听预加载脚本
	for _, script := range opts.GetPreLoadScripts() {
		key := resolveKey(paths, script)
		_ = watcher.StartWatch(ctx, key)
	}
}

// isHotReloadEnabled 检查是否启用了热加载。
func isHotReloadEnabled(opts *conf.Script_EngineOptions) bool {
	if opts == nil {
		return false
	}
	if hr := opts.GetHotReload(); hr != nil {
		return hr.Value
	}
	return false
}

// isEnabled 检查引擎是否启用。
// 未设置 enabled 字段时默认启用（presence 语义）。
func isEnabled(opts *conf.Script_EngineOptions) bool {
	if opts == nil {
		return true // 无 options 默认启用
	}
	if en := opts.GetEnabled(); en != nil {
		return en.Value
	}
	return true // 未设置默认启用
}

// readFileForLoad 是无 Source 时的回退方案，直接读取文件内容。
// 仅用于未配置 Source 但仍提供了脚本路径的场景。
func readFileForLoad(path string) (string, error) {
	// 委托给 FileSource 实现，避免重复 IO 逻辑
	src := source.NewFileSource()
	defer src.Close()
	ctx := context.Background()
	code, err := src.Load(ctx, path)
	if err != nil {
		return "", err
	}
	return code, nil
}
