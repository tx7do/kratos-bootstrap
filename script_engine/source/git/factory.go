// Package git 提供 Git 仓库作为脚本来源的集成。
//
// 独立 Go module，通过 init() 自动注册工厂。
// 用户只需导入即可启用：
//
//	import _ "github.com/tx7do/kratos-bootstrap/script_engine/source/git"
//
// YAML 配置示例:
//
//	script:
//	  source:
//	    type: GIT
//	    options:
//	      url: "https://github.com/org/scripts.git"
//	      branch: "main"
//	      subdir: "lua"
//	      cache_ttl: "30m"
package git

import (
	"fmt"

	"github.com/tx7do/go-scripts/source"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	scriptEngine "github.com/tx7do/kratos-bootstrap/script_engine"
	// TODO: 当 go-scripts/source/git 发布后，取消注释以下导入
	// gitSource "github.com/tx7do/go-scripts/source/git"
)

func init() {
	scriptEngine.MustRegisterSourceFactory(conf.Script_Source_GIT, NewSource)
}

// NewSource 根据配置创建 Git 脚本来源。
func NewSource(cfg *conf.Script_Source) (source.Reader, error) {
	if cfg == nil {
		return nil, fmt.Errorf("git source: config is nil")
	}

	opts := cfg.GetOptions()
	if opts == nil {
		return nil, fmt.Errorf("git source: options is required (url)")
	}
	m := opts.AsMap()

	url, _ := m["url"].(string)
	if url == "" {
		return nil, fmt.Errorf("git source: url is required")
	}

	branch, _ := m["branch"].(string)
	subdir, _ := m["subdir"].(string)
	_ = branch
	_ = subdir

	// TODO: 当 go-scripts/source/git 发布后，替换以下占位实现

	return nil, fmt.Errorf(
		"git source: go-scripts/source/git not yet published — " +
			"this sub-module is a placeholder waiting for upstream release",
	)
}
