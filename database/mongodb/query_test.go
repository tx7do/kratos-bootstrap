package mongodb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	bsonV2 "go.mongodb.org/mongo-driver/v2/bson"
)

func TestQueryBuilder(t *testing.T) {
	// 创建 QueryBuilder 实例
	qb := NewQuery()

	// 测试 SetFilter
	filter := bsonV2.M{"name": "test"}
	qb.SetFilter(filter)
	assert.Equal(t, filter, qb.filter)

	// 测试 SetLimit
	limit := int64(10)
	qb.SetLimit(limit)
	assert.NotNil(t, qb.opts.Limit)
	assert.Equal(t, &limit, qb.opts.Limit)

	// 测试 SetSort
	sort := bsonV2.D{{Key: "name", Value: 1}}
	qb.SetSort(sort)
	assert.NotNil(t, qb.opts.Sort)
	assert.Equal(t, sort, qb.opts.Sort)

	// 测试 Build
	finalFilter, finalOpts := qb.Build()
	assert.Equal(t, filter, finalFilter)
	assert.Equal(t, qb.opts, finalOpts)
}

func TestQueryBuilderMethods(t *testing.T) {
	qb := NewQuery()

	// 测试 SetFilter
	filter := bsonV2.M{"name": "test"}
	qb.SetFilter(filter)
	assert.Equal(t, filter, qb.filter)

	// 测试 SetNotEqual
	qb.SetNotEqual("status", "inactive")
	assert.Equal(t, bsonV2.M{OperatorNe: "inactive"}, qb.filter["status"])

	// 测试 SetGreaterThan
	qb.SetGreaterThan("age", 18)
	assert.Equal(t, bsonV2.M{OperatorGt: 18}, qb.filter["age"])

	// 测试 SetLessThan
	qb.SetLessThan("age", 30)
	assert.Equal(t, bsonV2.M{OperatorLt: 30}, qb.filter["age"])

	// 测试 SetExists
	qb.SetExists("email", true)
	assert.Equal(t, bsonV2.M{OperatorExists: true}, qb.filter["email"])

	// 测试 SetType
	qb.SetType("age", "int")
	assert.Equal(t, bsonV2.M{OperatorType: "int"}, qb.filter["age"])

	// 测试 SetBetween
	qb.SetBetween("price", 10, 100)
	assert.Equal(t, bsonV2.M{OperatorGte: 10, OperatorLte: 100}, qb.filter["price"])

	// 测试 SetOr
	orConditions := []bsonV2.M{
		{"status": "active"},
		{"status": "pending"},
	}
	qb.SetOr(orConditions)
	assert.Equal(t, orConditions, qb.filter[OperatorOr])

	// 测试 SetAnd
	andConditions := []bsonV2.M{
		{"age": bsonV2.M{OperatorGt: 18}},
		{"status": "active"},
	}
	qb.SetAnd(andConditions)
	assert.Equal(t, andConditions, qb.filter[OperatorAnd])

	// 测试 SetLimit
	limit := int64(10)
	qb.SetLimit(limit)
	assert.NotNil(t, qb.opts.Limit)
	assert.Equal(t, &limit, qb.opts.Limit)

	// 测试 SetSort
	sort := bsonV2.D{{Key: "name", Value: 1}}
	qb.SetSort(sort)
	assert.NotNil(t, qb.opts.Sort)
	assert.Equal(t, sort, qb.opts.Sort)

	// 测试 SetSortWithPriority
	sortWithPriority := []bsonV2.E{{Key: "priority", Value: -1}, {Key: "name", Value: 1}}
	qb.SetSortWithPriority(sortWithPriority)
	assert.Equal(t, bsonV2.D(sortWithPriority), qb.opts.Sort)

	// 测试 SetProjection
	projection := bsonV2.M{"name": 1, "age": 1}
	qb.SetProjection(projection)
	assert.Equal(t, projection, qb.opts.Projection)

	// 测试 SetSkip
	skip := int64(5)
	qb.SetSkip(skip)
	assert.NotNil(t, qb.opts.Skip)
	assert.Equal(t, &skip, qb.opts.Skip)

	// 测试 SetPage
	page, size := int64(2), int64(10)
	qb.SetPage(page, size)
	assert.Equal(t, &size, qb.opts.Limit)
	assert.Equal(t, int64(10), *qb.opts.Limit)
	assert.Equal(t, int64(10), *qb.opts.Skip)

	// 测试 SetRegex
	qb.SetRegex("name", "^test")
	assert.Equal(t, bsonV2.M{OperatorRegex: "^test"}, qb.filter["name"])

	// 测试 SetIn
	qb.SetIn("tags", []interface{}{"tag1", "tag2"})
	assert.Equal(t, bsonV2.M{OperatorIn: []interface{}{"tag1", "tag2"}}, qb.filter["tags"])

	// 测试 Build
	finalFilter, finalOpts := qb.Build()
	assert.Equal(t, qb.filter, finalFilter)
	assert.Equal(t, qb.opts, finalOpts)
}

