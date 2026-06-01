# langchaingo 包说明

## 概述

`langchaingo` 包基于 [LangChainGo](https://github.com/tmc/langchaingo) 框架封装，提供完整的 LLM 应用开发能力。

LangChainGo 是 Python LangChain 的 Go 移植版本，本包对其进行了统一封装，覆盖以下功能模块：

- **Model**：LLM 客户端创建（云端 / 本地模型）
- **Agent**：ReAct、Conversational、OpenAI Functions 等多种 Agent
- **Chain**：LLM Chain、对话链、摘要链、顺序链等
- **Memory**：对话记忆管理（缓冲 / 窗口 / Token 限制）
- **Embedding**：文本向量化
- **VectorStore**：向量存储与相似度搜索

## 子模块说明

| 文件 | 说明 |
|------|------|
| `client.go` | LLM 客户端创建（OpenAI 云端 / Ollama 本地） |
| `agent.go` | Agent 创建（OneShot / Conversational / OpenAI Functions） |
| `chain.go` | Chain 编排（LLM Chain / 对话链 / 摘要链 / 顺序链） |
| `memory.go` | 对话记忆管理（Buffer / Window / Token Buffer） |
| `embedding.go` | 文本向量化 |
| `vectorstore.go` | 向量存储与检索辅助方法 |
| `options.go` | 可选配置项（OpenAI / Ollama / HTTP 客户端） |

## API 概览

### Model

| 函数 | 说明 |
|------|------|
| `NewModel(cfg, opts...) (llms.Model, error)` | 根据配置创建 LLM 客户端 |
| `WithOpenaiOptions(opts...) Option` | 追加 OpenAI 原生选项 |
| `WithOllamaOptions(opts...) Option` | 追加 Ollama 原生选项 |
| `WithHTTPClient(client) Option` | 设置自定义 HTTP 客户端 |

### Agent

| 函数 | 说明 |
|------|------|
| `NewOneShotAgent(llm, tools, opts...) *OneShotZeroAgent` | 创建 ReAct 单次 Agent |
| `NewConversationalAgent(llm, tools, opts...) *ConversationalAgent` | 创建对话式 Agent |
| `NewOpenAIFunctionsAgent(llm, tools, opts...) *OpenAIFunctionsAgent` | 创建 Function Calling Agent |
| `NewExecutor(agent, opts...) *Executor` | 创建 Agent 执行器 |
| `NewOneShotExecutor(llm, tools, opts...) *Executor` | 便捷：创建 OneShot Agent + Executor |
| `NewConversationalExecutor(llm, tools, opts...) *Executor` | 便捷：创建 Conversational Agent + Executor |
| `NewOpenAIFunctionsExecutor(llm, tools, opts...) *Executor` | 便捷：创建 OpenAI Functions Agent + Executor |

### Chain

| 函数 | 说明 |
|------|------|
| `NewLLMChain(llm, prompt, opts...) *LLMChain` | 创建 LLM 链 |
| `NewConversationChain(llm, mem) LLMChain` | 创建对话链（内置记忆） |
| `NewSequentialChain(c, inputKeys, outputKeys, opts...)` | 创建顺序链 |
| `NewSimpleSequentialChain(c)` | 创建简单顺序链 |
| `NewStuffDocumentsChain(llmChain)` | 创建文档填充链 |
| `LoadStuffSummarization(llm)` | 加载摘要链 |
| `LoadRefineSummarization(llm)` | 加载精炼摘要链 |
| `LoadMapReduceSummarization(llm)` | 加载 MapReduce 摘要链 |

### Memory

| 函数 | 说明 |
|------|------|
| `NewChatMessageHistory(opts...) *ChatMessageHistory` | 创建聊天消息历史 |
| `NewConversationBuffer(opts...) *ConversationBuffer` | 完整对话缓冲 |
| `NewConversationWindowBuffer(windowSize, opts...)` | 窗口对话缓冲 |
| `NewConversationTokenBuffer(llm, maxToken, opts...)` | 基于 Token 的对话缓冲 |
| `NewSimpleMemory() Simple` | 空记忆 |

### Embedding

| 函数 | 说明 |
|------|------|
| `NewEmbedder(client, opts...) (*EmbedderImpl, error)` | 创建文本向量化器 |
| `EmbedQuery(ctx, client, text) ([]float32, error)` | 便捷：单文本向量化 |
| `EmbedDocuments(ctx, client, texts) ([][]float32, error)` | 便捷：多文本向量化 |

### VectorStore

| 函数 | 说明 |
|------|------|
| `ToRetriever(store, numDocuments, opts...) Retriever` | 将 VectorStore 转为 Retriever |
| `AddDocuments(ctx, store, docs, opts...) ([]string, error)` | 添加文档 |
| `SimilaritySearch(ctx, store, query, num, opts...) ([]Document, error)` | 相似度搜索 |

## 使用示例

### 创建 LLM 并执行 Agent

```go
package example

import (
    "context"

    aiLC "github.com/tx7do/kratos-bootstrap/ai/langchaingo"
    "github.com/tmc/langchaingo/tools"

    conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func ExampleAgent() {
    ctx := context.Background()

    cfg := &conf.AI_Model{
        Type:      conf.AI_Model_CLOUD_MODEL,
        ModelName: "gpt-4o",
        Cloud: &conf.AI_Model_Cloud{
            ApiKey:  "sk-xxx",
            BaseUrl: "https://api.openai.com/v1",
        },
    }

    llm, _ := aiLC.NewModel(cfg)

    agentTools := []tools.Tool{/* 定义工具 */}

    executor := aiLC.NewOpenAIFunctionsExecutor(llm, agentTools)

    result, _ := executor.Call(ctx, "今天北京的天气如何？")
}
```

### 对话链（带记忆）

```go
func ExampleConversation() {
    ctx := context.Background()

    llm, _ := aiLC.NewModel(cfg)

    mem := aiLC.NewConversationBuffer()

    chain := aiLC.NewConversationChain(llm, mem)

    result1, _ := chain.Call(ctx, map[string]any{"input": "我叫小明"})
    result2, _ := chain.Call(ctx, map[string]any{"input": "我叫什么名字？"})
}
```

### 文本向量化与相似度搜索

```go
func ExampleVectorSearch() {
    ctx := context.Background()

    llm, _ := aiLC.NewModel(cfg)
    embedder, _ := aiLC.NewEmbedder(llm)

    // 创建向量存储（以 Redis 为例）
    store, _ := lcRedis.New(ctx, "redis://localhost:6379", "my_collection", embedder)

    // 添加文档
    aiLC.AddDocuments(ctx, store, []schema.Document{
        {PageContent: "Go 是一门静态类型语言", Metadata: map[string]any{"source": "wiki"}},
    })

    // 语义搜索
    docs, _ := aiLC.SimilaritySearch(ctx, store, "什么是 Go", 3)
}
```
