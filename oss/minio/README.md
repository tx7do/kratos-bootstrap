# MinIO

`oss/minio` 提供基于 MinIO Go SDK 的 MinIO 客户端初始化与最小上传/下载封装。

## 功能

- 创建 MinIO SDK 客户端
- 提供轻量 `Storage` 封装
- 支持最小 `PutObject` / `GetObject` 操作
- 支持上传参数校验与通用 `io.Reader` 上传

## 配置示例

```yaml
oss:
  minio:
    endpoint: 127.0.0.1:9000
    access_key: minioadmin
    secret_key: minioadmin
    token: ""
    use_ssl: false
    upload_host: cdn-upload.example.com
    download_host: cdn.example.com
```

## 使用示例

```text
cfg := &conf.OSS{
    Minio: &conf.OSS_MinIO{
        Endpoint:  "127.0.0.1:9000",
        AccessKey: "minioadmin",
        SecretKey: "minioadmin",
        UseSsl:    false,
    },
}

storage := minio.NewStorage(cfg)
_, err := storage.PutObject(ctx, "example-bucket", "demo/hello.txt", strings.NewReader("hello world"), "text/plain")
if err != nil {
    panic(err)
}

obj, err := storage.GetObject(ctx, "example-bucket", "demo/hello.txt")
if err != nil {
    panic(err)
}
defer obj.Close()
```

## API

- `NewClient(cfg *conf.OSS) *minio.Client`
- `NewStorage(cfg *conf.OSS) *Storage`
- `(*Storage).SDK() *minio.Client`
- `(*Storage).PutObject(ctx, bucket, key, body, contentType)`
- `(*Storage).GetObject(ctx, bucket, key)`

## 与 `oss/s3` 的选型建议

`oss/minio` 和 `oss/s3` 都能操作 S3 协议对象存储，但定位不同。可以按下面方式选：

### 优先用 `oss/minio`（MinIO SDK）

- 你的主要目标是 MinIO（含自建 MinIO 集群）
- 你希望保持和 MinIO 生态的一致行为与调试体验
- 你当前项目已经大量使用 MinIO SDK，想减少迁移成本
- 你更关注在私有化、内网环境中稳定接入 MinIO

### 优先用 `oss/s3`（AWS SDK v2）

- 你主要接入 AWS S3，或同时接入多个 S3 兼容云厂商
- 你需要 AWS SDK v2 的统一生态能力（配置、鉴权、服务扩展）
- 你希望跨服务复用 AWS SDK 的基础设施与中间件能力
- 你希望使用 `oss/s3` 里已经整理好的多厂商 endpoint 配置经验

### 一句话建议

- **以 MinIO 为核心：选 `oss/minio`**
- **以 AWS S3 / 多云 S3 兼容接入为核心：选 `oss/s3`**

### 兼容性提醒

- 两个模块都支持基础上传下载，但底层 SDK 行为细节（重试、签名、默认 header）可能不同
- 如果你要在同一业务里混用两者，建议统一错误处理与重试策略，避免线上行为差异

## 迁移建议（`oss/minio` -> `oss/s3`）

下面是一份最小改造清单，适合把现有 MinIO 接入平滑迁移到 `oss/s3`。

### 1) 配置改造

- 将配置根从 `oss.minio` 迁移到 `oss.s3`
- 字段基本可复用：`endpoint`、`access_key`、`secret_key`、`token`、`use_ssl`
- 新增/重点字段：
  - `region`：建议显式填写（为空时 `oss/s3` 默认 `us-east-1`）
  - `bucket`：`oss/s3` 的 `Storage` 使用默认 bucket
  - `force_path_style`：MinIO/本地兼容服务一般建议 `true`

示例对照：

```yaml
# before: oss/minio
oss:
  minio:
    endpoint: 127.0.0.1:9000
    access_key: minioadmin
    secret_key: minioadmin
    token: ""
    use_ssl: false
```

```yaml
# after: oss/s3
oss:
  s3:
    endpoint: 127.0.0.1:9000
    region: us-east-1
    bucket: example-bucket
    access_key: minioadmin
    secret_key: minioadmin
    token: ""
    use_ssl: false
    force_path_style: true
```

### 2) 调用签名改造

- 初始化：`minio.NewStorage(cfg)` -> `s3.NewStorage(cfg)`
- 上传：
  - `minio`: `PutObject(ctx, bucket, key, body, contentType)`
  - `s3`: `PutObject(ctx, key, body, contentType)`（bucket 走配置默认值）
- 下载：
  - `minio`: `GetObject(ctx, bucket, key)`
  - `s3`: `GetObject(ctx, key)`

最小改造示例：

```text
// before (oss/minio)
storage := minio.NewStorage(cfg)
_, err := storage.PutObject(ctx, "example-bucket", "demo/hello.txt", strings.NewReader("hello"), "text/plain")
obj, err := storage.GetObject(ctx, "example-bucket", "demo/hello.txt")

// after (oss/s3)
storage := s3.NewStorage(cfg) // cfg.S3.Bucket = "example-bucket"
_, err := storage.PutObject(ctx, "demo/hello.txt", strings.NewReader("hello"), "text/plain")
obj, err := storage.GetObject(ctx, "demo/hello.txt")
```

### 3) 测试回归点

迁移后建议至少覆盖以下测试点：

- 配置加载：`region`、`bucket`、`force_path_style` 是否生效
- 上传下载：同一对象 `PutObject` 后 `GetObject` 内容一致
- endpoint 路由：本地/私有网关场景下 path-style 是否正确（尤其是 `force_path_style=true`）
- 证书与协议：`use_ssl=true/false` 切换是否符合部署环境
- 错误分支：空 key、空 bucket（s3 是配置缺失）、无效凭证、endpoint 不可达

### 4) 迁移策略建议

- 先在灰度环境“双写/双读对比”一小段时间，观察签名、重试、超时行为差异
- 对同一个 bucket/key 做一致性抽样校验（大小、ETag、内容）
- 最后切主流量到 `oss/s3`，保留短期回滚开关

## 注意事项

- `PutObject` 需要 `bucket` 与 `key` 非空
- `PutObject` 会对 `io.Reader` 做最小可读性处理，保证常见 reader 类型可上传
- `GetObject` 返回 `*minio.Object`，使用后请记得 `Close()`

