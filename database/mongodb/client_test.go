package mongodb

import (
	"context"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
	"github.com/tx7do/go-utils/trans"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Candle struct {
	Symbol    *string                `json:"s"`
	Open      *float64               `json:"o"`
	High      *float64               `json:"h"`
	Low       *float64               `json:"l"`
	Close     *float64               `json:"c"`
	Volume    *float64               `json:"v"`
	StartTime *timestamppb.Timestamp `json:"st"`
	EndTime   *timestamppb.Timestamp `json:"et"`
}

func createTestClient() *Client {
	cli, _ := NewClient(
		log.DefaultLogger,
		&conf.Bootstrap{
			Data: &conf.Data{
				Mongodb: &conf.Data_MongoDB{
					Uri:      "mongodb://root:123456@127.0.0.1:27017/?compressors=snappy,zlib,zstd",
					Database: trans.Ptr("finances"),
				},
			},
		},
	)
	return cli
}

func TestNewClient(t *testing.T) {
	client := createTestClient()
	assert.NotNil(t, client)

	client.CheckConnect()
}

func TestInsertOne(t *testing.T) {
	client := createTestClient()
	assert.NotNil(t, client)

	ctx := context.Background()

	candle := Candle{
		StartTime: timestamppb.New(time.Now()),
		Symbol:    trans.Ptr("AAPL"),
		Open:      trans.Ptr(1.0),
		High:      trans.Ptr(2.0),
		Low:       trans.Ptr(3.0),
		Close:     trans.Ptr(4.0),
		Volume:    trans.Ptr(1000.0),
	}

	_, err := client.InsertOne(ctx, "candles", candle)
	assert.NoError(t, err)
}
