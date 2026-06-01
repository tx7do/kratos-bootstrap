package langchaingo

import (
	"context"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

// VectorStore 是向量存储的通用接口，由各后端实现（Redis、pgvector、Milvus 等）。
// 用户可以通过此接口统一调用 AddDocuments 和 SimilaritySearch。
//
// 用法示例（以 Redis 为例）：
//
//	import lcRedis "github.com/tmc/langchaingo/vectorstores/redisvector"
//
//	llm, _ := langchaingo.NewModel(cfg)
//	embedder, _ := langchaingo.NewEmbedder(llm)
//
//	store, _ := lcRedis.New(
//	  context.Background(),
//	  "redis://localhost:6379",
//	  "my_collection",
//	  embedder,
//	)
//
//	// 添加文档
//	store.AddDocuments(ctx, []schema.Document{
//	  {PageContent: "Go 是一门静态类型语言", Metadata: map[string]any{"source": "wiki"}},
//	})
//
//	// 语义搜索
//	docs, _ := store.SimilaritySearch(ctx, "什么是 Go", 3)

// ToRetriever 将 VectorStore 转换为 Retriever，可嵌入 Chain 使用。
// numDocuments 指定检索时返回的最大文档数。
func ToRetriever(store vectorstores.VectorStore, numDocuments int, opts ...vectorstores.Option) vectorstores.Retriever {
	return vectorstores.ToRetriever(store, numDocuments, opts...)
}

// AddDocuments 向向量存储中添加文档。
// 这是 vectorstores.VectorStore 接口方法的快捷调用。
func AddDocuments(ctx context.Context, store vectorstores.VectorStore, docs []schema.Document, opts ...vectorstores.Option) ([]string, error) {
	return store.AddDocuments(ctx, docs, opts...)
}

// SimilaritySearch 在向量存储中进行相似度搜索。
// query 为查询文本，numDocuments 为返回的最大文档数。
func SimilaritySearch(ctx context.Context, store vectorstores.VectorStore, query string, numDocuments int, opts ...vectorstores.Option) ([]schema.Document, error) {
	return store.SimilaritySearch(ctx, query, numDocuments, opts...)
}

// --- VectorStore Option 透传 ---

// WithNameSpace 设置向量存储的命名空间。
func WithNameSpace(nameSpace string) vectorstores.Option {
	return vectorstores.WithNameSpace(nameSpace)
}

// WithScoreThreshold 设置相似度分数阈值，低于此阈值的文档将被过滤。
func WithScoreThreshold(scoreThreshold float32) vectorstores.Option {
	return vectorstores.WithScoreThreshold(scoreThreshold)
}

// WithFilters 设置元数据过滤条件。
func WithFilters(filters any) vectorstores.Option {
	return vectorstores.WithFilters(filters)
}

// WithEmbedder 设置用于向量化的 Embedder。
// 当需要在单个 VectorStore 中使用多个 LLM 时很有用。
func WithEmbedder(embedder embeddings.Embedder) vectorstores.Option {
	return vectorstores.WithEmbedder(embedder)
}

// WithDeduplicater 设置去重函数，在添加文档时跳过已存在的文档。
func WithDeduplicater(fn func(ctx context.Context, doc schema.Document) bool) vectorstores.Option {
	return vectorstores.WithDeduplicater(fn)
}
