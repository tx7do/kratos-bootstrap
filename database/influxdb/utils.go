package influxdb

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func BuildQuery(
	table string,
	filters map[string]interface{},
	operators map[string]string,
	fields []string,
) (string, []interface{}) {
	var queryBuilder strings.Builder
	args := make([]interface{}, 0)

	// 构建 SELECT 语句
	queryBuilder.WriteString("SELECT ")
	if len(fields) > 0 {
		queryBuilder.WriteString(strings.Join(fields, ", "))
	} else {
		queryBuilder.WriteString("*")
	}
	queryBuilder.WriteString(fmt.Sprintf(" FROM %s", table))

	// 构建 WHERE 条件
	if len(filters) > 0 {
		queryBuilder.WriteString(" WHERE ")
		var conditions []string
		var operator string
		for key, value := range filters {
			operator = "=" // 默认操作符
			if op, exists := operators[key]; exists {
				operator = op
			}
			conditions = append(conditions, fmt.Sprintf("%s %s ?", key, operator))
			args = append(args, value)
		}
		queryBuilder.WriteString(strings.Join(conditions, " AND "))
	}

	return queryBuilder.String(), args
}

func GetPointTag(point *influxdb3.Point, name string) *string {
	if point == nil {
		return nil
	}
	tagValue, ok := point.GetTag(name)
	if !ok || tagValue == "" {
		return nil
	}
	return &tagValue
}

func GetBoolPointTag(point *influxdb3.Point, name string) *bool {
	if point == nil {
		return nil
	}
	tagValue, ok := point.GetTag(name)
	if !ok || tagValue == "" {
		return nil
	}

	value := tagValue == "true"
	return &value
}

func GetUint32PointTag(point *influxdb3.Point, name string) *uint32 {
	if point == nil {
		return nil
	}
	tagValue, ok := point.GetTag(name)
	if !ok || tagValue == "" {
		return nil
	}

	value, err := strconv.ParseUint(tagValue, 10, 64)
	if err != nil {
		return nil
	}
	value32 := uint32(value)
	return &value32
}

func GetUint64PointTag(point *influxdb3.Point, name string) *uint64 {
	if point == nil {
		return nil
	}
	tagValue, ok := point.GetTag(name)
	if !ok || tagValue == "" {
		return nil
	}

	value, err := strconv.ParseUint(tagValue, 10, 64)
	if err != nil {
		return nil
	}
	return &value
}

func GetEnumPointTag[T ~int32](point *influxdb3.Point, name string, valueMap map[string]int32) *T {
	if point == nil {
		return nil
	}
	tagValue, ok := point.GetTag(name)
	if !ok || tagValue == "" {
		return nil
	}
	enumValue, exists := valueMap[tagValue]
	if !exists {
		return nil
	}

	enumType := T(enumValue)
	return &enumType
}

func GetTimestampField(point *influxdb3.Point, name string) *timestamppb.Timestamp {
	if point == nil {
		return nil
	}

	value := point.GetField(name)
	if value == nil {
		return nil
	}
	if timestamp, ok := value.(*timestamppb.Timestamp); ok {
		return timestamp
	}
	if timeValue, ok := value.(time.Time); ok {
		return timestamppb.New(timeValue)
	}
	return nil
}

func GetUint32Field(point *influxdb3.Point, name string) *uint32 {
	if point == nil {
		return nil
	}

	value := point.GetUIntegerField(name)
	if value == nil {
		return nil
	}
	uint32Value := uint32(*value)
	if uint32Value == 0 {
		return nil
	}
	return &uint32Value
}

func BoolToString(value *bool) string {
	if value == nil {
		return "false"
	}
	if *value {
		return "true"
	}
	return "false"
}

func Uint64ToString(value *uint64) string {
	if value == nil {
		return "0"
	}
	return fmt.Sprintf("%d", *value)
}

func BuildQueryWithParams(
	table string,
	filters map[string]interface{},
	operators map[string]string,
	fields []string,
) string {
	var queryBuilder strings.Builder

	// 构建 SELECT 语句
	queryBuilder.WriteString("SELECT ")
	if len(fields) > 0 {
		queryBuilder.WriteString(strings.Join(fields, ", "))
	} else {
		queryBuilder.WriteString("*")
	}
	queryBuilder.WriteString(fmt.Sprintf(" FROM %s", table))

	// 构建 WHERE 条件
	if len(filters) > 0 {
		var operator string
		queryBuilder.WriteString(" WHERE ")
		var conditions []string
		for key, value := range filters {
			operator = "=" // 默认操作符
			if op, exists := operators[key]; exists {
				operator = op
			}
			conditions = append(conditions, fmt.Sprintf("%s %s %v", key, operator, value))
		}
		queryBuilder.WriteString(strings.Join(conditions, " AND "))
	}

	return queryBuilder.String()
}
