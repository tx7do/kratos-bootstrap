<p align="center">
  <h1 align="center">kratos-bootstrap</h1>
  <p align="center">
    go-kratos ベースのマイクロサービスアプリケーションブートストラップフレームワーク
  </p>
  <p align="center">
    <em>ワンストップのインフラストラクチャブートストラップ — マイクロサービス開発の繰り返し設定から解放</em>
  </p>
</p>

<p align="center">
  <a href="README.md">中文</a> · <a href="README_en.md">English</a> · <a href="README_ja.md">日本語</a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=Go" alt="Go Version" />
  <img src="https://img.shields.io/badge/Kratos-v2-00ADD8?style=flat-square" alt="Kratos" />
  <img src="https://img.shields.io/badge/License-MIT-green?style=flat-square" alt="License" />
  <img src="https://img.shields.io/badge/PRs-Welcome-brightgreen?style=flat-square" alt="PRs Welcome" />
</p>

---

## プロジェクトハイライト

- **統合ブートストラップエントリ**：設定読み込み、ロガー初期化、サービス登録、分散トレーシングを一括カプセル化 — 1回の呼び出しでマイクロサービスを起動
- **9種のデータベースエンジン**：ClickHouse、Doris、Elasticsearch、OpenSearch、MongoDB、InfluxDB、Cassandra、Ent、GORM など主要データベースをネイティブサポート
- **7種のサービスレジストリ**：Consul、Etcd、Nacos、Zookeeper、Eureka、Polaris、ServiceComb など主要サービスディスカバリコンポーネントを統合
- **6種のリモート設定センター**：Apollo、Consul KV、Etcd、Nacos、Polaris、Kubernetes ConfigMap のリモート設定管理をサポート
- **5種のトランスポート層**：Kafka、Asynq、MCP、MQTT、SSE などメッセージング・ストリーミングコンポーネントを内蔵
- **3種のAIフレームワーク**：go-openai、LangChainGo、ByteDance Eino を統合し、クラウドモデルとローカルモデル（Ollama）のシームレスな切り替えを実現
- **多言語スクリプトエンジン**：JavaScript（ES6+）と Lua スクリプトによるビジネスロジックの動的拡張をサポート
- **フルチェーンオブザーバビリティ**：OpenTelemetry ベースの分散トレーシング、Jaeger など多様な Exporter に対応
- **6種のログフレームワーク**：Zap、Zerolog、Logrus、Fluentd、Alibaba Cloud SLS、Tencent Cloud CLS のログ出力をサポート
- **設定契約ファースト**：全設定を Protobuf で定義 — 型安全、強制制約、自動生成

---

## 技術スタック

| 層 | 技術 | 説明 |
| --- | --- | --- |
| 言語 | Go 1.23+ | 高パフォーマンスコンパイル言語 |
| フレームワーク | go-kratos v2 | Bilibili発のオープンソースマイクロサービスフレームワーク |
| 設定定義 | Protobuf + buf.build | 契約ファースト、型安全な設定管理 |
| トレーシング | OpenTelemetry | 分散オブザーバビリティ標準 |
| オブジェクトストレージ | MinIO / S3 | S3互換オブジェクトストレージ |
| キャッシュ | Redis | go-redis v9、トレーシング・メトリクス採取対応 |
| AIフレームワーク | go-openai / LangChainGo / Eino | マルチフレームワーク LLM 統合（クラウド＆ローカル） |
| メッセージキュー | Kafka | 高スループットイベントストリーミング |
| 非同期タスク | Asynq | Redis ベースの非同期タスクキュー |
| AIツールプロトコル | MCP (Model Context Protocol) | AIエージェントツール呼び出し標準プロトコル |
| リアルタイムプッシュ | SSE | サーバーセントイベント |
| IoT | MQTT | 軽量メッセージングプロトコル |
| スクリプトエンジン | goja (JS) / gopher-lua | 多言語動的スクリプト実行 |

---

## プロジェクト構成

