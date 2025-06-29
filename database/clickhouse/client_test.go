package clickhouse

import (
	"context"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

type Candle struct {
	Symbol    string    `json:"symbol" ch:"symbol"`
	Open      float64   `json:"open" ch:"open"`
	High      float64   `json:"high" ch:"high"`
	Low       float64   `json:"low" ch:"low"`
	Close     float64   `json:"close" ch:"close"`
	Volume    float64   `json:"volume" ch:"volume"`
	Timestamp time.Time `json:"timestamp" ch:"timestamp"`
}

func createTestClient() *Client {
	database := "finances"
	username := "default"
	password := "*Abcd123456"
	cli, _ := NewClient(
		log.DefaultLogger,
		&conf.Bootstrap{
			Data: &conf.Data{
				Clickhouse: &conf.Data_ClickHouse{
					Addresses: []string{"localhost:9000"},
					Database:  &database,
					Username:  &username,
					Password:  &password,
				},
			},
		},
	)
	return cli
}

func createCandlesTable(client *Client) {
	// 创建表的 SQL 语句
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS candles (
			timestamp DateTime,
			symbol String,
			open Float64,
			high Float64,
			low Float64,
			close Float64,
			volume Float64
		) ENGINE = MergeTree()
		ORDER BY timestamp
	`
	err := client.Exec(context.Background(), createTableQuery)
	if err != nil {
		log.Errorf("Failed to create candles table: %v", err)
		return
	}
}

func TestNewClient(t *testing.T) {
	client := createTestClient()
	assert.NotNil(t, client)

	// 测试 CheckConnection
	err := client.CheckConnection(context.Background())
	assert.NoError(t, err, "CheckConnection 应该成功执行")

	// 测试 GetServerVersion
	version := client.GetServerVersion()
	assert.NotEmpty(t, version, "GetServerVersion 应该返回非空值")

	createCandlesTable(client)
}

func TestAsyncInsert(t *testing.T) {
	client := createTestClient()
	assert.NotNil(t, client)

	// 测试异步插入
	err := client.AsyncInsert(context.Background(), "INSERT INTO test_table (id, name) VALUES (?, ?)", true, 1, "example")
	assert.NoError(t, err, "AsyncInsert 应该成功执行")
}

func TestBatchInsert(t *testing.T) {
	client := createTestClient()
	assert.NotNil(t, client)

	// 测试数据
	data := [][]interface{}{
		{1, "example1"},
		{2, "example2"},
		{3, "example3"},
	}

	// 测试批量插入
	err := client.BatchInsert(context.Background(), "INSERT INTO test_table (id, name) VALUES (?, ?)", data)
	assert.NoError(t, err, "BatchInsert 应该成功执行")
}

func TestInsertIntoCandlesTable(t *testing.T) {
	client := createTestClient()
	assert.NotNil(t, client)

	createCandlesTable(client)

	// 插入数据的 SQL 语句
	insertQuery := `
		INSERT INTO candles (timestamp, symbol, open, high, low, close, volume)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	// 测试数据
	err := client.AsyncInsert(context.Background(), insertQuery, true,
		"2023-10-01 12:00:00", "AAPL", 100.5, 105.0, 99.5, 102.0, 1500.0)
	assert.NoError(t, err, "InsertIntoCandlesTable 应该成功执行")
}

func TestQueryCandlesTable(t *testing.T) {
	client := createTestClient()
	assert.NotNil(t, client)

	createCandlesTable(client)

	// 查询数据的 SQL 语句
	query := `
		SELECT timestamp, symbol, open, high, low, close, volume
		FROM candles
	`

	// 定义结果集
	var results []any

	// 执行查询
	err := client.Query(context.Background(), func() interface{} { return &Candle{} }, &results, query)
	assert.NoError(t, err, "QueryCandlesTable 应该成功执行")
	assert.NotEmpty(t, results, "QueryCandlesTable 应该返回结果")
}

func TestSelectCandlesTable(t *testing.T) {
	client := createTestClient()
	assert.NotNil(t, client)

	createCandlesTable(client)

	// 查询数据的 SQL 语句
	query := `
		SELECT timestamp, symbol, open, high, low, close, volume
		FROM candles
	`

	// 定义结果集
	var results []Candle

	// 执行查询
	err := client.Select(context.Background(), &results, query)
	assert.NoError(t, err, "QueryCandlesTable 应该成功执行")
	assert.NotEmpty(t, results, "QueryCandlesTable 应该返回结果")
}

func TestQueryRow(t *testing.T) {
	client := createTestClient()
	assert.NotNil(t, client)

	createCandlesTable(client)

	// 插入测试数据
	insertQuery := `
		INSERT INTO candles (timestamp, symbol, open, high, low, close, volume)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	err := client.AsyncInsert(context.Background(), insertQuery, true,
		"2023-10-01 12:00:00", "AAPL", 100.5, 105.0, 99.5, 102.0, 1500.0)
	assert.NoError(t, err, "数据插入失败")

	// 查询单行数据
	query := `
		SELECT timestamp, symbol, open, high, low, close, volume
		FROM candles
		WHERE symbol = ?
	`
	var result Candle

	err = client.QueryRow(context.Background(), &result, query, "AAPL")
	assert.NoError(t, err, "QueryRow 应该成功执行")
	assert.Equal(t, "AAPL", result.Symbol, "symbol 列值应该为 AAPL")
	assert.Equal(t, 100.5, result.Open, "open 列值应该为 100.5")
	assert.Equal(t, 1500.0, result.Volume, "volume 列值应该为 1500.0")
}

func TestDropCandlesTable(t *testing.T) {
	client := createTestClient()
	assert.NotNil(t, client)

	// 删除表的 SQL 语句
	dropTableQuery := `DROP TABLE IF EXISTS candles`

	// 执行删除表操作
	err := client.Exec(context.Background(), dropTableQuery)
	assert.NoError(t, err, "DropCandlesTable 应该成功执行")
}

func TestAggregateCandlesTable(t *testing.T) {
	client := createTestClient()
	assert.NotNil(t, client)

	createCandlesTable(client)

	// 聚合查询的 SQL 语句
	query := `
		SELECT symbol, 
		       MAX(high) AS max_high, 
		       MIN(low) AS min_low, 
		       AVG(close) AS avg_close, 
		       SUM(volume) AS total_volume
		FROM candles
		GROUP BY symbol
	`

	// 定义结果集
	var results []struct {
		Symbol      string  `ch:"symbol"`
		MaxHigh     float64 `ch:"max_high"`
		MinLow      float64 `ch:"min_low"`
		AvgClose    float64 `ch:"avg_close"`
		TotalVolume float64 `ch:"total_volume"`
	}

	// 执行查询
	err := client.Select(context.Background(), &results, query)
	assert.NoError(t, err, "AggregateCandlesTable 应该成功执行")
	assert.NotEmpty(t, results, "AggregateCandlesTable 应该返回结果")
}

func TestBatchInsertCandlesTable(t *testing.T) {
	client := createTestClient()
	assert.NotNil(t, client)

	createCandlesTable(client)

	// 测试数据
	data := [][]interface{}{
		{"2023-10-01 12:00:00", "AAPL", 100.5, 105.0, 99.5, 102.0, 1500.0},
		{"2023-10-01 12:01:00", "GOOG", 200.5, 205.0, 199.5, 202.0, 2500.0},
		{"2023-10-01 12:02:00", "MSFT", 300.5, 305.0, 299.5, 302.0, 3500.0},
	}

	// 批量插入数据
	err := client.BatchInsert(context.Background(), `
		INSERT INTO candles (timestamp, symbol, open, high, low, close, volume)
		VALUES (?, ?, ?, ?, ?, ?, ?)`, data)
	assert.NoError(t, err, "BatchInsertCandlesTable 应该成功执行")
}
