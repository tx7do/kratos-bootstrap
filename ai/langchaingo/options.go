package langchaingo

import lcOpenai "github.com/tmc/langchaingo/llms/openai"

// Option 是 LangChainGo 客户端的可选配置项。
type Option func(*options)

type options struct {
	openaiOpts []lcOpenai.Option
}

// WithOpenaiOptions 追加 LangChainGo OpenAI 原生选项。
func WithOpenaiOptions(opts ...lcOpenai.Option) Option {
	return func(o *options) {
		o.openaiOpts = append(o.openaiOpts, opts...)
	}
}
