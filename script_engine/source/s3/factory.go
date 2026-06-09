// Package s3 提供 S3 / 兼容对象存储作为脚本来源的集成。
//
// 独立 Go module，通过 init() 自动注册工厂。
// 用户只需导入即可启用：
//
//	import _ "github.com/tx7do/kratos-bootstrap/script_engine/source/s3"
//
// YAML 配置示例:
//
//	script:
//	  source:
//	    type: S3
//	    options:
//	      bucket: "my-scripts"
//	      region: "us-east-1"
//	      prefix: "scripts/"
//	      cache_ttl: "5m"
package s3

import (
	"fmt"

	"github.com/tx7do/go-scripts/source"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	scriptEngine "github.com/tx7do/kratos-bootstrap/script_engine"
	// TODO: 当 go-scripts/source/s3 发布后，取消注释以下导入
	// s3Source "github.com/tx7do/go-scripts/source/s3"
)

func init() {
	scriptEngine.MustRegisterSourceFactory(conf.Script_Source_S3, NewSource)
}

// NewSource 根据配置创建 S3 脚本来源。
func NewSource(cfg *conf.Script_Source) (source.Reader, error) {
	if cfg == nil {
		return nil, fmt.Errorf("s3 source: config is nil")
	}

	opts := cfg.GetOptions()
	if opts == nil {
		return nil, fmt.Errorf("s3 source: options is required (bucket, region)")
	}
	m := opts.AsMap()

	bucket, _ := m["bucket"].(string)
	if bucket == "" {
		return nil, fmt.Errorf("s3 source: bucket is required")
	}

	region, _ := m["region"].(string)
	prefix, _ := m["prefix"].(string)

	_ = region
	_ = prefix

	// TODO: 当 go-scripts/source/s3 发布后，替换以下占位实现:
	//
	//   return s3Source.New(ctx, bucket,
	//       s3Source.WithRegion(region),
	//       s3Source.WithPrefix(prefix),
	//   )

	return nil, fmt.Errorf(
		"s3 source: go-scripts/source/s3 not yet published — " +
			"this sub-module is a placeholder waiting for upstream release",
	)
}
