package langchaingo

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/tmc/langchaingo/llms"
	lcOllama "github.com/tmc/langchaingo/llms/ollama"
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

	// 设置 HTTP 客户端
	httpClient := o.httpClient
	if httpClient == nil {
		timeout := 30 * time.Second
		if cfg.TimeoutSeconds > 0 {
			timeout = time.Duration(cfg.TimeoutSeconds) * time.Second
		}
		httpClient = &http.Client{Timeout: timeout}
	}
	opts = append(opts, lcOpenai.WithHTTPClient(httpClient))

	// 追加用户自定义的 OpenAI 选项
	opts = append(opts, o.openaiOpts...)

	return lcOpenai.New(opts...)
}

// newOllamaModel 创建本地模型（基于 LangChainGo Ollama 实现）。
func newOllamaModel(cfg *conf.AI_Model, o *options) (llms.Model, error) {
	if cfg.Local == nil {
		return nil, errors.New("local config is nil")
	}

	host := cfg.Local.Host
	if host == "" {
		host = "localhost"
	}
	port := cfg.Local.Port
	if port == 0 {
		port = 11434
	}

	opts := []lcOllama.Option{
		lcOllama.WithModel(cfg.GetModelName()),
		lcOllama.WithServerURL(fmt.Sprintf("http://%s:%d", host, port)),
	}

	// 设置 HTTP 客户端
	if o.httpClient != nil {
		opts = append(opts, lcOllama.WithHTTPClient(o.httpClient))
	}

	// 追加用户自定义的 Ollama 选项
	opts = append(opts, o.ollamaOpts...)

	return lcOllama.New(opts...)
}

// applyOptions 应用可选配置项。
func applyOptions(opts []Option) *options {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}