func TestSetGeoWithin(t *testing.T) {
	qb := NewQuery()

	field := "location"
	geometry := bsonV2.M{"type": "Polygon", "coordinates": []interface{}{
		[]interface{}{
			[]float64{40.0, -70.0},
			[]float64{41.0, -70.0},
			[]float64{41.0, -71.0},
			[]float64{40.0, -71.0},
			[]float64{40.0, -70.0},
		},
	}}

	qb.SetGeoWithin(field, geometry)

	expected := bsonV2.M{
		OperatorGeoWithin: bsonV2.M{
			OperatorGeometry: geometry,
		},
	}

	assert.Equal(t, expected, qb.filter[field])
}

func TestSetGeoIntersects(t *testing.T) {
	qb := NewQuery()

	field := "location"
	geometry := bsonV2.M{"type": "LineString", "coordinates": [][]float64{
		{40.0, -70.0},
		{41.0, -71.0},
	}}

	qb.SetGeoIntersects(field, geometry)

	expected := bsonV2.M{
		OperatorGeoIntersects: bsonV2.M{
			OperatorGeometry: geometry,
		},
	}

	assert.Equal(t, expected, qb.filter[field])
}

func TestSetNear(t *testing.T) {
	qb := NewQuery()

	field := "location"
	point := bsonV2.M{"type": "Point", "coordinates": []float64{40.7128, -74.0060}}
	maxDistance := 500.0
	minDistance := 50.0

	qb.SetNear(field, point, maxDistance, minDistance)

	expected := bsonV2.M{
		OperatorNear: bsonV2.M{
			OperatorGeometry:    point,
			OperatorMaxDistance: maxDistance,
			OperatorMinDistance: minDistance,
		},
	}

	assert.Equal(t, expected, qb.filter[field])
}

func TestSetNearSphere(t *testing.T) {
	qb := NewQuery()

	field := "location"
	point := bsonV2.M{"type": "Point", "coordinates": []float64{40.7128, -74.0060}}
	maxDistance := 1000.0
	minDistance := 100.0

	qb.SetNearSphere(field, point, maxDistance, minDistance)

	expected := bsonV2.M{
		OperatorNearSphere: bsonV2.M{
			OperatorGeometry:    point,
			OperatorMaxDistance: maxDistance,
			OperatorMinDistance: minDistance,
		},
	}

	assert.Equal(t, expected, qb.filter[field])
}

func TestQueryBuilderPipeline(t *testing.T) {
	// 创建 QueryBuilder 实例
	qb := NewQuery()

	// 添加聚合阶段
	matchStage := bsonV2.M{OperatorMatch: bsonV2.M{"status": "active"}}
	groupStage := bsonV2.M{OperatorGroup: bsonV2.M{"_id": "$category", "count": bsonV2.M{OperatorSum: 1}}}
	sortStage := bsonV2.M{OperatorSortAgg: bsonV2.M{"count": -1}}

	qb.AddStage(matchStage).AddStage(groupStage).AddStage(sortStage)

	// 构建 Pipeline
	pipeline := qb.BuildPipeline()

	// 验证 Pipeline
	expectedPipeline := []bsonV2.M{matchStage, groupStage, sortStage}
	assert.Equal(t, expectedPipeline, pipeline)
}
