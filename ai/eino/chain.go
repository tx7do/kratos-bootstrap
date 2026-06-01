package eino

import (
	"context"

	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/indexer"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/compose"
)

// NewChain 创建一个泛型 Chain，用于将多个组件按顺序串联。
// Chain 是 Eino 最核心的编排方式，通过 Append* 方法链式添加节点，最后 Compile 成可执行的 Runnable。
//
// 用法示例：
//
//	chain := eino.NewChain[map[string]any, *schema.Message]()
//	chain.
//	    AppendChatTemplate(chatTemplate).
//	    AppendChatModel(chatModel)
//	runnable, _ := chain.Compile(ctx)
//	result, _ := runnable.Invoke(ctx, map[string]any{"question": "什么是 Go？"})
func NewChain[I, O any](opts ...compose.NewGraphOption) *compose.Chain[I, O] {
	return compose.NewChain[I, O](opts...)
}

// NewGraph 创建一个泛型 Graph，用于构建更复杂的 DAG 工作流。
// Graph 比 Chain 更灵活，支持并行、分支、循环等高级编排模式。
//
// 用法示例：
//
//	graph := eino.NewGraph[map[string]any, *schema.Message]()
//	graph.AddChatTemplateNode("prompt", chatTemplate)
//	graph.AddChatModelNode("model", chatModel)
//	graph.AddEdge(compose.START, "prompt")
//	graph.AddEdge("prompt", "model")
//	graph.AddEdge("model", compose.END)
//	runnable, _ := graph.Compile(ctx)
func NewGraph[I, O any](opts ...compose.NewGraphOption) *compose.Graph[I, O] {
	return compose.NewGraph[I, O](opts...)
}

// START 和 END 是 Graph 中用于标识起始和结束节点的常量。
const (
	START = compose.START
	END   = compose.END
)

// --- Chain 节点追加方法（包装 compose.Chain 的 Append* 方法）---

// AppendChatModel 向 Chain 追加 ChatModel 节点。
func AppendChatModel[I, O any](chain *compose.Chain[I, O], node model.BaseChatModel, opts ...compose.GraphAddNodeOpt) *compose.Chain[I, O] {
	return chain.AppendChatModel(node, opts...)
}

// AppendChatTemplate 向 Chain 追加 ChatTemplate 节点。
func AppendChatTemplate[I, O any](chain *compose.Chain[I, O], node prompt.ChatTemplate, opts ...compose.GraphAddNodeOpt) *compose.Chain[I, O] {
	return chain.AppendChatTemplate(node, opts...)
}

// AppendToolsNode 向 Chain 追加 ToolsNode 节点。
func AppendToolsNode[I, O any](chain *compose.Chain[I, O], node *compose.ToolsNode, opts ...compose.GraphAddNodeOpt) *compose.Chain[I, O] {
	return chain.AppendToolsNode(node, opts...)
}

// AppendEmbedding 向 Chain 追加 Embedding 节点。
func AppendEmbedding[I, O any](chain *compose.Chain[I, O], node embedding.Embedder, opts ...compose.GraphAddNodeOpt) *compose.Chain[I, O] {
	return chain.AppendEmbedding(node, opts...)
}

// AppendRetriever 向 Chain 追加 Retriever 节点。
func AppendRetriever[I, O any](chain *compose.Chain[I, O], node retriever.Retriever, opts ...compose.GraphAddNodeOpt) *compose.Chain[I, O] {
	return chain.AppendRetriever(node, opts...)
}

// AppendLoader 向 Chain 追加文档加载器节点。
func AppendLoader[I, O any](chain *compose.Chain[I, O], node document.Loader, opts ...compose.GraphAddNodeOpt) *compose.Chain[I, O] {
	return chain.AppendLoader(node, opts...)
}

// AppendDocumentTransformer 向 Chain 追加文档转换器节点。
func AppendDocumentTransformer[I, O any](chain *compose.Chain[I, O], node document.Transformer, opts ...compose.GraphAddNodeOpt) *compose.Chain[I, O] {
	return chain.AppendDocumentTransformer(node, opts...)
}

// AppendIndexer 向 Chain 追加 Indexer 节点。
func AppendIndexer[I, O any](chain *compose.Chain[I, O], node indexer.Indexer, opts ...compose.GraphAddNodeOpt) *compose.Chain[I, O] {
	return chain.AppendIndexer(node, opts...)
}

// CompileChain 编译 Chain 为可执行的 Runnable。
func CompileChain[I, O any](ctx context.Context, chain *compose.Chain[I, O], opts ...compose.GraphCompileOption) (compose.Runnable[I, O], error) {
	return chain.Compile(ctx, opts...)
}
