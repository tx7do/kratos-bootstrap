package clickhouse

import (
	"database/sql"
	"reflect"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	timeFormat = "2006-01-02 15:04:05.000000000"
)

func structToValueArray(input any) []any {
	// 检查是否是指针类型，如果是则解引用
	val := reflect.ValueOf(input)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// 确保输入是结构体
	if val.Kind() != reflect.Struct {
		return nil
	}

	var values []any
	for i := 0; i < val.NumField(); i++ {
		value := val.Field(i).Interface()

		switch v := value.(type) {
		case *sql.NullString:
			if v.Valid {
				values = append(values, v.String)
			} else {
				values = append(values, nil)
			}
		case *sql.NullInt64:
			if v.Valid {
				values = append(values, v.Int64)
			} else {
				values = append(values, nil)
			}
		case *sql.NullFloat64:
			if v.Valid {
				values = append(values, v.Float64)
			} else {
				values = append(values, nil)
			}
		case *sql.NullBool:
			if v.Valid {
				values = append(values, v.Bool)
			} else {
				values = append(values, nil)
			}

		case *sql.NullTime:
			if v != nil && v.Valid {
				values = append(values, v.Time.Format(timeFormat))
			} else {
				values = append(values, nil)
			}

		case *time.Time:
			if v != nil {
				values = append(values, v.Format(timeFormat))
			} else {
				values = append(values, nil)
			}

		case time.Time:
			// 处理 time.Time 类型
			if !v.IsZero() {
				values = append(values, v.Format(timeFormat))
			} else {
				values = append(values, nil) // 如果时间为零值，插入 NULL
			}

		case timestamppb.Timestamp:
			// 处理 timestamppb.Timestamp 类型
			if !v.IsValid() {
				values = append(values, v.AsTime().Format(timeFormat))
			} else {
				values = append(values, nil) // 如果时间为零值，插入 NULL
			}

		case *timestamppb.Timestamp:
			// 处理 *timestamppb.Timestamp 类型
			if v != nil && v.IsValid() {
				values = append(values, v.AsTime().Format(timeFormat))
			} else {
				values = append(values, nil) // 如果时间为零值，插入 NULL
			}

		case durationpb.Duration:
			// 处理 timestamppb.Duration 类型
			if v.AsDuration() != 0 {
				values = append(values, v.AsDuration().String())
			} else {
				values = append(values, nil) // 如果时间为零值，插入 NULL
			}

		case *durationpb.Duration:
			// 处理 *timestamppb.Duration 类型
			if v != nil && v.AsDuration() != 0 {
				values = append(values, v.AsDuration().String())
			} else {
				values = append(values, nil) // 如果时间为零值，插入 NULL
			}

		case []any:
			// 处理切片类型
			if len(v) > 0 {
				for _, item := range v {
					if item == nil {
						values = append(values, nil)
					} else {
						values = append(values, item)
					}
				}
			} else {
				values = append(values, nil) // 如果切片为空，插入 NULL
			}

		case [][]any:
			// 处理二维切片类型
			if len(v) > 0 {
				for _, item := range v {
					if len(item) > 0 {
						values = append(values, item)
					} else {
						values = append(values, nil) // 如果子切片为空，插入 NULL
					}
				}
			} else {
				values = append(values, nil) // 如果二维切片为空，插入 NULL
			}

		default:
			values = append(values, v)
		}
	}

	return values
}
