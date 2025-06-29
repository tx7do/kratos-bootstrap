package clickhouse

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	clickhouseV2 "github.com/ClickHouse/clickhouse-go/v2"
	driverV2 "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

// BatchInserter 批量插入器
type BatchInserter struct {
	conn       clickhouseV2.Conn
	tableName  string
	columns    []string
	batchSize  int
	rows       []interface{}
	insertStmt string
	mu         sync.Mutex
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewBatchInserter 创建新的批量插入器
func NewBatchInserter(
	ctx context.Context,
	conn clickhouseV2.Conn,
	tableName string,
	batchSize int,
	columns []string,
) (*BatchInserter, error) {
	if batchSize <= 0 {
		batchSize = 1000 // 默认批量大小
	}

	if len(columns) == 0 {
		return nil, errors.New("必须指定列名")
	}

	// 构建INSERT语句
	placeholders := make([]string, len(columns))
	for i := range placeholders {
		placeholders[i] = "?"
	}

	insertStmt := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	ctx, cancel := context.WithCancel(ctx)

	return &BatchInserter{
		conn:       conn,
		tableName:  tableName,
		columns:    columns,
		batchSize:  batchSize,
		rows:       make([]interface{}, 0, batchSize),
		insertStmt: insertStmt,
		ctx:        ctx,
		cancel:     cancel,
	}, nil
}

// Add 添加数据行
func (bi *BatchInserter) Add(row interface{}) error {
	bi.mu.Lock()
	defer bi.mu.Unlock()

	// 检查上下文是否已取消
	if bi.ctx.Err() != nil {
		return bi.ctx.Err()
	}

	bi.rows = append(bi.rows, row)

	// 达到批量大小时自动提交
	if len(bi.rows) >= bi.batchSize {
		return bi.flush()
	}

	return nil
}

// Flush 强制提交当前批次
func (bi *BatchInserter) Flush() error {
	bi.mu.Lock()
	defer bi.mu.Unlock()

	return bi.flush()
}

// Close 关闭插入器并提交剩余数据
func (bi *BatchInserter) Close() error {
	defer bi.cancel()

	bi.mu.Lock()
	defer bi.mu.Unlock()

	return bi.flush()
}

// flush 内部提交方法
func (bi *BatchInserter) flush() error {
	if len(bi.rows) == 0 {
		return nil
	}

	// 创建批量
	batch, err := bi.conn.PrepareBatch(bi.ctx, bi.insertStmt)
	if err != nil {
		return ErrBatchPrepareFailed
	}

	// 添加所有行
	for _, row := range bi.rows {
		// 使用反射获取字段值
		if err = appendStructToBatch(batch, row, bi.columns); err != nil {
			return ErrBatchAppendFailed
		}
	}

	// 提交批量
	if err = batch.Send(); err != nil {
		return ErrBatchSendFailed
	}

	// 清空批次
	bi.rows = bi.rows[:0]
	return nil
}

// appendStructToBatch 使用反射将结构体字段添加到批次
func appendStructToBatch(batch driverV2.Batch, obj interface{}, columns []string) error {
	v := reflect.ValueOf(obj)

	// 如果是指针，获取指针指向的值
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return errors.New("nil指针")
		}
		v = v.Elem()
	}

	// 必须是结构体
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("期望结构体类型，得到 %v", v.Kind())
	}

	// 获取结构体类型
	t := v.Type()

	// 准备参数值
	values := make([]interface{}, len(columns))

	// 映射列名到结构体字段
	for i, col := range columns {
		// 查找匹配的字段
		found := false
		for j := 0; j < v.NumField(); j++ {
			field := t.Field(j)

			// 检查ch标签
			if tag := field.Tag.Get("ch"); strings.TrimSpace(tag) == col {
				values[i] = v.Field(j).Interface()
				found = true
				break
			}

			// 检查json标签
			jsonTags := strings.Split(field.Tag.Get("json"), ",")
			if len(jsonTags) > 0 && strings.TrimSpace(jsonTags[0]) == col {
				values[i] = v.Field(j).Interface()
				found = true
				break
			}

			// 检查字段名
			if field.Name == col {
				values[i] = v.Field(j).Interface()
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("未找到列 %s 对应的结构体字段", col)
		}
	}

	// 添加到批次
	return batch.Append(values...)
}
