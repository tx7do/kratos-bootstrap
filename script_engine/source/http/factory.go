// Package httpsource 提供 HTTP 远程拉取作为脚本来源的集成。
//
// 独立 Go module，通过 init() 自动注册工厂。
// 用户只需导入即可启用：
//
//	import _ "github.com/tx7do/kratos-bootstrap/script_engine/source/http"
//
// YAML 配置示例:
//
//	script:
//	  source:
//	    type: HTTP
//	    options:
//	      base_url: "https://example.com/scripts/"
//	      cache_ttl: "10m"
package httpsource

import (
	"fmt"

	"github.com/tx7do/go-scripts/source"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	scriptEngine "github.com/tx7do/kratos-bootstrap/script_engine"
	// TODO: 当 go-scripts/source/http 发布后，取消注释以下导入
	// httpSource "github.com/tx7do/go-scripts/source/http"
)

func init() {
	scriptEngine.MustRegisterSourceFactory(conf.Script_Source_HTTP, NewSource)
}

// NewSource 根据配置创建 HTTP 脚本来源。
func NewSource(cfg *conf.Script_Source) (source.Reader, error) {
	if cfg == nil {
		return nil, fmt.Errorf("http source: config is nil")
	}

	opts := cfg.GetOptions()
	if opts == nil {
		return nil, fmt.Errorf("http source: options is required (base_url)")
	}
	m := opts.AsMap()

	baseURL, _ := m["base_url"].(string)
	if baseURL == "" {
		return nil, fmt.Errorf("http source: base_url is required")
	}

	// TODO: 当 go-scripts/source/http 发布后，替换以下占位实现

	return nil, fmt.Errorf(
		"http source: go-scripts/source/http not yet published — " +
			"this sub-module is a placeholder waiting for upstream release",
	)
}
