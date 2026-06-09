# Script Engine

<p align="center">
  <a href="README.md">中文</a> · <a href="README_en.md">English</a> · <a href="README_ja.md">日本語</a>
</p>

A multi-language embedded scripting engine wrapper built on [go-scripts](https://github.com/tx7do/go-scripts), providing unified script execution capabilities for Kratos applications.

## Supported Script Engines (9 types)

| Engine | Enum | Library | Use Cases |
| --- | --- | --- | --- |
| **Lua** | `LUA` | gopher-lua | Game scripting, config logic, embedded extensions |
| **JavaScript** | `JAVASCRIPT` | goja | Frontend reuse, rule engines, rapid prototyping |
| **Python** | `GPYTHON` | gpython | Data processing, ops scripts, algorithm validation |
| **Go** | `YAEGI` | Yaegi | Dynamic plugins, DevOps toolchain |
| **WebAssembly** | `WAZERO` | wazero | High-performance sandbox, cross-language module reuse |
| **CEL** | `CEL` | cel-go | Policy engine, permission rules, conditional logic |
| **Expr** | `EXPR` | expr-lang | Business expressions, template engine, data filtering |
| **Starlark** | `STARLARK` | starlark-go | Build tools, secure scripting, Bazel rules |
| **TCL** | `TCL` | modernc/tcl | Legacy system compatibility, network device scripting |

> Lua and JavaScript are **full engines** (support script loading, execution, function registration, module registration, hot reload).
> Other engines are **lightweight engines** (only support core lifecycle + script execution).

## Core Features

- **Source Mode**: Unified script loading via `FileSource`, `MemSource`, `EmbedSource`, `MultiSource`
- **Extended Sources**: 7 extended sources (S3, etcd, Consul, Redis, HTTP, Git, Database) via factory registration
- **Caching Layer**: Remote sources (S3/etcd etc.) automatically wrapped with `CachedSource` to reduce network IO
- **Engine Pool**: Fixed-size pool `EnginePool` + auto-scaling pool `AutoGrowEnginePool`
- **Hot Reload**: Enable automatic script reloading on file changes with `hot_reload: true`
- **Kratos Logger Integration**: `KratosLogger` adapter seamlessly bridges Kratos unified logging

## Configuration Examples

### YAML Configuration

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

### Using go:embed for Embedded Scripts

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
    // Register embed FS provider
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

    // Create engine (auto-loads scripts from embed)
    eng, _ := script_engine.NewScriptEngine(ctx, cfg)
}
```

## Usage Examples

### Single Engine Instance

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

    // Register host function
    _ = eng.RegisterFunction("greet", func(name string) string {
        return fmt.Sprintf("Hello, %s!", name)
    })

    // Execute entry script
    result, err := eng.Execute(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result)
}
```

### Engine Pool (Recommended for Production)

```go
// Auto-scaling pool (initial 2, max 16)
pool, err := script_engine.NewAutoGrowEnginePool(ctx, cfg)
if err != nil {
    log.Fatal(err)
}
defer pool.Close()

// Execute via pool (auto Acquire/Release)
result, err := pool.ExecuteString(ctx, "hello.js", `greet("world")`)
```

### Kratos Logger Integration

```go
import (
    "github.com/go-kratos/kratos/v2/log"
    "github.com/tx7do/kratos-bootstrap/script_engine"
)

// Inject Kratos Logger into go-scripts
script_engine.SetLogger(kratosLogger)
```

## File Structure

| File / Directory | Description |
| --- | --- |
| `script_engine.go` | Core wrapper: NewScriptEngine / NewEnginePool / NewAutoGrowEnginePool |
| `source.go` | Source factory registry, source creation, path resolution, cache wrapping, hot reload |
| `logger.go` | Kratos Logger → go-scripts Logger adapter |
| `utils.go` | Engine type mapping, capability detection |
| `source/etcd/` | etcd source sub-module (independent go.mod) |
| `source/s3/` | S3 source sub-module (independent go.mod) |
| `source/redis/` | Redis source sub-module (independent go.mod) |
| `source/consul/` | Consul source sub-module (independent go.mod) |
| `source/http/` | HTTP source sub-module (independent go.mod) |
| `source/git/` | Git source sub-module (independent go.mod) |
| `source/database/` | SQL database source sub-module (independent go.mod) |

