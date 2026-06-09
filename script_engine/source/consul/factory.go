// Package consul 提供 Consul KV 作为脚本来源的集成。
//
// 独立 Go module，通过 init() 自动注册工厂。
// 用户只需导入即可启用：
//
//	import _ "github.com/tx7do/kratos-bootstrap/script_engine/source/consul"
//
// YAML 配置示例:
//
//	script:
//	  source:
//	    type: CONSUL
//	    options:
//	      address: "localhost:8500"
//	      prefix: "scripts/"
package consul

import (
	"fmt"

	"github.com/tx7do/go-scripts/source"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	scriptEngine "github.com/tx7do/kratos-bootstrap/script_engine"
	// TODO: 当 go-scripts/source/consul 发布后，取消注释以下导入
	// consulSource "github.com/tx7do/go-scripts/source/consul"
)

func init() {
	scriptEngine.MustRegisterSourceFactory(conf.Script_Source_CONSUL, NewSource)
}

// NewSource 根据配置创建 Consul 脚本来源。
func NewSource(cfg *conf.Script_Source) (source.Reader, error) {
	if cfg == nil {
		return nil, fmt.Errorf("consul source: config is nil")
	}

	opts := cfg.GetOptions()
	if opts == nil {
		return nil, fmt.Errorf("consul source: options is required (address)")
	}
	m := opts.AsMap()

	address, _ := m["address"].(string)
	if address == "" {
		return nil, fmt.Errorf("consul source: address is required")
	}

	prefix, _ := m["prefix"].(string)
	_ = prefix

	// TODO: 当 go-scripts/source/consul 发布后，替换以下占位实现

	return nil, fmt.Errorf(
		"consul source: go-scripts/source/consul not yet published — " +
			"this sub-module is a placeholder waiting for upstream release",
	)
}
