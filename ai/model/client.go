package model

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sashabaranov/go-openai"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewClient 根据配置创建 OpenAI 兼容的客户端。
// 支持云端模型（OpenAI、通义千问等）和本地模型（Ollama）。
func NewClient(cfg *conf.AI_Model, opts ...Option) (*openai.Client, error) {
	if cfg == nil {
		return nil, errors.New("ai config is nil")
	}

	o := applyOptions(opts)

	switch cfg.Type {
	case conf.AI_Model_LOCAL_MODEL:
		return newLocalClient(cfg, o)
	case conf.AI_Model_CLOUD_MODEL:
		return newCloudClient(cfg, o)
	default:
		return nil, fmt.Errorf("unsupported ai model type: %v", cfg.Type)
	}
}

// newCloudClient 创建云端模型客户端。
func newCloudClient(cfg *conf.AI_Model, o *options) (*openai.Client, error) {
	if cfg.Cloud == nil {
		return nil, errors.New("cloud config is nil")
	}

	c := openai.DefaultConfig(cfg.Cloud.ApiKey)

	if cfg.Cloud.BaseUrl != "" {
		c.BaseURL = cfg.Cloud.BaseUrl
	}
	if cfg.Cloud.Organization != "" {
		c.OrgID = cfg.Cloud.Organization
	}

	setHTTPClient(cfg, o, &c)

	client := openai.NewClientWithConfig(c)
	return client, nil
}

// newLocalClient 创建本地模型客户端（Ollama）。
func newLocalClient(cfg *conf.AI_Model, o *options) (*openai.Client, error) {
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

	c := openai.DefaultConfig("none")
	c.BaseURL = fmt.Sprintf("http://%s:%d/v1", host, port)

	setHTTPClient(cfg, o, &c)

	client := openai.NewClientWithConfig(c)
	return client, nil
}

// setHTTPClient 设置自定义 HTTP 客户端（含超时）。
func setHTTPClient(cfg *conf.AI_Model, o *options, c *openai.ClientConfig) {
	if o.httpClient != nil {
		c.HTTPClient = o.httpClient
		return
	}

	timeout := 30 * time.Second
	if cfg.TimeoutSeconds > 0 {
		timeout = time.Duration(cfg.TimeoutSeconds) * time.Second
	}
	c.HTTPClient = &http.Client{Timeout: timeout}
}

// applyOptions 应用可选配置项。
func applyOptions(opts []Option) *options {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}
