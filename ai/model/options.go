package model

import "net/http"

// Option 是 Agent 客户端的可选配置项。
type Option func(*options)

type options struct {
	httpClient *http.Client
}

// WithHTTPClient 设置自定义 HTTP 客户端。
func WithHTTPClient(httpClient *http.Client) Option {
	return func(o *options) {
		o.httpClient = httpClient
	}
}
