// Package database 提供 SQL 数据库作为脚本来源的集成。
//
// 独立 Go module，通过 init() 自动注册工厂。
// 用户只需导入即可启用：
//
//	import _ "github.com/tx7do/kratos-bootstrap/script_engine/source/database"
//
// YAML 配置示例:
//
//	script:
//	  source:
//	    type: DATABASE
//	    options:
//	      dsn: "postgres://user:pass@localhost:5432/scripts"
//	      table: "scripts"
//	      key_column: "name"
//	      content_column: "code"
//	      cache_ttl: "5m"
package database

import (
	"fmt"

	"github.com/tx7do/go-scripts/source"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	scriptEngine "github.com/tx7do/kratos-bootstrap/script_engine"
	// TODO: 当 go-scripts/source/database 发布后，取消注释以下导入
	// dbSource "github.com/tx7do/go-scripts/source/database"
)

func init() {
	scriptEngine.MustRegisterSourceFactory(conf.Script_Source_DATABASE, NewSource)
}

// NewSource 根据配置创建数据库脚本来源。
func NewSource(cfg *conf.Script_Source) (source.Reader, error) {
	if cfg == nil {
		return nil, fmt.Errorf("database source: config is nil")
	}

	opts := cfg.GetOptions()
	if opts == nil {
		return nil, fmt.Errorf("database source: options is required (dsn, table)")
	}
	m := opts.AsMap()

	dsn, _ := m["dsn"].(string)
	if dsn == "" {
		return nil, fmt.Errorf("database source: dsn is required")
	}

	table, _ := m["table"].(string)
	if table == "" {
		return nil, fmt.Errorf("database source: table is required")
	}

	keyColumn, _ := m["key_column"].(string)
	contentColumn, _ := m["content_column"].(string)
	_ = keyColumn
	_ = contentColumn

	// TODO: 当 go-scripts/source/database 发布后，替换以下占位实现

	return nil, fmt.Errorf(
		"database source: go-scripts/source/database not yet published — " +
			"this sub-module is a placeholder waiting for upstream release",
	)
}