## Script Sources

### Built-in Sources (supported in go-scripts v0.0.6)

| Type | Enum | Hot Reload | Description |
| --- | --- | --- | --- |
| **FILE** | `FILE` | mtime polling | Local filesystem |
| **MEMORY** | `MEMORY` | channel notification | In-memory storage |
| **EMBED** | `EMBED` | Not supported | go:embed embedded |
| **MULTI** | `MULTI` | Pass-through sub-sources | Multi-source aggregation (Fallback / FirstOK) |

### Extended Sources (Independent Sub-modules)

Each extended source is an **independent go.mod sub-module** — import on demand, no impact on main package size:

| Type | Enum | Sub-module Path | Library | Hot Reload |
| --- | --- | --- | --- | --- |
| **S3** | `S3` | `script_engine/source/s3` | AWS SDK v2 | ETag comparison |
| **ETCD** | `ETCD` | `script_engine/source/etcd` | etcd v3 client | Native Watch |
| **CONSUL** | `CONSUL` | `script_engine/source/consul` | Consul API | ModifyIndex |
| **REDIS** | `REDIS` | `script_engine/source/redis` | go-redis v9 | Value comparison |
| **HTTP** | `HTTP` | `script_engine/source/http` | net/http | CRC32 checksum |
| **GIT** | `GIT` | `script_engine/source/git` | go-git v6 | commit hash |
| **DATABASE** | `DATABASE` | `script_engine/source/database` | database/sql | checksum |

> Extended source sub-modules are currently **placeholder implementations** (factory is registered, but returns a "not yet published" error).
> Once go-scripts publishes the corresponding source sub-package, uncomment the TODO in factory.go to enable.

### Using Extended Sources

Simply import the corresponding sub-module — the factory auto-registers in `init()`:

```go
import (
    _ "github.com/tx7do/kratos-bootstrap/script_engine/source/etcd"   // Enable etcd
    _ "github.com/tx7do/kratos-bootstrap/script_engine/source/s3"    // Enable S3
    // Import other sources as needed
)
```

Corresponding YAML configuration:

```yaml
script:
  engine: LUA
  source:
    type: ETCD                              # Source type
    options:
      endpoints: ["localhost:2379"]        # etcd-specific config
      prefix: "/scripts/"
      cache_ttl: "5m"                       # Auto-wrap with CachedSource
```

### Custom Sources

If you need to integrate sources not yet covered by go-scripts (e.g., custom RPC), register directly in your application code:

```go
import (
    conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
    "github.com/tx7do/kratos-bootstrap/script_engine"
    "github.com/tx7do/go-scripts/source"
)

func init() {
    script_engine.MustRegisterSourceFactory(conf.Script_Source_S3,
        func(cfg *conf.Script_Source) (source.Reader, error) {
            // Custom creation logic
            return myCustomSource.New(cfg.GetOptions().AsMap())
        })
}
```

### S3 + Local Cache (MULTI Heterogeneous Aggregation)

```yaml
script:
  engine: LUA
  source:
    type: MULTI
    strategy: FALLBACK
    options:
      cache_ttl: "5m"           # S3 remote source auto-wrapped with cache
      sources:                   # Heterogeneous sub-source list
        - type: S3
          options:
            bucket: "my-scripts"
            region: "us-east-1"
            prefix: "lua/"
        - type: FILE
          paths:
            - "./scripts"       # Local fallback
  options:
    entry:
      value: "main.lua"
```

### Redis Source + Cache

```yaml
script:
  engine: JAVASCRIPT
  source:
    type: REDIS
    options:
      addr: "localhost:6379"
      db: 0
      prefix: "script:"
      cache_ttl: "2m"          # Local cache for 2 minutes
```

## Engine Registration

Register desired engines via `init()` before use (just import the corresponding sub-package):

```go
import (
    _ "github.com/tx7do/go-scripts/lua"        // Register Lua engine
    _ "github.com/tx7do/go-scripts/javascript" // Register JavaScript engine
    _ "github.com/tx7do/go-scripts/cel"        // Register CEL engine
    // ... import other engines as needed
)
```
