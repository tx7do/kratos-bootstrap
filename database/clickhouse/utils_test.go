package clickhouse

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Ptr returns a pointer to the provided value.
func Ptr[T any](v T) *T {
	return &v
}

func TestStructToValueArrayWithCandle(t *testing.T) {
	now := time.Now()

	candle := Candle{
		Timestamp: Ptr(now),
		Symbol:    Ptr("AAPL"),
		Open:      Ptr(100.5),
		High:      Ptr(105.0),
		Low:       Ptr(99.5),
		Close:     Ptr(102.0),
		Volume:    Ptr(1500.0),
	}

	values := structToValueArray(candle)
	assert.NotNil(t, values, "Values should not be nil")
	assert.Len(t, values, 7, "Expected 7 fields in the Candle struct")
	assert.Equal(t, now.Format(timeFormat), values[0].(string), "Timestamp should match")
	assert.Equal(t, *candle.Symbol, *(values[1].(*string)), "Symbol should match")
	assert.Equal(t, *candle.Open, *values[2].(*float64), "Open price should match")
	assert.Equal(t, *candle.High, *values[3].(*float64), "High price should match")
	assert.Equal(t, *candle.Low, *values[4].(*float64), "Low price should match")
	assert.Equal(t, *candle.Close, *values[5].(*float64), "Close price should match")
	assert.Equal(t, *candle.Volume, *values[6].(*float64), "Volume should match")

	t.Logf("QueryRow Result: [%v] Candle: %s, Open: %f, High: %f, Low: %f, Close: %f, Volume: %f\n",
		values[0],
		*(values[1].(*string)),
		*values[2].(*float64),
		*values[3].(*float64),
		*values[4].(*float64),
		*values[5].(*float64),
		*values[6].(*float64),
	)
}

func TestStructToValueArrayWithCandlePtr(t *testing.T) {
	now := time.Now()

	candle := &Candle{
		Timestamp: Ptr(now),
		Symbol:    Ptr("AAPL"),
		Open:      Ptr(100.5),
		High:      Ptr(105.0),
		Low:       Ptr(99.5),
		Close:     Ptr(102.0),
		Volume:    Ptr(1500.0),
	}

	values := structToValueArray(candle)
	assert.NotNil(t, values, "Values should not be nil")
	assert.Len(t, values, 7, "Expected 7 fields in the Candle struct")
	assert.Equal(t, now.Format(timeFormat), values[0].(string), "Timestamp should match")
	assert.Equal(t, *candle.Symbol, *(values[1].(*string)), "Symbol should match")
	assert.Equal(t, *candle.Open, *values[2].(*float64), "Open price should match")
	assert.Equal(t, *candle.High, *values[3].(*float64), "High price should match")
	assert.Equal(t, *candle.Low, *values[4].(*float64), "Low price should match")
	assert.Equal(t, *candle.Close, *values[5].(*float64), "Close price should match")
	assert.Equal(t, *candle.Volume, *values[6].(*float64), "Volume should match")

	t.Logf("QueryRow Result: [%v] Candle: %s, Open: %f, High: %f, Low: %f, Close: %f, Volume: %f\n",
		values[0],
		*(values[1].(*string)),
		*values[2].(*float64),
		*values[3].(*float64),
		*values[4].(*float64),
		*values[5].(*float64),
		*values[6].(*float64),
	)
}
