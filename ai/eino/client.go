package eino

import (
	"context"
	"errors"
	"fmt"
	"time"

	einoOpenai "github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewChatModel 根据配置创建 Eino ChatModel。
// 支持云端模型（OpenAI 兼容 API）和本地模型（Ollama）。
func NewChatModel(ctx context.Context, cfg *conf.AI_Model) (model.ChatModel, error) {
	if cfg == nil {
		return nil, errors.New("ai model config is nil")
	}

	switch cfg.Type {
	case conf.AI_Model_LOCAL_MODEL:
		return newOllamaChatModel(ctx, cfg)
	case conf.AI_Model_CLOUD_MODEL:
		return newCloudChatModel(ctx, cfg)
	default:
		return nil, fmt.Errorf("unsupported ai model type: %v", cfg.Type)
	}
}

// newCloudChatModel 创建云端模型（基于 Eino OpenAI 实现）。
func newCloudChatModel(ctx context.Context, cfg *conf.AI_Model) (model.ChatModel, error) {
	if cfg.Cloud == nil {
		return nil, errors.New("cloud config is nil")
	}

	config := &einoOpenai.ChatModelConfig{
		APIKey: cfg.Cloud.ApiKey,
		Model:  cfg.GetModelName(),
	}

	if cfg.Cloud.BaseUrl != "" {
		config.BaseURL = cfg.Cloud.BaseUrl
	}
	if cfg.GetTimeoutSeconds() > 0 {
		config.Timeout = time.Duration(cfg.GetTimeoutSeconds()) * time.Second
	}

	return einoOpenai.NewChatModel(ctx, config)
}

// newOllamaChatModel 创建本地模型（基于 Eino OpenAI 实现，兼容 Ollama）。
func newOllamaChatModel(ctx context.Context, cfg *conf.AI_Model) (model.ChatModel, error) {
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

	config := &einoOpenai.ChatModelConfig{
		APIKey:  "ollama",
		BaseURL: fmt.Sprintf("http://%s:%d/v1", host, port),
		Model:   cfg.GetModelName(),
	}
	if cfg.GetTimeoutSeconds() > 0 {
		config.Timeout = time.Duration(cfg.GetTimeoutSeconds()) * time.Second
	}

	return einoOpenai.NewChatModel(ctx, config)
}
