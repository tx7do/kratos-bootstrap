# Script Engine

<p align="center">
  <a href="README.md">中文</a> · <a href="README_en.md">English</a> · <a href="README_ja.md">日本語</a>
</p>

[go-scripts](https://github.com/tx7do/go-scripts) ベースの多言語組み込みスクリプトエンジンラッパー。Kratos アプリケーションに統一的なスクリプト実行機能を提供します。

## 対応スクリプトエンジン（9種）

| エンジン | 列挙値 | ライブラリ | ユースケース |
| --- | --- | --- | --- |
| **Lua** | `LUA` | gopher-lua | ゲームスクリプト、設定ロジック、組み込み拡張 |
| **JavaScript** | `JAVASCRIPT` | goja | フロントエンド再利用、ルールエンジン、迅速なプロトタイピング |
| **Python** | `GPYTHON` | gpython | データ処理、運用スクリプト、アルゴリズム検証 |
| **Go** | `YAEGI` | Yaegi | 動的プラグイン、DevOpsツールチェーン |
| **WebAssembly** | `WAZERO` | wazero | 高パフォーマンスサンドボックス、cross-languageモジュール再利用 |
| **CEL** | `CEL` | cel-go | ポリシーエンジン、権限ルール、条件判定 |
| **Expr** | `EXPR` | expr-lang | ビジネス式、テンプレートエンジン、データフィルタリング |
| **Starlark** | `STARLARK` | starlark-go | ビルドツール、セキュアスクリプト、Bazelルール |
| **TCL** | `TCL` | modernc/tcl | レガシーシステム互換、ネットワーク機器スクリプト |

> Lua と JavaScript は**フルエンジン**（スクリプト読み込み、実行、関数登録、モジュール登録、ホットリロード対応）。
> その他のエンジンは**軽量エンジン**（コアライフサイクル + スクリプト実行のみ対応）。

## コア機能

- **Source モード**：`FileSource`、`MemSource`、`EmbedSource`、`MultiSource` による統一スクリプト読み込み
- **拡張ソース**：S3、etcd、Consul、Redis、HTTP、Git、Database など7種の拡張ソースをファクトリ登録機構で統合
- **キャッシュレイヤー**：リモートソース（S3/etcd等）は自動的に `CachedSource` でラップし、ネットワーク IO を削減
- **エンジンプール**：固定サイズプール `EnginePool` + 自動拡張プール `AutoGrowEnginePool`
- **ホットリロード**：`hot_reload: true` を設定するだけでスクリプト変更時の自動リロードを有効化
- **Kratos ロガー統合**：`KratosLogger` アダプターが Kratos 統合ログ体系にシームレスに接続

## 設定例

### YAML 設定

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

### go:embed によるスクリプト埋め込み

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
    // embed FS プロバイダーを登録
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

    // エンジン作成（embed から自動的にスクリプトを読み込み）
    eng, _ := script_engine.NewScriptEngine(ctx, cfg)
}
```

## 使用例

### 単一エンジンインスタンス

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

    // ホスト関数を登録
    _ = eng.RegisterFunction("greet", func(name string) string {
        return fmt.Sprintf("Hello, %s!", name)
    })

    // エントリスクリプトを実行
    result, err := eng.Execute(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result)
}
```

### エンジンプール（本番推奨）

```go
// 自動拡張プール（初期2、上限16）
pool, err := script_engine.NewAutoGrowEnginePool(ctx, cfg)
if err != nil {
    log.Fatal(err)
}
defer pool.Close()

// プール経由でスクリプト実行（自動 Acquire/Release）
result, err := pool.ExecuteString(ctx, "hello.js", `greet("world")`)
```

### Kratos ロガー統合

```go
import (
    "github.com/go-kratos/kratos/v2/log"
    "github.com/tx7do/kratos-bootstrap/script_engine"
)

// Kratos Logger を go-scripts に注入
script_engine.SetLogger(kratosLogger)
```

## ファイル構成

| ファイル / ディレクトリ | 説明 |
| --- | --- |
| `script_engine.go` | コアラッパー：NewScriptEngine / NewEnginePool / NewAutoGrowEnginePool |
| `source.go` | Source ファクトリ登録、ソース作成、パス解決、キャッシュラップ、ホットリロード |
| `logger.go` | Kratos Logger → go-scripts Logger アダプター |
| `utils.go` | エンジンタイプマッピング、機能検出 |
| `source/etcd/` | etcd ソースサブモジュール（独立 go.mod） |
| `source/s3/` | S3 ソースサブモジュール（独立 go.mod） |
| `source/redis/` | Redis ソースサブモジュール（独立 go.mod） |
| `source/consul/` | Consul ソースサブモジュール（独立 go.mod） |
| `source/http/` | HTTP ソースサブモジュール（独立 go.mod） |
| `source/git/` | Git ソースサブモジュール（独立 go.mod） |
| `source/database/` | SQLデータベースソースサブモジュール（独立 go.mod） |

