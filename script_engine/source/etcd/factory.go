// Package etcd 提供 etcd 作为脚本来源的集成。
//
// 这是一个独立 Go module，通过 init() 自动注册工厂到 script_engine 主包。
// 用户只需导入即可启用：
//
//	import _ "github.com/tx7do/kratos-bootstrap/script_engine/source/etcd"
//
// YAML 配置示例:
//
//	script:
//	  source:
//	    type: ETCD
//	    options:
//	      endpoints: ["localhost:2379"]
//	      prefix: "/scripts/"
//	      cache_ttl: "5m"
package etcd

import (
	"fmt"

	"github.com/tx7do/go-scripts/source"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	scriptEngine "github.com/tx7do/kratos-bootstrap/script_engine"
	// TODO: 当 go-scripts/source/etcd 发布后，取消注释以下导入
	// etcdSource "github.com/tx7do/go-scripts/source/etcd"
)

func init() {
	scriptEngine.MustRegisterSourceFactory(conf.Script_Source_ETCD, NewSource)
}

// NewSource 根据配置创建 etcd 脚本来源。
func NewSource(cfg *conf.Script_Source) (source.Reader, error) {
	if cfg == nil {
		return nil, fmt.Errorf("etcd source: config is nil")
	}

	opts := cfg.GetOptions()
	if opts == nil {
		return nil, fmt.Errorf("etcd source: options is required (endpoints, prefix)")
	}
	m := opts.AsMap()

	// 解析配置
	endpoints := toStringSlice(m["endpoints"])
	if len(endpoints) == 0 {
		return nil, fmt.Errorf("etcd source: at least one endpoint is required")
	}

	prefix, _ := m["prefix"].(string)

	_ = endpoints
	_ = prefix

	// TODO: 当 go-scripts/source/etcd 发布后，替换以下占位实现:
	//
	//   return etcdSource.New(
	//       etcdSource.WithEndpoints(endpoints...),
	//       etcdSource.WithPrefix(prefix),
	//   )

	return nil, fmt.Errorf(
		"etcd source: go-scripts/source/etcd not yet published — " +
			"this sub-module is a placeholder waiting for upstream release",
	)
}

// toStringSlice 将 any 转为 []string。
func toStringSlice(v any) []string {
	if v == nil {
		return nil
	}
	if slice, ok := v.([]any); ok {
		result := make([]string, 0, len(slice))
		for _, item := range slice {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	}
	if s, ok := v.(string); ok {
		return []string{s}
	}
	return nil
}
