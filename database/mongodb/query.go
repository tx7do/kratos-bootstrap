package mongodb

import (
	bsonV2 "go.mongodb.org/mongo-driver/v2/bson"
	optionsV2 "go.mongodb.org/mongo-driver/v2/mongo/options"
)

type QueryBuilder struct {
	filter   bsonV2.M
	opts     *optionsV2.FindOptions
	pipeline []bsonV2.M
}

func NewQuery() *QueryBuilder {
	return &QueryBuilder{
		filter: bsonV2.M{},
		opts:   &optionsV2.FindOptions{},
	}
}

// SetFilter 设置查询过滤条件
func (qb *QueryBuilder) SetFilter(filter bsonV2.M) *QueryBuilder {
	qb.filter = filter
	return qb
}

// SetOr 设置多个条件的逻辑或
func (qb *QueryBuilder) SetOr(conditions []bsonV2.M) *QueryBuilder {
	qb.filter[OperatorOr] = conditions
	return qb
}

// SetAnd 设置多个条件的逻辑与
func (qb *QueryBuilder) SetAnd(conditions []bsonV2.M) *QueryBuilder {
	qb.filter[OperatorAnd] = conditions
	return qb
}

// SetNotEqual 设置字段的不等于条件
func (qb *QueryBuilder) SetNotEqual(field string, value interface{}) *QueryBuilder {
	qb.filter[field] = bsonV2.M{OperatorNe: value}
	return qb
}

// SetGreaterThan 设置字段的大于条件
func (qb *QueryBuilder) SetGreaterThan(field string, value interface{}) *QueryBuilder {
	qb.filter[field] = bsonV2.M{OperatorGt: value}
	return qb
}

// SetLessThan 设置字段的小于条件
func (qb *QueryBuilder) SetLessThan(field string, value interface{}) *QueryBuilder {
	qb.filter[field] = bsonV2.M{OperatorLt: value}
	return qb
}

// SetExists 设置字段是否存在条件
func (qb *QueryBuilder) SetExists(field string, exists bool) *QueryBuilder {
	qb.filter[field] = bsonV2.M{OperatorExists: exists}
	return qb
}

// SetType 设置字段的类型条件
func (qb *QueryBuilder) SetType(field string, typeValue interface{}) *QueryBuilder {
	qb.filter[field] = bsonV2.M{OperatorType: typeValue}
	return qb
}

// SetBetween 设置字段的范围查询条件
func (qb *QueryBuilder) SetBetween(field string, start, end interface{}) *QueryBuilder {
	qb.filter[field] = bsonV2.M{
		OperatorGte: start,
		OperatorLte: end,
	}
	return qb
}

// SetIn 设置字段的包含条件
func (qb *QueryBuilder) SetIn(field string, values []interface{}) *QueryBuilder {
	qb.filter[field] = bsonV2.M{OperatorIn: values}
	return qb
}

// SetNotIn 设置字段的排除条件
func (qb *QueryBuilder) SetNotIn(field string, values []interface{}) *QueryBuilder {
	qb.filter[field] = bsonV2.M{OperatorNin: values}
	return qb
}

// SetElemMatch 设置数组字段的匹配条件
func (qb *QueryBuilder) SetElemMatch(field string, match bsonV2.M) *QueryBuilder {
	qb.filter[field] = bsonV2.M{OperatorElemMatch: match}
	return qb
}

// SetAll 设置字段必须包含所有指定值的条件
func (qb *QueryBuilder) SetAll(field string, values []interface{}) *QueryBuilder {
	qb.filter[field] = bsonV2.M{OperatorAll: values}
	return qb
}

// SetSize 设置数组字段的大小条件
func (qb *QueryBuilder) SetSize(field string, size int) *QueryBuilder {
	qb.filter[field] = bsonV2.M{OperatorSize: size}
	return qb
}

// SetCurrentDate 设置字段为当前日期
func (qb *QueryBuilder) SetCurrentDate(field string) *QueryBuilder {
	qb.filter[field] = bsonV2.M{OperatorCurrentDate: true}
	return qb
}

// SetTextSearch 设置文本搜索条件
func (qb *QueryBuilder) SetTextSearch(search string) *QueryBuilder {
	qb.filter[OperatorText] = bsonV2.M{OperatorSearch: search}
	return qb
}

