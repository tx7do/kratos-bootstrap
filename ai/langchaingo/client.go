package langchaingo

import (
	"errors"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	lcOpenai "github.com/tmc/langchaingo/llms/openai"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewModel 根据配置创建 LangChainGo LLM 客户端。
// 支持云端模型（OpenAI 兼容 API）和本地模型（Ollama）。
func NewModel(cfg *conf.AI_Model, opts ...Option) (llms.Model, error) {
	if cfg == nil {
		return nil, errors.New("ai model config is nil")
	}

	o := applyOptions(opts)

	switch cfg.Type {
	case conf.AI_Model_LOCAL_MODEL:
		return newOllamaModel(cfg, o)
	case conf.AI_Model_CLOUD_MODEL:
		return newCloudModel(cfg, o)
	default:
		return nil, fmt.Errorf("unsupported ai model type: %v", cfg.Type)
	}
}

// newCloudModel 创建云端模型（基于 LangChainGo OpenAI 实现）。
func newCloudModel(cfg *conf.AI_Model, o *options) (llms.Model, error) {
	if cfg.Cloud == nil {
		return nil, errors.New("cloud config is nil")
	}

	opts := []lcOpenai.Option{
		lcOpenai.WithToken(cfg.Cloud.ApiKey),
		lcOpenai.WithModel(cfg.GetModelName()),
	}

	if cfg.Cloud.BaseUrl != "" {
		opts = append(opts, lcOpenai.WithBaseURL(cfg.Cloud.BaseUrl))
	}

	// 追加用户自定义的 OpenAI 选项
	opts = append(opts, o.openaiOpts...)

	return lcOpenai.New(opts...)
}

// newOllamaModel 创建本地模型（基于 LangChainGo Ollama 实现）。
func newOllamaModel(cfg *conf.AI_Model, o *options) (llms.Model, error) {
	return nil, errors.New("langchaingo ollama: not yet implemented")
}

// applyOptions 应用可选配置项。
func applyOptions(opts []Option) *options {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}
