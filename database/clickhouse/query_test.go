package clickhouse

import (
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
)

func TestQueryBuilder(t *testing.T) {
	logger := log.NewHelper(log.DefaultLogger)
	qb := NewQueryBuilder("test_table", logger)

	// 测试 Select 方法
	qb.Select("id", "name")
	query, params := qb.Build()
	assert.Contains(t, query, "SELECT id, name FROM test_table")

	// 测试 Distinct 方法
	qb.Distinct()
	query, _ = qb.Build()
	assert.Contains(t, query, "SELECT DISTINCT id, name FROM test_table")

	// 测试 Where 方法
	qb.Where("id > ?", 10).Where("name = ?", "example")
	query, params = qb.Build()
	assert.Contains(t, query, "WHERE id > ? AND name = ?")
	assert.Equal(t, []interface{}{10, "example"}, params)

	// 测试 OrderBy 方法
	qb.OrderBy("name ASC")
	query, _ = qb.Build()
	assert.Contains(t, query, "ORDER BY name ASC")

	// 测试 GroupBy 方法
	qb.GroupBy("category")
	query, _ = qb.Build()
	assert.Contains(t, query, "GROUP BY category")

	// 测试 Having 方法
	qb.Having("COUNT(id) > ?", 5)
	query, params = qb.Build()
	assert.Contains(t, query, "HAVING COUNT(id) > ?")
	assert.Equal(t, []interface{}{10, "example", 5}, params)

	// 测试 Join 方法
	qb.Join("INNER", "other_table", "test_table.id = other_table.id")
	query, _ = qb.Build()
	assert.Contains(t, query, "INNER JOIN other_table ON test_table.id = other_table.id")

	// 测试 With 方法
	qb.With("temp AS (SELECT id FROM another_table WHERE status = 'active')")
	query, _ = qb.Build()
	assert.Contains(t, query, "WITH temp AS (SELECT id FROM another_table WHERE status = 'active')")

	// 测试 Union 方法
	qb.Union("SELECT id FROM another_table")
	query, _ = qb.Build()
	assert.Contains(t, query, "UNION SELECT id FROM another_table")

	// 测试 Limit 和 Offset 方法
	qb.Limit(10).Offset(20)
	query, _ = qb.Build()
	assert.Contains(t, query, "LIMIT 10 OFFSET 20")

	// 测试 UseIndex 方法
	qb.UseIndex("idx_name")
	query, _ = qb.Build()
	assert.Contains(t, query, "USE INDEX (idx_name)")

	// 测试 CacheResult 方法
	qb.CacheResult()
	query, _ = qb.Build()
	assert.Contains(t, query, "/* CACHE */")

	// 测试 EnableDebug 方法
	qb.EnableDebug()
	assert.True(t, qb.debug)

	// 测试 ArrayJoin 方法
	qb.ArrayJoin("array_column")
	query, _ = qb.Build()
	assert.Contains(t, query, "ARRAY JOIN array_column")

	// 测试 Final 方法
	qb.Final()
	query, _ = qb.Build()
	assert.Contains(t, query, "test_table FINAL")

	// 测试 Sample 方法
	qb.Sample(0.1)
	query, _ = qb.Build()
	assert.Contains(t, query, "test_table SAMPLE 0.100000")

	// 测试 LimitBy 方法
	qb.LimitBy(5, "name")
	query, _ = qb.Build()
	assert.Contains(t, query, "LIMIT BY 5 (name)")

	// 测试 PreWhere 方法
	qb.PreWhere("status = ?", "active")
	query, params = qb.Build()
	assert.Contains(t, query, "PREWHERE status = ?")
	assert.Equal(t, []interface{}{"active"}, params)

	// 测试 Format 方法
	qb.Format("JSON")
	query, _ = qb.Build()
	assert.Contains(t, query, "FORMAT JSON")

	// 测试边界情况：空列名
	assert.Panics(t, func() {
		qb.Select("")
	}, "应该抛出异常：无效的列名")

	// 测试边界情况：无效条件
	assert.Panics(t, func() {
		qb.Where("id = 1; DROP TABLE test_table")
	}, "应该抛出异常：无效的条件")
}
