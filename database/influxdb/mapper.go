package influxdb

import (
	"context"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
)

// Mapper 数据转换的接口
type Mapper[T any] interface {
	// ToPoint 将数据转换为InfluxDB的Point格式
	ToPoint(data *T) *influxdb3.Point

	// ToData 将InfluxDB的Point转换为原始数据
	ToData(point *influxdb3.Point) *T
}

// Insert 插入数据
func Insert[T any](ctx context.Context, c *Client, data *T, mapper Mapper[T]) error {
	if c.cli == nil {
		return ErrClientNotConnected
	}

	if data == nil {
		return ErrEmptyData
	}

	point := mapper.ToPoint(data)
	if point == nil {
		return ErrInvalidPoint
	}

	err := c.Insert(ctx, point)
	if err != nil {
		return err
	}

	return nil
}

// BatchInsert 批量插入数据
func BatchInsert[T any](ctx context.Context, c *Client, data []*T, mapper Mapper[T]) error {
	if c.cli == nil {
		return ErrClientNotConnected
	}

	if len(data) == 0 {
		return ErrEmptyData
	}

	points := make([]*influxdb3.Point, len(data))
	for i, d := range data {
		point := mapper.ToPoint(d)
		if point == nil {
			return ErrInvalidPoint
		}
		points[i] = point
	}

	err := c.BatchInsert(ctx, points)
	if err != nil {
		return err
	}

	return nil
}

// Query 查询数据
func Query[T any](ctx context.Context, c *Client, query string, mapper Mapper[T]) ([]*T, error) {
	if c.cli == nil {
		return nil, ErrClientNotConnected
	}

	iterator, err := c.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var dataset []*T

	for iterator.Next() {
		point, _ := iterator.AsPoints().AsPoint()
		if point == nil {
			return nil, ErrInvalidPoint
		}

		data := mapper.ToData(point)
		dataset = append(dataset, data)
	}

	if iterator.Err() != nil {
		return nil, iterator.Err()
	}

	return dataset, nil
}
