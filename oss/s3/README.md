# S3

`oss/s3` 提供基于 AWS SDK for Go v2 的 S3 客户端初始化能力。

## 功能

- 支持 AWS S3
- 支持兼容 S3 协议的对象存储（通过 `endpoint`）
- 支持静态 `access_key` / `secret_key` / `token`
- 支持 `region`
- 支持 `force_path_style`
- 支持根据 `use_ssl` 自动补齐 `endpoint` 协议头

## 配置示例

```yaml
oss:
  s3:
    endpoint: s3.ap-southeast-1.amazonaws.com
    region: ap-southeast-1
    bucket: example-bucket
    access_key: your-access-key
    secret_key: your-secret-key
    token: ""
    use_ssl: true
    force_path_style: false
    upload_host: cdn-upload.example.com
    download_host: cdn.example.com
```

如果接的是 MinIO、Ceph、LocalStack、RustFS 等兼容 S3 的服务，可配置自定义 `endpoint`，例如：

```yaml
oss:
  s3:
    endpoint: 127.0.0.1:9000
    region: us-east-1
    bucket: example-bucket
    access_key: minioadmin
    secret_key: minioadmin
    use_ssl: false
    force_path_style: true
```

## 兼容服务说明

除了 AWS S3，本模块也适用于越来越多实现了 S3 API 的对象存储服务。常见接入建议如下：

| 服务 | 是否兼容 S3 | 配置建议 |
| --- | --- | --- |
| AWS S3 | 是 | `endpoint` 可留空或使用官方域名；`region` 建议按真实地域填写；`force_path_style` 一般为 `false` |
| MinIO | 是 | 通常需要配置自定义 `endpoint`；多数场景建议 `force_path_style: true` |
| Ceph RGW | 是 | 建议显式配置 `endpoint`；部分部署更适合 `force_path_style: true` |
| LocalStack | 是 | 本地开发常用，通常配置 `endpoint: 127.0.0.1:4566`、`use_ssl: false`、`force_path_style: true` |
| RustFS | 是 | Rust 实现的 S3 兼容对象存储，建议显式配置 `endpoint`，通常推荐 `force_path_style: true` |
| Garage | 是 | 社区常见轻量对象存储，建议显式配置 `endpoint`，通常推荐 `force_path_style: true` |
| SeaweedFS S3 | 是 | 启用 S3 Gateway 后可接入，通常建议 `force_path_style: true` |
| Wasabi | 是 | 使用 Wasabi 提供的区域 endpoint；通常保持 `force_path_style: false`，按 Wasabi 控制台地域填写 `region` |
| Backblaze B2 S3 | 是 | 使用 B2 提供的 S3-compatible endpoint；一般建议按官方 endpoint 填写，`region` 依服务信息填写 |
| Tencent COS | 兼容 | 使用 COS 的 S3 API endpoint；多数情况下更适合保持默认主机风格访问，`force_path_style` 通常为 `false` |
| Aliyun OSS（S3兼容访问） | 视部署而定 | 若你的 OSS 接入层或网关提供 S3 兼容 endpoint，可按该 endpoint 配置；是否启用 path-style 取决于该兼容层实现 |
| Cloudflare R2 | 基本兼容 | 一般使用厂商提供的专属 endpoint；通常不需要 path-style，`region` 可按厂商要求填写 |

## 各厂商 endpoint 示例

下面是一些常见服务的 `endpoint` 填写方式示例。实际使用时请以服务商控制台、官方文档或你的私有部署地址为准。

### AWS S3

```yaml
oss:
  s3:
    endpoint: s3.ap-southeast-1.amazonaws.com
    region: ap-southeast-1
    force_path_style: false
```

也可以不显式填写 `endpoint`，仅通过 `region` 交给 AWS SDK 自动解析。

### MinIO

```yaml
oss:
  s3:
    endpoint: 127.0.0.1:9000
    region: us-east-1
    use_ssl: false
    force_path_style: true
```

### LocalStack

```yaml
oss:
  s3:
    endpoint: 127.0.0.1:4566
    region: us-east-1
    use_ssl: false
    force_path_style: true
```

### RustFS

```yaml
oss:
  s3:
    endpoint: rustfs.example.local:9000
    region: us-east-1
    use_ssl: false
    force_path_style: true
```

### SeaweedFS S3 Gateway

```yaml
oss:
  s3:
    endpoint: seaweedfs.example.local:8333
    region: us-east-1
    use_ssl: false
    force_path_style: true
```

### Garage

```yaml
oss:
  s3:
    endpoint: garage.example.local:3900
    region: us-east-1
    use_ssl: false
    force_path_style: true
```

### Wasabi

```yaml
oss:
  s3:
    endpoint: s3.ap-northeast-1.wasabisys.com
    region: ap-northeast-1
    use_ssl: true
    force_path_style: false
```

### Backblaze B2 S3

```yaml
oss:
  s3:
    endpoint: s3.us-west-004.backblazeb2.com
    region: us-west-004
    use_ssl: true
    force_path_style: false
```

### Tencent COS

```yaml
oss:
  s3:
    endpoint: cos.ap-guangzhou.myqcloud.com
    region: ap-guangzhou
    use_ssl: true
    force_path_style: false
```

