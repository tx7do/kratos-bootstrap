package langchaingo

import (
	"context"

	"github.com/tmc/langchaingo/embeddings"
)

// NewEmbedder 基于 EmbedderClient 创建一个文本向量化器（Embedder）。
// OpenAI LLM 和 Ollama LLM 都实现了 EmbedderClient 接口，可直接使用。
//
// 用法示例：
//
//	llm, _ := openai.New(openai.WithToken("sk-xxx"))  // 或 ollama.New(ollama.WithModel("llama3"))
//	embedder, _ := langchaingo.NewEmbedder(llm)
//	vector, _ := embedder.EmbedQuery(ctx, "什么是 Go 语言？")
//	vectors, _ := embedder.EmbedDocuments(ctx, []string{"文档1", "文档2"})
func NewEmbedder(client embeddings.EmbedderClient, opts ...embeddings.Option) (*embeddings.EmbedderImpl, error) {
	return embeddings.NewEmbedder(client, opts...)
}

// EmbedQuery 便捷方法：基于 EmbedderClient 对单个文本进行向量化。
// 无需先创建 Embedder，适合简单的单次查询场景。
func EmbedQuery(ctx context.Context, client embeddings.EmbedderClient, text string) ([]float32, error) {
	embedder, err := embeddings.NewEmbedder(client)
	if err != nil {
		return nil, err
	}
	return embedder.EmbedQuery(ctx, text)
}

// EmbedDocuments 便捷方法：基于 EmbedderClient 对多个文本进行向量化。
func EmbedDocuments(ctx context.Context, client embeddings.EmbedderClient, texts []string) ([][]float32, error) {
	embedder, err := embeddings.NewEmbedder(client)
	if err != nil {
		return nil, err
	}
	return embedder.EmbedDocuments(ctx, texts)
}

// --- Embedding Option 透传 ---

// WithStripNewLines 设置是否去除文本中的换行符（默认 true）。
func WithStripNewLines(stripNewLines bool) embeddings.Option {
	return embeddings.WithStripNewLines(stripNewLines)
}

// WithBatchSize 设置批量处理的文档数量（默认 512）。
func WithBatchSize(batchSize int) embeddings.Option {
	return embeddings.WithBatchSize(batchSize)
}
