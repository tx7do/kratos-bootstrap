package clickhouse

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

type QueryBuilder struct {
	table       string
	columns     []string
	distinct    bool
	conditions  []string
	orderBy     []string
	groupBy     []string
	having      []string
	joins       []string
	with        []string
	union       []string
	limit       int
	offset      int
	params      []interface{} // 用于存储参数
	useIndex    string        // 索引提示
	cacheResult bool          // 是否缓存查询结果
	debug       bool          // 是否启用调试
	log         *log.Helper
}

// NewQueryBuilder 创建一个新的 QueryBuilder 实例
func NewQueryBuilder(table string, log *log.Helper) *QueryBuilder {
	return &QueryBuilder{
		log:    log,
		table:  table,
		params: []interface{}{},
	}
}

// EnableDebug 启用调试模式
func (qb *QueryBuilder) EnableDebug() *QueryBuilder {
	qb.debug = true
	return qb
}

// logDebug 打印调试信息
func (qb *QueryBuilder) logDebug(message string) {
	if qb.debug {
		qb.log.Debug("[QueryBuilder Debug]:", message)
	}
}

// Select 设置查询的列
func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
	for _, column := range columns {
		if !isValidIdentifier(column) {
			panic("Invalid column name")
		}
	}

	qb.columns = columns
	return qb
}

// Distinct 设置 DISTINCT 查询
func (qb *QueryBuilder) Distinct() *QueryBuilder {
	qb.distinct = true
	return qb
}

// Where 添加查询条件并支持参数化
func (qb *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
	if !isValidCondition(condition) {
		panic("Invalid condition")
	}

	qb.conditions = append(qb.conditions, condition)
	qb.params = append(qb.params, args...)
	return qb
}

// OrderBy 设置排序条件
func (qb *QueryBuilder) OrderBy(order string) *QueryBuilder {
	qb.orderBy = append(qb.orderBy, order)
	return qb
}

// GroupBy 设置分组条件
func (qb *QueryBuilder) GroupBy(columns ...string) *QueryBuilder {
	qb.groupBy = append(qb.groupBy, columns...)
	return qb
}

// Having 添加分组后的过滤条件并支持参数化
func (qb *QueryBuilder) Having(condition string, args ...interface{}) *QueryBuilder {
	qb.having = append(qb.having, condition)
	qb.params = append(qb.params, args...)
	return qb
}

// Join 添加 JOIN 操作
func (qb *QueryBuilder) Join(joinType, table, onCondition string) *QueryBuilder {
	join := fmt.Sprintf("%s JOIN %s ON %s", joinType, table, onCondition)
	qb.joins = append(qb.joins, join)
	return qb
}

// With 添加 WITH 子句
func (qb *QueryBuilder) With(expression string) *QueryBuilder {
	qb.with = append(qb.with, expression)
	return qb
}

// Union 添加 UNION 操作
func (qb *QueryBuilder) Union(query string) *QueryBuilder {
	qb.union = append(qb.union, query)
	return qb
}

// Limit 设置查询结果的限制数量
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limit = limit
	return qb
}

// Offset 设置查询结果的偏移量
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offset = offset
	return qb
}

// UseIndex 设置索引提示
func (qb *QueryBuilder) UseIndex(index string) *QueryBuilder {
	qb.useIndex = index
	return qb
}

// CacheResult 启用查询结果缓存
func (qb *QueryBuilder) CacheResult() *QueryBuilder {
	qb.cacheResult = true
	return qb
}

// ArrayJoin 添加 ARRAY JOIN 子句
func (qb *QueryBuilder) ArrayJoin(expression string) *QueryBuilder {
	qb.joins = append(qb.joins, fmt.Sprintf("ARRAY JOIN %s", expression))
	return qb
}

// Final 添加 FINAL 修饰符
func (qb *QueryBuilder) Final() *QueryBuilder {
	qb.table = fmt.Sprintf("%s FINAL", qb.table)
	return qb
}

// Sample 添加 SAMPLE 子句
func (qb *QueryBuilder) Sample(sampleRate float64) *QueryBuilder {
	qb.table = fmt.Sprintf("%s SAMPLE %f", qb.table, sampleRate)
	return qb
}

// LimitBy 添加 LIMIT BY 子句
func (qb *QueryBuilder) LimitBy(limit int, columns ...string) *QueryBuilder {
	qb.limit = limit
	qb.orderBy = append(qb.orderBy, fmt.Sprintf("LIMIT BY %d (%s)", limit, strings.Join(columns, ", ")))
	return qb
}

// PreWhere 添加 PREWHERE 子句
func (qb *QueryBuilder) PreWhere(condition string, args ...interface{}) *QueryBuilder {
	qb.conditions = append([]string{condition}, qb.conditions...)
	qb.params = append(args, qb.params...)
	return qb
}

// Format 添加 FORMAT 子句
func (qb *QueryBuilder) Format(format string) *QueryBuilder {
	qb.union = append(qb.union, fmt.Sprintf("FORMAT %s", format))
	return qb
}

// Build 构建最终的 SQL 查询
func (qb *QueryBuilder) Build() (string, []interface{}) {
	query := ""

	if qb.cacheResult {
		query += "/* CACHE */ "
	}

	query += "SELECT "
	if qb.distinct {
		query += "DISTINCT "
	}
	query += qb.buildColumns()
	query += fmt.Sprintf(" FROM %s", qb.table)

	if qb.useIndex != "" {
		query += fmt.Sprintf(" USE INDEX (%s)", qb.useIndex)
	}

	if len(qb.conditions) > 0 {
		query += fmt.Sprintf(" WHERE %s", strings.Join(qb.conditions, " AND "))
	}

	if len(qb.groupBy) > 0 {
		query += fmt.Sprintf(" GROUP BY %s", strings.Join(qb.groupBy, ", "))
	}

	if len(qb.having) > 0 {
		query += fmt.Sprintf(" HAVING %s", strings.Join(qb.having, " AND "))
	}

	if len(qb.orderBy) > 0 {
		query += fmt.Sprintf(" ORDER BY %s", strings.Join(qb.orderBy, ", "))
	}

	if qb.limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", qb.limit)
	}

	if qb.offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", qb.offset)
	}

	return query, qb.params
}

func (qb *QueryBuilder) buildColumns() string {
	if len(qb.columns) == 0 {
		return "*"
	}
	return strings.Join(qb.columns, ", ")
}

// isValidIdentifier 验证表名或列名是否合法
func isValidIdentifier(identifier string) bool {
	// 仅允许字母、数字、下划线，且不能以数字开头
	matched, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, identifier)
	return matched
}

// isValidCondition 验证条件语句是否合法
func isValidCondition(condition string) bool {
	// 简单验证条件中是否包含危险字符
	return !strings.Contains(condition, ";") && !strings.Contains(condition, "--")
}
