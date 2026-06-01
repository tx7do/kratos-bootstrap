# model 包说明

## 概述

`model` 包基于 [go-openai](https://github.com/sashabaranov/go-openai) 封装，提供 OpenAI 兼容的 LLM 客户端创建能力。

支持两种部署模式：

- **云端模型**（`CLOUD_MODEL`）：OpenAI、通义千问等兼容 OpenAI API 的云端服务
- **本地模型**（`LOCAL_MODEL`）：通过 Ollama 等工具部署的本地模型

## 特性

- 通过 Protobuf 配置（`conf.AI_Model`）统一初始化，无需硬编码连接参数
- 支持自定义 Base URL，兼容各种 OpenAI 兼容 API
- 支持自定义 HTTP 客户端和超时设置
- 支持多组织（Organization）配置

## API 概览

| 函数 | 说明 |
|------|------|
| `NewClient(cfg *conf.AI_Model, opts ...Option) (*openai.Client, error)` | 根据配置创建客户端 |
| `WithHTTPClient(httpClient *http.Client) Option` | 设置自定义 HTTP 客户端 |

返回的 `*openai.Client` 可直接使用 go-openai 的全部 API，包括 Chat Completion、Embedding、Function Calling 等。

## 使用示例

### 云端模型

```go
package example

import (
    openai "github.com/sashabaranov/go-openai"
    aiModel "github.com/tx7do/kratos-bootstrap/ai/model"

    conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func NewCloudClient() (*openai.Client, error) {
    cfg := &conf.AI_Model{
        Type:      conf.AI_Model_CLOUD_MODEL,
        ModelName: "gpt-4o",
        Cloud: &conf.AI_Model_Cloud{
            ApiKey:    "sk-xxx",
            BaseUrl:   "https://api.openai.com/v1",
        },
        TimeoutSeconds: 60,
    }
    return aiModel.NewClient(cfg)
}
```

### 本地模型（Ollama）

```go
func NewLocalClient() (*openai.Client, error) {
    cfg := &conf.AI_Model{
        Type:      conf.AI_Model_LOCAL_MODEL,
        ModelName: "llama3",
        Local: &conf.AI_Model_Local{
            Host: "localhost",
            Port: 11434,
        },
        TimeoutSeconds: 120,
    }
    return aiModel.NewClient(cfg)
}
```

### 调用 Chat Completion

```go
client, _ := aiModel.NewClient(cfg)

resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
    Model: cfg.GetModelName(),
    Messages: []openai.ChatCompletionMessage{
        {Role: openai.ChatMessageRoleSystem, Content: "你是一个翻译助手。"},
        {Role: openai.ChatMessageRoleUser, Content: "将以下内容翻译为英文：你好世界"},
    },
})
```