### Cloudflare R2

```yaml
oss:
  s3:
    endpoint: <account-id>.r2.cloudflarestorage.com
    region: auto
    use_ssl: true
    force_path_style: false
```

### Ceph RGW / 私有 S3 网关

```yaml
oss:
  s3:
    endpoint: ceph-rgw.example.local
    region: us-east-1
    use_ssl: true
    force_path_style: true
```

### RustFS 配置示例

`RustFS` 属于 S3 兼容对象存储，接入方式与 MinIO 类似。一般建议：

- 配置服务暴露出来的 `endpoint`
- 如果服务没有做虚拟主机风格 bucket 路由，建议开启 `force_path_style: true`
- 如果是自签名证书或内网 HTTPS，建议结合上层网络/证书配置一起验证

示例：

```yaml
oss:
  s3:
    endpoint: rustfs.example.local:9000
    region: us-east-1
    bucket: example-bucket
    access_key: your-access-key
    secret_key: your-secret-key
    token: ""
    use_ssl: false
    force_path_style: true
```

如果你的 RustFS 部署支持标准域名风格 bucket 访问，也可以将 `force_path_style` 调整为 `false`。

## 什么时候要开启 `force_path_style`

`force_path_style` 会影响 bucket 在请求 URL 中的组织方式：

- `false`：更接近 AWS S3 默认风格，bucket 会出现在主机名里，例如 `https://my-bucket.s3.amazonaws.com/object.txt`
- `true`：使用路径风格，请求形式更像 `https://s3.example.com/my-bucket/object.txt`

通常可以按下面经验判断：

### 建议开启 `force_path_style: true`

- 你接的是 MinIO、RustFS、SeaweedFS S3、Garage、LocalStack 这类自建或本地 S3 兼容服务
- 服务只有一个统一域名或 IP:Port，例如 `127.0.0.1:9000`
- 你的证书并没有覆盖 `bucket.endpoint` 这种泛域名访问
- 反向代理、网关、Ingress 没有配置 bucket 子域名路由
- 兼容服务文档明确要求使用 path-style

### 通常保持 `force_path_style: false`

- 你接的是 AWS S3 官方 endpoint
- 你使用 Wasabi、Cloudflare R2、Tencent COS 这类云厂商提供的标准兼容 endpoint，且文档推荐主机风格访问
- 你的域名、证书、网关已经支持 `bucket.example.com` 这类虚拟主机风格 bucket 路由

### 不确定时怎么选

如果你不确定，建议优先按服务官方文档配置；若遇到以下问题，可以优先尝试打开 `force_path_style`：

- 上传或下载时出现签名不匹配
- 请求被路由到错误的 bucket
- HTTPS 证书域名不匹配
- 使用内网域名、IP、端口直连时访问失败

一个简单经验是：

- **云厂商官方 S3 服务优先尝试 `false`**
- **自建/本地/网关型 S3 兼容服务优先尝试 `true`**

## 使用示例

```text
package main

import (
    conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
    kss3 "github.com/tx7do/kratos-bootstrap/oss/s3"
)

func main() {
    client := kss3.NewClient(&conf.OSS{
        S3: &conf.OSS_S3{
            Endpoint:       "s3.ap-southeast-1.amazonaws.com",
            Region:         "ap-southeast-1",
            Bucket:         "example-bucket",
            AccessKey:      "your-access-key",
            SecretKey:      "your-secret-key",
            UseSsl:         true,
            ForcePathStyle: false,
        },
    })

    _ = client
}
```

## 最小上传/下载封装

`oss/s3` 还提供了一个轻量封装 `Storage`，会自动复用配置中的默认 `bucket`，适合直接做基础对象上传下载。

### 初始化 `Storage`

```text
storage := kss3.NewStorage(&conf.OSS{
    S3: &conf.OSS_S3{
        Endpoint:       "127.0.0.1:9000",
        Region:         "us-east-1",
        Bucket:         "example-bucket",
        AccessKey:      "minioadmin",
        SecretKey:      "minioadmin",
        UseSsl:         false,
        ForcePathStyle: true,
    },
})
```

### 上传对象

```text
_, err := storage.PutObject(ctx, "demo/hello.txt", strings.NewReader("hello world"), "text/plain")
if err != nil {
    panic(err)
}
```

### 下载对象

```text
resp, err := storage.GetObject(ctx, "demo/hello.txt")
if err != nil {
    panic(err)
}
defer resp.Body.Close()
```

### 封装说明

- `NewStorage`：创建带默认 bucket 的轻量封装
- `Storage.PutObject`：上传对象到配置中的默认 bucket
- `Storage.GetObject`：从配置中的默认 bucket 下载对象
- `Storage.SDK()`：拿到底层 `*s3.Client`，便于继续调用 AWS SDK 原生能力

## 说明

- 当 `region` 为空时，默认使用 `us-east-1`
- 当 `endpoint` 未携带协议头时，会根据 `use_ssl` 自动补齐为 `http://` 或 `https://`
- `bucket` 当前用于业务侧默认桶配置，客户端初始化时不会自动发起桶级操作