```
kratos-bootstrap/
├── api/                                # Protobuf API定義と生成コード
│   ├── protos/conf/v1/                 # .protoソースファイル（設定構造定義）
│   └── gen/go/conf/v1/                 # buf生成Goコード
├── bootstrap/                          # アプリケーションブートストラップコア
│   ├── bootstrap.go                    # ブートストラップエントリ（設定/ロガー/レジストリ/トレーサー）
│   ├── cli.go                          # CLIフレームワーク（Cobra）
│   ├── context.go                      # ブートストラップコンテキスト
│   └── daemon.go                       # デーモンプロセスサポート
├── ai/                                 # AI大規模モデル統合
│   ├── model/                          # go-openaiネイティブクライアント
│   ├── langchaingo/                    # LangChainGoフレームワークラッパー
│   └── eino/                           # ByteDance Einoフレームワークラッパー
├── config/                             # リモート設定センター
│   ├── apollo/                         # Apollo
│   ├── consul/                         # Consul KV
│   ├── etcd/                           # Etcd
│   ├── nacos/                          # Nacos
│   ├── polaris/                        # Polaris
│   └── kubernetes/                     # Kubernetes ConfigMap
├── registry/                           # サービス登録・ディスカバリ
│   ├── consul/                         # Consul
│   ├── etcd/                           # Etcd
│   ├── nacos/                          # Nacos
│   ├── zookeeper/                      # Zookeeper
│   ├── eureka/                         # Eureka
│   ├── polaris/                        # Polaris
│   ├── servicecomb/                    # ServiceComb
│   └── kubernetes/                     # Kubernetes
├── database/                           # データベースクライアント
│   ├── clickhouse/                     # ClickHouse（OLAP）
│   ├── doris/                          # Apache Doris（OLAP）
│   ├── elasticsearch/                  # Elasticsearch
│   ├── opensearch/                     # OpenSearch
│   ├── mongodb/                        # MongoDB
│   ├── influxdb/                       # InfluxDB（時系列DB）
│   ├── cassandra/                      # Cassandra
│   ├── ent/                            # Ent ORM
│   └── gorm/                           # GORM ORM
├── transport/                          # トランスポート層
│   ├── kafka/                          # Kafka
│   ├── asynq/                          # Asynq非同期タスク
│   ├── mcp/                            # MCP AIツールプロトコル
│   ├── mqtt/                           # MQTT IoTプロトコル
│   └── sse/                            # SSEサーバープッシュ
├── logger/                             # ログフレームワーク統合
│   ├── zap/                            # Zap
│   ├── zerolog/                        # Zerolog
│   ├── logrus/                         # Logrus
│   ├── fluent/                         # Fluentd
│   ├── aliyun/                         # Alibaba Cloud SLS
│   └── tencent/                        # Tencent Cloud CLS
├── cache/                              # キャッシュ
│   └── redis/                          # Redisクライアント
├── tracer/                             # 分散トレーシング
│   └── exporter.go                     # OpenTelemetry Exporterファクトリ
├── rpc/                                # RPC通信
│   ├── grpc.go                         # gRPCクライアント/サーバー
│   ├── rest.go                         # REST (HTTP) クライアント/サーバー
│   └── middleware/                     # RPCミドルウェア
│       ├── validate/                   # Protobufバリデーション
│       └── requestid/                  # リクエストIDミドルウェア
├── oss/                                # オブジェクトストレージ
│   ├── minio/                          # MinIOクライアント
│   └── s3/                             # S3互換クライアント
└── script_engine/                      # スクリプトエンジン
    └── script_engine.go                # JavaScript / Luaスクリプトエンジン
```

---

## コア機能

### アプリケーションブートストラップ

| 機能 | 説明 |
| --- | --- |
| 統合エントリ | `Bootstrap`関数一つで設定読み込み、ロガー初期化、サービス登録、トレーシングを完結 |
| CLIフレームワーク | CobraベースのCLI、サブコマンドカスタマイズとフラグ注入に対応 |
| デーモンモード | ネイティブデーモンサポート、バックグラウンド実行とPID管理 |
| アプリメタデータ | アプリ名、バージョン、インスタンスID、ネームスペースの一元管理 |
| グレースフルシャットダウン | シグナルキャプチャと安全なサービス終了を内蔵 |

### サービス登録・ディスカバリ

| レジストリ | 説明 |
| --- | --- |
| Consul | HashiCorp サービスディスカバリ & KVストア |
| Etcd | 高可用分散KV（Kubernetesバックボーン） |
| Nacos | Alibaba Cloud マイクロサービス登録 & 設定センター |
| Zookeeper | Apache 分散調整サービス |
| Eureka | Netflix サービスディスカバリ |
| Polaris | Tencent Cloud Polaris サービスガバナンス |
| ServiceComb | Huawei Cloud マイクロサービスエンジン |
| Kubernetes | ネイティブ Kubernetes サービスディスカバリ |

### リモート設定センター

| 設定センター | 説明 |
| --- | --- |
| Apollo | Ctrip 分散設定管理 |
| Consul | Consul KV 設定ストレージ |
| Etcd | Etcd 分散設定 |
| Nacos | Nacos 設定管理 |
| Polaris | Polaris 設定管理 |
| Kubernetes | Kubernetes ConfigMap |

### データベースサポート

| データベース | タイプ | 説明 |
| --- | --- | --- |
| ClickHouse | OLAP | カラムナストレージ、極致の分析パフォーマンス |
| Apache Doris | OLAP | 高パフォーマンスリアルタイム分析エンジン |
| Elasticsearch | 検索エンジン | フルテキスト検索 & ログ分析 |
| OpenSearch | 検索エンジン | Elasticsearch オープンソースフォーク |
| MongoDB | ドキュメントDB | フレキシブルなドキュメントモデル |
| InfluxDB | 時系列DB | 時系列データの保存 & クエリ |
| Cassandra | ワイドカラムDB | 高可用分散ストレージ |
| Ent | ORM | Go エンティティフレームワーク |
| GORM | ORM | 最も人気のある Go ORM |

