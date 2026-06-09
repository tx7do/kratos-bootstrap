# Script Engine

<p align="center">
  <a href="README.md">中文</a> · <a href="README_en.md">English</a> · <a href="README_ja.md">日本語</a>
</p>

基于 [go-scripts](https://github.com/tx7do/go-scripts) 的多语言嵌入式脚本引擎封装，为 Kratos 应用提供统一的脚本执行能力。

## 支持的脚本引擎（9 种）

| 引擎 | 枚举值 | 底层库 | 适用场景 |
| --- | --- | --- | --- |
| **Lua** | `LUA` | gopher-lua | 游戏脚本、配置逻辑、嵌入式扩展 |
| **JavaScript** | `JAVASCRIPT` | goja | 前端复用、规则引擎、快速原型 |
| **Python** | `GPYTHON` | gpython | 数据处理、运维脚本、算法验证 |
| **Go** | `YAEGI` | Yaegi | 动态插件、DevOps 工具链 |
| **WebAssembly** | `WAZERO` | wazero | 高性能沙箱、跨语言模块复用 |
| **CEL** | `CEL` | cel-go | 策略引擎、权限规则、条件判断 |
| **Expr** | `EXPR` | expr-lang | 业务表达式、模板引擎、数据筛选 |
| **Starlark** | `STARLARK` | starlark-go | 构建工具、安全脚本、Bazel 规则 |
| **TCL** | `TCL` | modernc/tcl | 传统系统兼容、网络设备脚本 |

> Lua 和 JavaScript 是**完整引擎**（支持脚本加载、执行、函数注册、模块注册、热更新）。
> 其他引擎为**轻量引擎**（仅支持核心生命周期 + 脚本执行）。

## 核心功能

- **Source 模式**：通过 `FileSource`、`MemSource`、`EmbedSource`、`MultiSource` 统一脚本加载
- **扩展来源**：S3、etcd、Consul、Redis、HTTP、Git、Database 等 7 种扩展来源，通过工厂注册机制接入
- **缓存层**：远程源（S3/etcd 等）自动包裹 `CachedSource`，减少网络 IO
- **引擎池**：固定大小池 `EnginePool` + 自动扩容池 `AutoGrowEnginePool`
- **热更新**：配置 `hot_reload: true` 即可启用脚本变更自动重载
- **Kratos 日志集成**：`KratosLogger` 适配器无缝接入 Kratos 统一日志体系

## 配置示例

### YAML 配置

```yaml
script:
  engine: JAVASCRIPT
  options:
    enabled: true
    paths:
      - "./scripts"
    entry:
      value: "main.js"
    pre_load_scripts:
      - "utils.js"
      - "config.js"
    hot_reload:
      value: true
  source:
    type: FILE
    paths:
      - "./scripts"
  pool:
    initial:
      value: 2
    max:
      value: 10
```

### 使用 go:embed 嵌入脚本

```yaml
script:
  engine: LUA
  source:
    type: EMBED
  options:
    entry:
      value: "main.lua"
    pre_load_scripts:
      - "lib.lua"
```

```go
//go:embed scripts/*.lua
var embedFS embed.FS

func main() {
    // 注册 embed FS 提供者
    script_engine.SetEmbedFSProvider(func() (map[string]string, error) {
        scripts := make(map[string]string)
        entries, _ := embedFS.ReadDir("scripts")
        for _, entry := range entries {
            if entry.IsDir() { continue }
            data, _ := embedFS.ReadFile("scripts/" + entry.Name())
            scripts[entry.Name()] = string(data)
        }
        return scripts, nil
    })

    // 创建引擎（自动从 embed 加载脚本）
    eng, _ := script_engine.NewScriptEngine(ctx, cfg)
}
```

## 示例用法

### 单引擎实例

```go
import (
    "context"
    "fmt"

    _ "github.com/tx7do/go-scripts/javascript"
    conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
    "github.com/tx7do/kratos-bootstrap/script_engine"
)

func main() {
    cfg := &conf.Script{
        Engine: conf.Script_JAVASCRIPT,
        Options: &conf.Script_EngineOptions{
            Paths: []string{"./scripts"},
            Entry: &wrapperspb.StringValue{Value: "main.js"},
        },
        Source: &conf.Script_Source{
            Type: conf.Script_Source_FILE,
        },
    }

    ctx := context.Background()
    eng, err := script_engine.NewScriptEngine(ctx, cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer eng.Close()

    // 注册宿主函数
    _ = eng.RegisterFunction("greet", func(name string) string {
        return fmt.Sprintf("Hello, %s!", name)
    })

    // 执行入口脚本
    result, err := eng.Execute(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result)
}
```

### 引擎池（生产推荐）

```go
// 自动扩容池（初始 2，上限 16）
pool, err := script_engine.NewAutoGrowEnginePool(ctx, cfg)
if err != nil {
    log.Fatal(err)
}
defer pool.Close()

// 通过池执行脚本（自动 Acquire/Release）
result, err := pool.ExecuteString(ctx, "hello.js", `greet("world")`)
```

### 集成 Kratos 日志

```go
import (
    "github.com/go-kratos/kratos/v2/log"
    "github.com/tx7do/kratos-bootstrap/script_engine"
)

// 将 Kratos Logger 注入 go-scripts
script_engine.SetLogger(kratosLogger)
```

## 文件结构

| 文件 / 目录 | 说明 |
| --- | --- |
| `script_engine.go` | 核心封装：NewScriptEngine / NewEnginePool / NewAutoGrowEnginePool |
| `source.go` | Source 工厂注册表、来源创建、路径解析、缓存包装、热更新 |
| `logger.go` | Kratos Logger → go-scripts Logger 适配器 |
| `utils.go` | 引擎类型映射、能力检测 |
| `source/etcd/` | etcd 来源子模块（独立 go.mod） |
| `source/s3/` | S3 来源子模块（独立 go.mod） |
| `source/redis/` | Redis 来源子模块（独立 go.mod） |
| `source/consul/` | Consul 来源子模块（独立 go.mod） |
| `source/http/` | HTTP 来源子模块（独立 go.mod） |
| `source/git/` | Git 来源子模块（独立 go.mod） |
| `source/database/` | SQL 数据库来源子模块（独立 go.mod） |

## 脚本来源（Source）

### 内置来源（go-scripts v0.0.6 已支持）

| 类型 | 枚举值 | 热更新 | 说明 |
| --- | --- | --- | --- |
| **FILE** | `FILE` | mtime 轮询 | 本地文件系统 |
| **MEMORY** | `MEMORY` | channel 通知 | 内存存储 |
| **EMBED** | `EMBED` | 不支持 | go:embed 嵌入 |
| **MULTI** | `MULTI` | 透传子源 | 多源聚合（Fallback / FirstOK） |

### 扩展来源（独立子模块）

每个扩展来源是一个**独立 go.mod 子模块**，按需 import 即可启用，不影响主包体积：

| 类型 | 枚举值 | 子模块路径 | 底层库 | 热更新 |
| --- | --- | --- | --- | --- |
| **S3** | `S3` | `script_engine/source/s3` | AWS SDK v2 | ETag 比对 |
| **ETCD** | `ETCD` | `script_engine/source/etcd` | etcd v3 client | 原生 Watch |
| **CONSUL** | `CONSUL` | `script_engine/source/consul` | Consul API | ModifyIndex |
| **REDIS** | `REDIS` | `script_engine/source/redis` | go-redis v9 | 值比对 |
| **HTTP** | `HTTP` | `script_engine/source/http` | net/http | CRC32 校验 |
| **GIT** | `GIT` | `script_engine/source/git` | go-git v6 | commit hash |
| **DATABASE** | `DATABASE` | `script_engine/source/database` | database/sql | checksum |

> 扩展来源子模块目前是**占位实现**（工厂已注册，但返回 "not yet published" 错误）。
> 当 go-scripts 发布对应 source 子包后，取消 factory.go 中的 TODO 注释即可启用。

### 使用扩展来源

只需 import 对应子模块，工厂会在 `init()` 中自动注册：

```go
import (
    _ "github.com/tx7do/kratos-bootstrap/script_engine/source/etcd"   // 启用 etcd
    _ "github.com/tx7do/kratos-bootstrap/script_engine/source/s3"    // 启用 S3
    // 按需导入其他来源
)
```

对应的 YAML 配置：

```yaml
script:
  engine: LUA
  source:
    type: ETCD                              # 来源类型
    options:
      endpoints: ["localhost:2379"]        # etcd 特有配置
      prefix: "/scripts/"
      cache_ttl: "5m"                       # 自动包裹 CachedSource
```

### 自定义来源

如果需要接入 go-scripts 尚未覆盖的来源（如自研 RPC），可在宿主代码中直接注册：

```go
import (
    conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
    "github.com/tx7do/kratos-bootstrap/script_engine"
    "github.com/tx7do/go-scripts/source"
)

func init() {
    script_engine.MustRegisterSourceFactory(conf.Script_Source_S3,
        func(cfg *conf.Script_Source) (source.Reader, error) {
            // 自定义创建逻辑
            return myCustomSource.New(cfg.GetOptions().AsMap())
        })
}
```

### S3 + 本地缓存（MULTI 异构聚合）

```yaml
script:
  engine: LUA
  source:
    type: MULTI
    strategy: FALLBACK
    options:
      cache_ttl: "5m"           # S3 远程源自动包裹缓存层
      sources:                   # 异构子源列表
        - type: S3
          options:
            bucket: "my-scripts"
            region: "us-east-1"
            prefix: "lua/"
        - type: FILE
          paths:
            - "./scripts"       # 本地回退
  options:
    entry:
      value: "main.lua"
```

### Redis 来源 + 缓存

```yaml
script:
  engine: JAVASCRIPT
  source:
    type: REDIS
    options:
      addr: "localhost:6379"
      db: 0
      prefix: "script:"
      cache_ttl: "2m"          # 本地缓存 2 分钟
```

## 引擎注册

使用前需通过 `init()` 注册所需引擎（导入对应子包即可）：

```go
import (
    _ "github.com/tx7do/go-scripts/lua"        // 注册 Lua 引擎
    _ "github.com/tx7do/go-scripts/javascript" // 注册 JavaScript 引擎
    _ "github.com/tx7do/go-scripts/cel"        // 注册 CEL 引擎
    // ... 其他引擎按需导入
)
```
