package langchaingo

import (
	"net/http"

	lcOllama "github.com/tmc/langchaingo/llms/ollama"
	lcOpenai "github.com/tmc/langchaingo/llms/openai"
)

// Option 是 LangChainGo 客户端的可选配置项。
type Option func(*options)

type options struct {
	// OpenAI 云端模型原生选项
	openaiOpts []lcOpenai.Option

	// Ollama 本地模型原生选项
	ollamaOpts []lcOllama.Option

	// HTTP 客户端（用于 OpenAI 和 Ollama）
	httpClient *http.Client
}

// WithOpenaiOptions 追加 LangChainGo OpenAI 原生选项。
func WithOpenaiOptions(opts ...lcOpenai.Option) Option {
	return func(o *options) {
		o.openaiOpts = append(o.openaiOpts, opts...)
	}
}

// WithOllamaOptions 追加 LangChainGo Ollama 原生选项。
func WithOllamaOptions(opts ...lcOllama.Option) Option {
	return func(o *options) {
		o.ollamaOpts = append(o.ollamaOpts, opts...)
	}
}

// WithHTTPClient 设置自定义 HTTP 客户端（用于 OpenAI 和 Ollama）。
func WithHTTPClient(httpClient *http.Client) Option {
	return func(o *options) {
		o.httpClient = httpClient
	}
}