### AI / LLM統合

| モジュール | フレームワーク | 説明 |
| --- | --- | --- |
| model | go-openai | OpenAI互換APIネイティブクライアント |
| langchaingo | LangChainGo | Chain / Agent / Embedding / VectorStore / Memory |
| eino | ByteDance Eino | Chain / Tool / Prompt / Compose |

全AIモジュール対応：
- **クラウドモデル**：OpenAI、Qwen、DeepSeek、その他 OpenAI API 互換モデルサービス
- **ローカルモデル**：Ollama ローカルデプロイ、クラウド依存ゼロ

### トランスポート層

| コンポーネント | 説明 |
| --- | --- |
| Kafka | 高スループットメッセージキュー、SASL/SCRAM認証 & TLS暗号化対応 |
| Asynq | Redisベース非同期タスクキュー（スタンドアロン/クラスタ/センチネル） |
| MCP | Model Context Protocol — AIエージェントツール呼び出し（HTTP/SSE/Stdio） |
| MQTT | 軽量IoTメッセージングプロトコル |
| SSE | Server-Sent Events リアルタイムプッシュ |

### ログフレームワーク

| フレームワーク | 説明 |
| --- | --- |
| Zap | Uber の高性能構造化ロガー |
| Zerolog | ゼロアロケーション JSONロガー |
| Logrus | 構造化ロガー |
| Fluent | Fluentd ログ収集 |
| Alibaba Cloud SLS | Alibaba Cloud ログサービス |
| Tencent Cloud CLS | Tencent Cloud ログサービス |

### 分散トレーシング

OpenTelemetry標準ベース：
- 設定可能なサンプリングレート（TraceIDRatioBased サンプリング）
- W3C TraceContext & Baggage 伝搬
- プラグイン可能な Exporter ファクトリパターン
- サービスメタデータの自動注入（サービス名、バージョン、インスタンスID、環境）

### RPC通信

| 機能 | 説明 |
| --- | --- |
| gRPC サーバー/クライアント | ネイティブ gRPC、TLS / タイムアウト / ミドルウェアチェーン内蔵 |
| REST (HTTP) サーバー | HTTP/RESTful APIサービス、CORS & pprof デバッグ対応 |
| ミドルウェアチェーン | Recovery / Tracing / Validate / RateLimit / Metadata |
| レートリミット | BBR（Google BBRアルゴリズム）アダプティブレートリミット |
| ホワイトリスト | メソッド名ベースのミドルウェアホワイトリスト機構 |

### スクリプトエンジン

| 言語 | エンジン | 説明 |
| --- | --- | --- |
| JavaScript | goja | ES6+サポート、ルールエンジンや動的計算に最適 |
| Lua | gopher-lua | 軽量組み込みスクリプト言語 |

スクリプトプール管理（固定サイズ / 自動スケーリング）、プリロードスクリプト、エントリスクリプトに対応。

### オブジェクトストレージ

| ストレージ | 説明 |
| --- | --- |
| MinIO | セルフホスト S3互換オブジェクトストレージ |
| S3 | Amazon S3 / 互換ストレージ |

---

## 設計原則

- **設定駆動**：全コンポーネントを Protobuf 定義の設定で初期化 — 型安全、ハードコーディング排除
- **非侵入ラッピング**：ラッパー層は設定変換とインスタンス生成のみを行い、フレームワークネイティブAPIを維持
- **ファクトリパターン**：レジストリ、ロガー、トレーサーなどのコンポーネントにファクトリ登録パターンを採用し、オンデマンド読み込みと疎結合を実現
- **モジュラー設計**：各機能モジュールは独立した Go Module — 必要なものだけをインポート、不要な依存関係を排除
- **防御的プログラミング**：包括的なnilチェックとエラーハンドリングにより、コンポーネント欠落時にメインフローを妨げない設計

---

## ユースケース

- マイクロサービスアーキテクチャにおける標準化されたアプリケーションブートストラップとインフラ管理
- 複数データベースとミドルウェアの統合が必要な複雑なビジネスシステム
- マルチLLM・エージェントフレームワーク統合を必要とするAIアプリケーション開発
- MQTTメッセージングとエッジコンピューティングを必要とするIoTシナリオ
- Kafka / Asynq 非同期タスクスケジューリングによるリアルタイムデータ処理
- 動的スクリプトによるビジネスルール拡張（ルールエンジン、動的計算など）

---

## ライセンス

本プロジェクトは MIT ライセンスの下で公開されています。詳細は [LICENSE](LICENSE) ファイルをご参照ください。
