// Package redis 提供 Redis 作为脚本来源的集成。
//
// 独立 Go module，通过 init() 自动注册工厂。
// 用户只需导入即可启用：
//
//	import _ "github.com/tx7do/kratos-bootstrap/script_engine/source/redis"
//
// YAML 配置示例:
//
//	script:
//	  source:
//	    type: REDIS
//	    options:
//	      addr: "localhost:6379"
//	      db: 0
//	      prefix: "script:"
//	      cache_ttl: "2m"
package redis

import (
	"fmt"

	"github.com/tx7do/go-scripts/source"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	scriptEngine "github.com/tx7do/kratos-bootstrap/script_engine"
	// TODO: 当 go-scripts/source/redis 发布后，取消注释以下导入
	// redisSource "github.com/tx7do/go-scripts/source/redis"
)

func init() {
	scriptEngine.MustRegisterSourceFactory(conf.Script_Source_REDIS, NewSource)
}

// NewSource 根据配置创建 Redis 脚本来源。
func NewSource(cfg *conf.Script_Source) (source.Reader, error) {
	if cfg == nil {
		return nil, fmt.Errorf("redis source: config is nil")
	}

	opts := cfg.GetOptions()
	if opts == nil {
		return nil, fmt.Errorf("redis source: options is required (addr)")
	}
	m := opts.AsMap()

	addr, _ := m["addr"].(string)
	if addr == "" {
		return nil, fmt.Errorf("redis source: addr is required")
	}

	prefix, _ := m["prefix"].(string)
	_ = prefix

	// TODO: 当 go-scripts/source/redis 发布后，替换以下占位实现:
	//
	//   return redisSource.New(addr,
	//       redisSource.WithPrefix(prefix),
	//   )

	return nil, fmt.Errorf(
		"redis source: go-scripts/source/redis not yet published — " +
			"this sub-module is a placeholder waiting for upstream release",
	)
}
