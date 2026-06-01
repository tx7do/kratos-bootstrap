package eino

import (
	einoOpenai "github.com/cloudwego/eino-ext/components/model/openai"
)

// Option 是 Eino 客户端的可选配置项。
type Option func(*options)

type options struct {
	// OpenAI ChatModel 原生配置修饰器
	// 可用于修改 ChatModelConfig 的额外字段
	configModifier func(*einoOpenai.ChatModelConfig)
}

// WithConfigModifier 提供一个函数在创建 ChatModel 前修改 ChatModelConfig。
// 可用于设置 einoOpenai.ChatModelConfig 中的额外字段（如 Temperature、MaxTokens 等）。
//
// 用法示例：
//
//	langchaingo.WithConfigModifier(func(cfg *einoOpenai.ChatModelConfig) {
//	    cfg.Temperature = ptr.To(float32(0.7))
//	    cfg.MaxTokens = ptr.To(4096)
//	})
func WithConfigModifier(modifier func(*einoOpenai.ChatModelConfig)) Option {
	return func(o *options) {
		o.configModifier = modifier
	}
}

// applyOptions 应用可选配置项。
func applyOptions(opts []Option) *options {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}
