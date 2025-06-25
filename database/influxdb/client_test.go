package influxdb

import (
	"context"
	"testing"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
	"github.com/tx7do/go-utils/trans"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Candle struct {
	Symbol    *string
	Open      *float64
	High      *float64
	Low       *float64
	Close     *float64
	Volume    *float64
	StartTime *timestamppb.Timestamp
}

func (c *Candle) GetSymbol() string {
	if c.Symbol != nil {
		return *c.Symbol
	}
	return ""
}

func (c *Candle) GetOpen() float64 {
	if c.Open != nil {
		return *c.Open
	}
	return 0.0
}

func (c *Candle) GetHigh() float64 {
	if c.High != nil {
		return *c.High
	}
	return 0.0
}

func (c *Candle) GetLow() float64 {
	if c.Low != nil {
		return *c.Low
	}
	return 0.0
}

func (c *Candle) GetClose() float64 {
	if c.Close != nil {
		return *c.Close
	}
	return 0.0
}

func (c *Candle) GetVolume() float64 {
	if c.Volume != nil {
		return *c.Volume
	}
	return 0.0
}

func (c *Candle) GetStartTime() *timestamppb.Timestamp {
	if c.StartTime != nil {
		return c.StartTime
	}
	return timestamppb.Now()
}

type CandleMapper struct{}

var candleMapper CandleMapper

func (m *CandleMapper) ToPoint(data *Candle) *influxdb3.Point {
	p := influxdb3.NewPoint(
		"candles",
		map[string]string{"s": data.GetSymbol()},
		nil,
		data.StartTime.AsTime(),
	)

	p.
		SetDoubleField("o", data.GetOpen()).
		SetDoubleField("h", data.GetHigh()).
		SetDoubleField("l", data.GetLow()).
		SetDoubleField("c", data.GetClose()).
		SetDoubleField("v", data.GetVolume())

	return p
}

func (m *CandleMapper) ToData(point *influxdb3.Point) *Candle {
	symbol, _ := point.GetTag("s")

	return &Candle{
		Symbol:    &symbol,
		Open:      point.GetDoubleField("o"),
		High:      point.GetDoubleField("h"),
		Low:       point.GetDoubleField("l"),
		Close:     point.GetDoubleField("c"),
		Volume:    point.GetDoubleField("v"),
		StartTime: timestamppb.New(point.Values.Timestamp),
	}
}

func createTestClient() *Client {
	cli, _ := NewClient(
		log.DefaultLogger,
		&conf.Bootstrap{
			Data: &conf.Data{
				Influxdb: &conf.Data_InfluxDB{
					Host:         "http://localhost:8181",
					Token:        "apiv3_yYde4mJo0BYC7Ipi_00ZEex-A8if4swdqTBXiO-lCUDKhsIavHlRCQfo3p_DzI7S34ADHOC7Qxf600VVgW6LQQ",
					Database:     "finances",
					Organization: "primary",
				},
			},
		},
	)
	return cli
}

func TestNewClient(t *testing.T) {
	client := createTestClient()
	assert.NotNil(t, client)
}

func TestClient_Insert(t *testing.T) {
	client := createTestClient()
	assert.NotNil(t, client)

	item := &Candle{
		StartTime: timestamppb.New(time.Now()),
		Symbol:    trans.Ptr("AAPL"),
		Open:      trans.Ptr(1.0),
		High:      trans.Ptr(2.0),
		Low:       trans.Ptr(3.0),
		Close:     trans.Ptr(4.0),
		Volume:    trans.Ptr(1000.0),
	}

	point := candleMapper.ToPoint(item)

	err := client.Insert(context.Background(), point)
	assert.NoError(t, err)
}

func TestClient_BatchInsert(t *testing.T) {
	client := createTestClient()
	assert.NotNil(t, client)

	items := []*Candle{
		{
			StartTime: timestamppb.New(time.Now()),
			Symbol:    trans.Ptr("AAPL"),
			Open:      trans.Ptr(1.0),
			High:      trans.Ptr(2.0),
			Low:       trans.Ptr(3.0),
			Close:     trans.Ptr(4.0),
			Volume:    trans.Ptr(1000.0),
		},
	}

	var points []*influxdb3.Point
	for _, item := range items {
		point := candleMapper.ToPoint(item)
		points = append(points, point)
	}

	err := client.BatchInsert(
		context.Background(),
		points,
	)
	assert.NoError(t, err)
}

func TestClient_Query(t *testing.T) {
	client := createTestClient()
	assert.NotNil(t, client)

	ctx := context.Background()

	sql := `SELECT * FROM candles`

	iterator, err := client.Query(ctx, sql)
	assert.NoError(t, err)

	for iterator.Next() {
		point, _ := iterator.AsPoints().AsPoint()
		candle := candleMapper.ToData(point)
		t.Logf("[%v] Candle: %s, Open: %f, High: %f, Low: %f, Close: %f, Volume: %f\n",
			candle.GetStartTime().AsTime().String(),
			candle.GetSymbol(),
			candle.GetOpen(), candle.GetHigh(), candle.GetLow(), candle.GetClose(), candle.GetVolume(),
		)
	}

	candles, err := Query(ctx, client, sql, &candleMapper)
	assert.NoError(t, err)
	for _, candle := range candles {
		t.Logf("Candle: %s, Open: %f, High: %f, Low: %f, Close: %f, Volume: %f\n",
			candle.GetSymbol(),
			candle.GetOpen(), candle.GetHigh(), candle.GetLow(), candle.GetClose(), candle.GetVolume(),
		)
	}
}