// SetMod 设置字段的模运算条件
func (qb *QueryBuilder) SetMod(field string, divisor, remainder int) *QueryBuilder {
	qb.filter[field] = bsonV2.M{OperatorMod: []int{divisor, remainder}}
	return qb
}

// SetRegex 设置正则表达式查询条件
func (qb *QueryBuilder) SetRegex(field string, pattern string) *QueryBuilder {
	qb.filter[field] = bsonV2.M{OperatorRegex: pattern}
	return qb
}

// SetGeoWithin 设置地理位置范围查询条件
func (qb *QueryBuilder) SetGeoWithin(field string, geometry bsonV2.M) *QueryBuilder {
	qb.filter[field] = bsonV2.M{OperatorGeoWithin: geometry}
	return qb
}

// SetGeoIntersects 设置地理位置相交查询条件
func (qb *QueryBuilder) SetGeoIntersects(field string, geometry bsonV2.M) *QueryBuilder {
	qb.filter[field] = bsonV2.M{OperatorGeoIntersects: geometry}
	return qb
}

// SetNear 设置地理位置附近查询条件
func (qb *QueryBuilder) SetNear(field string, point bsonV2.M, maxDistance, minDistance float64) *QueryBuilder {
	qb.filter[field] = bsonV2.M{
		OperatorNear: bsonV2.M{
			OperatorGeometry:    point,
			OperatorMaxDistance: maxDistance,
			OperatorMinDistance: minDistance,
		},
	}
	return qb
}

// SetNearSphere 设置球面距离附近查询条件
func (qb *QueryBuilder) SetNearSphere(field string, point bsonV2.M, maxDistance, minDistance float64) *QueryBuilder {
	qb.filter[field] = bsonV2.M{
		OperatorNearSphere: bsonV2.M{
			OperatorGeometry:    point,
			OperatorMaxDistance: maxDistance,
			OperatorMinDistance: minDistance,
		},
	}
	return qb
}

// SetLimit 设置查询结果的限制数量
func (qb *QueryBuilder) SetLimit(limit int64) *QueryBuilder {
	if qb.opts == nil {
		qb.opts = &optionsV2.FindOptions{}
	}
	qb.opts.Limit = &limit
	return qb
}

// SetSort 设置查询结果的排序条件
func (qb *QueryBuilder) SetSort(sort bsonV2.D) *QueryBuilder {
	if qb.opts == nil {
		qb.opts = &optionsV2.FindOptions{}
	}
	qb.opts.Sort = sort
	return qb
}

// SetSortWithPriority 设置查询结果的排序条件，并指定优先级
func (qb *QueryBuilder) SetSortWithPriority(sortFields []bsonV2.E) *QueryBuilder {
	if qb.opts == nil {
		qb.opts = &optionsV2.FindOptions{}
	}
	qb.opts.Sort = bsonV2.D(sortFields)
	return qb
}

// SetProjection 设置查询结果的字段投影
func (qb *QueryBuilder) SetProjection(projection bsonV2.M) *QueryBuilder {
	qb.opts.Projection = projection
	return qb
}

// SetSkip 设置查询结果的跳过数量
func (qb *QueryBuilder) SetSkip(skip int64) *QueryBuilder {
	qb.opts.Skip = &skip
	return qb
}

// SetPage 设置分页功能，page 从 1 开始，size 为每页的文档数量
func (qb *QueryBuilder) SetPage(page, size int64) *QueryBuilder {
	offset := (page - 1) * size
	qb.opts.Skip = &offset
	qb.opts.Limit = &size
	return qb
}

// AddStage 添加聚合阶段到管道
func (qb *QueryBuilder) AddStage(stage bsonV2.M) *QueryBuilder {
	qb.pipeline = append(qb.pipeline, stage)
	return qb
}

// BuildPipeline 返回最终的聚合管道
func (qb *QueryBuilder) BuildPipeline() []bsonV2.M {
	return qb.pipeline
}

// Build 返回最终的过滤条件和查询选项
func (qb *QueryBuilder) Build() (bsonV2.M, *optionsV2.FindOptions) {
	return qb.filter, qb.opts
}
