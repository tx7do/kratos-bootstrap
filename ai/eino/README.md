# eino 包说明

## 概述

`eino` 包基于字节跳动开源的 [Eino](https://github.com/cloudwego/eino) 框架封装，提供 LLM 应用开发能力。

Eino 是一个 Go 原生的 LLM 应用框架，支持以下核心能力：

- **ChatModel**：统一的模型调用接口
- **Chain / Graph / Workflow**：灵活的组件编排方式
- **Prompt 模板**：变量注入与消息模板
- **Tool**：工具定义与调用
- **Lambda**：自定义逻辑节点
- **分支与并行**：条件路由和并行执行

## 子模块说明

| 文件 | 说明 |
|------|------|
| `client.go` | ChatModel 客户端创建（云端 / Ollama 本地模型） |
| `chain.go` | Chain / Graph 编排与节点追加 |
| `compose.go` | Lambda、Parallel、Branch、Workflow、StateGraph 等编排组件 |
| `prompt.go` | 提示模板创建与消息构造辅助方法 |
| `tool.go` | 工具节点、工具信息与参数定义 |
| `options.go` | 可选配置项（ChatModelConfig 修饰器） |

## API 概览

### ChatModel

| 函数 | 说明 |
|------|------|
| `NewChatModel(ctx, cfg, opts...) (model.ChatModel, error)` | 根据配置创建 ChatModel |
| `WithConfigModifier(fn) Option` | 创建前修改 ChatModelConfig（如设置 Temperature） |

### Chain / Graph 编排

| 函数 | 说明 |
|------|------|
| `NewChain[I, O](opts...) *Chain` | 创建链式编排 |
| `NewGraph[I, O](opts...) *Graph` | 创建 DAG 工作流 |
| `NewWorkflow[I, O](opts...) *Workflow` | 创建带依赖管理的工作流 |
| `CompileChain(ctx, chain, opts...) (Runnable, error)` | 编译 Chain 为可执行单元 |

### Prompt 模板

| 函数 | 说明 |
|------|------|
| `FromMessages(formatType, templates...) *DefaultChatTemplate` | 从消息模板创建 ChatTemplate |
| `SystemMessage(content) *Message` | 创建系统消息 |
| `UserMessage(content) *Message` | 创建用户消息 |
| `AssistantMessage(content) *Message` | 创建助手消息 |
| `MessagesPlaceholder(name, optional) MessagesTemplate` | 创建消息占位符 |

### Tool

| 函数 | 说明 |
|------|------|
| `NewToolNode(ctx, config) (*ToolsNode, error)` | 创建工具执行节点 |
| `NewToolInfo(name, desc, params) *ToolInfo` | 创建工具描述 |
| `NewParameterInfo(typ, desc, required, opts...) *ParameterInfo` | 创建参数描述 |

### 编排组件

| 函数 | 说明 |
|------|------|
| `InvokableLambda(fn, opts...) *Lambda` | 创建同步 Lambda 节点 |
| `StreamableLambda(fn, opts...) *Lambda` | 创建流式 Lambda 节点 |
| `NewParallel() *Parallel` | 创建并行结构 |
| `NewChainBranch(cond) *ChainBranch` | 创建条件分支 |
| `WithGenLocalState(fn) NewGraphOption` | 为 Graph 注册本地状态 |

## 使用示例

### 创建 ChatModel 并对话

```go
package example

import (
    "context"

    aiEino "github.com/tx7do/kratos-bootstrap/ai/eino"

    conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func ExampleChat() {
    ctx := context.Background()

    cfg := &conf.AI_Model{
        Type:      conf.AI_Model_CLOUD_MODEL,
        ModelName: "gpt-4o",
        Cloud: &conf.AI_Model_Cloud{
            ApiKey:  "sk-xxx",
            BaseUrl: "https://api.openai.com/v1",
        },
    }

    chatModel, _ := aiEino.NewChatModel(ctx, cfg)

    result, _ := chatModel.Generate(ctx, []*schema.Message{
        {Role: schema.User, Content: "什么是 Go 语言？"},
    })
    fmt.Println(result.Content)
}
```

### 使用 Chain 编排

```go
func ExampleChain() {
    ctx := context.Background()

    chatModel, _ := aiEino.NewChatModel(ctx, cfg)

    tpl := aiEino.FromMessages(schema.FString,
        aiEino.SystemMessage("你是一个{role}助手。"),
        aiEino.UserMessage("{question}"),
    )

    chain := aiEino.NewChain[map[string]any, *schema.Message]()
    chain.
        AppendChatTemplate(tpl).
        AppendChatModel(chatModel)

    runnable, _ := aiEino.CompileChain(ctx, chain)

    result, _ := runnable.Invoke(ctx, map[string]any{
        "role":     "翻译",
        "question": "翻译这段话",
    })
}
```

### 使用 Graph + Tool 构建 Agent

```go
func ExampleAgent() {
    ctx := context.Background()

    chatModel, _ := aiEino.NewChatModel(ctx, cfg)

    // 定义工具
    searchTool := aiEino.NewToolInfo("search", "搜索信息", aiEino.NewParamsOneOfByParams(
        map[string]*schema.ParameterInfo{
            "query": aiEino.NewParameterInfo(schema.String, "搜索关键词", true),
        },
    ))

    toolsNode, _ := aiEino.NewToolNode(ctx, &aiEino.ToolsNodeConfig{
        Tools: []tool.BaseTool{/* 实现的工具 */},
    })

    // 构建 Graph
    graph := aiEino.NewGraph[map[string]any, *schema.Message]()
    graph.AddChatModelNode("model", chatModel)
    graph.AddToolsNode("tools", toolsNode)
    graph.AddEdge(aiEino.START, "model")
    graph.AddEdge("model", "tools")
    graph.AddEdge("tools", aiEino.END)

    runnable, _ := graph.Compile(ctx)
    result, _ := runnable.Invoke(ctx, map[string]any{"question": "今天天气如何？"})
}
```