## スクリプトソース（Source）

### 組み込みソース（go-scripts v0.0.6 対応済み）

| タイプ | 列挙値 | ホットリロード | 説明 |
| --- | --- | --- | --- |
| **FILE** | `FILE` | mtime ポーリング | ローカルファイルシステム |
| **MEMORY** | `MEMORY` | channel 通知 | メモリストレージ |
| **EMBED** | `EMBED` | 非対応 | go:embed 埋め込み |
| **MULTI** | `MULTI` | サブソース透過 | マルチソース集約（Fallback / FirstOK） |

### 拡張ソース（独立サブモジュール）

各拡張ソースは**独立した go.mod サブモジュール**で、必要に応じて import するだけで有効化でき、メインパッケージのサイズに影響しません：

| タイプ | 列挙値 | サブモジュールパス | ライブラリ | ホットリロード |
| --- | --- | --- | --- | --- |
| **S3** | `S3` | `script_engine/source/s3` | AWS SDK v2 | ETag 比較 |
| **ETCD** | `ETCD` | `script_engine/source/etcd` | etcd v3 client | ネイティブ Watch |
| **CONSUL** | `CONSUL` | `script_engine/source/consul` | Consul API | ModifyIndex |
| **REDIS** | `REDIS` | `script_engine/source/redis` | go-redis v9 | 値比較 |
| **HTTP** | `HTTP` | `script_engine/source/http` | net/http | CRC32 チェックサム |
| **GIT** | `GIT` | `script_engine/source/git` | go-git v6 | commit hash |
| **DATABASE** | `DATABASE` | `script_engine/source/database` | database/sql | checksum |

> 拡張ソースサブモジュールは現在**プレースホルダー実装**です（ファクトリは登録済みですが、"not yet published" エラーを返します）。
> go-scripts が対応する source サブパッケージを公開後、factory.go の TODO コメントを解除するだけで有効化できます。

### 拡張ソースの使用

対応するサブモジュールを import するだけで、ファクトリが `init()` で自動登録されます：

```go
import (
    _ "github.com/tx7do/kratos-bootstrap/script_engine/source/etcd"   // etcd を有効化
    _ "github.com/tx7do/kratos-bootstrap/script_engine/source/s3"    // S3 を有効化
    // その他のソースは必要に応じて import
)
```

対応する YAML 設定：

```yaml
script:
  engine: LUA
  source:
    type: ETCD                              # ソースタイプ
    options:
      endpoints: ["localhost:2379"]        # etcd 固有の設定
      prefix: "/scripts/"
      cache_ttl: "5m"                       # CachedSource で自動ラップ
```

### カスタムソース

go-scripts が未対応のソース（独自 RPC など）を統合する場合、アプリケーションコード内で直接登録できます：

```go
import (
    conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
    "github.com/tx7do/kratos-bootstrap/script_engine"
    "github.com/tx7do/go-scripts/source"
)

func init() {
    script_engine.MustRegisterSourceFactory(conf.Script_Source_S3,
        func(cfg *conf.Script_Source) (source.Reader, error) {
            // カスタム作成ロジック
            return myCustomSource.New(cfg.GetOptions().AsMap())
        })
}
```

### S3 + ローカルキャッシュ（MULTI 異種集約）

```yaml
script:
  engine: LUA
  source:
    type: MULTI
    strategy: FALLBACK
    options:
      cache_ttl: "5m"           # S3 リモートソースを自動的にキャッシュラップ
      sources:                   # 異種サブソースリスト
        - type: S3
          options:
            bucket: "my-scripts"
            region: "us-east-1"
            prefix: "lua/"
        - type: FILE
          paths:
            - "./scripts"       # ローカルフォールバック
  options:
    entry:
      value: "main.lua"
```

### Redis ソース + キャッシュ

```yaml
script:
  engine: JAVASCRIPT
  source:
    type: REDIS
    options:
      addr: "localhost:6379"
      db: 0
      prefix: "script:"
      cache_ttl: "2m"          # ローカルキャッシュ 2分間
```

## エンジン登録

使用前に `init()` で必要なエンジンを登録します（対応するサブパッケージを import するだけ）：

```go
import (
    _ "github.com/tx7do/go-scripts/lua"        // Lua エンジンを登録
    _ "github.com/tx7do/go-scripts/javascript" // JavaScript エンジンを登録
    _ "github.com/tx7do/go-scripts/cel"        // CEL エンジンを登録
    // ... その他のエンジンは必要に応じて import
)
```
