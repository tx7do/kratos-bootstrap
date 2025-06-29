package influxdb

import (
	"reflect"
	"testing"
)

func TestBuildQuery(t *testing.T) {
	tests := []struct {
		name          string
		table         string
		filters       map[string]interface{}
		operators     map[string]string
		fields        []string
		expectedQuery string
		expectedArgs  []interface{}
	}{
		{
			name:          "Basic query with filters and fields",
			table:         "candles",
			filters:       map[string]interface{}{"s": "AAPL", "o": 150.0},
			fields:        []string{"s", "o", "h", "l", "c", "v"},
			expectedQuery: "SELECT s, o, h, l, c, v FROM candles WHERE s = ? AND o = ?",
			expectedArgs:  []interface{}{"AAPL", 150.0},
		},
		{
			name:          "Query with no filters",
			table:         "candles",
			filters:       map[string]interface{}{},
			fields:        []string{"s", "o", "h"},
			expectedQuery: "SELECT s, o, h FROM candles",
			expectedArgs:  []interface{}{},
		},
		{
			name:          "Query with no fields",
			table:         "candles",
			filters:       map[string]interface{}{"s": "AAPL"},
			fields:        []string{},
			expectedQuery: "SELECT * FROM candles WHERE s = ?",
			expectedArgs:  []interface{}{"AAPL"},
		},
		{
			name:          "Empty table name",
			table:         "",
			filters:       map[string]interface{}{"s": "AAPL"},
			fields:        []string{"s", "o"},
			expectedQuery: "SELECT s, o FROM  WHERE s = ?",
			expectedArgs:  []interface{}{"AAPL"},
		},
		{
			name:          "Special characters in filters",
			table:         "candles",
			filters:       map[string]interface{}{"name": "O'Reilly"},
			fields:        []string{"name"},
			expectedQuery: "SELECT name FROM candles WHERE name = ?",
			expectedArgs:  []interface{}{"O'Reilly"},
		},
		{
			name:          "Query with interval filters",
			table:         "candles",
			filters:       map[string]interface{}{"time": "now() - interval '15 minutes'"},
			fields:        []string{"*"},
			operators:     map[string]string{"time": ">="},
			expectedQuery: "SELECT * FROM candles WHERE time >= ?",
			expectedArgs:  []interface{}{"now() - interval '15 minutes'"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, args := BuildQuery(tt.table, tt.filters, tt.operators, tt.fields)

			if query != tt.expectedQuery {
				t.Errorf("expected query %s, got %s", tt.expectedQuery, query)
			}

			if !reflect.DeepEqual(args, tt.expectedArgs) {
				t.Errorf("expected args %v, got %v", tt.expectedArgs, args)
			}
		})
	}
}

func TestBuildQueryWithParams(t *testing.T) {
	tests := []struct {
		name          string
		table         string
		filters       map[string]interface{}
		operators     map[string]string
		fields        []string
		expectedQuery string
	}{
		{
			name:          "Basic query with filters and fields",
			table:         "candles",
			filters:       map[string]interface{}{"s": "'AAPL'", "o": 150.0},
			operators:     map[string]string{"o": ">"},
			fields:        []string{"s", "o", "h", "l", "c", "v"},
			expectedQuery: "SELECT s, o, h, l, c, v FROM candles WHERE s = 'AAPL' AND o > 150",
		},
		{
			name:          "Query with no filters",
			table:         "candles",
			filters:       map[string]interface{}{},
			operators:     map[string]string{},
			fields:        []string{"s", "o", "h"},
			expectedQuery: "SELECT s, o, h FROM candles",
		},
		{
			name:          "Query with no fields",
			table:         "candles",
			filters:       map[string]interface{}{"s": "'AAPL'"},
			operators:     map[string]string{},
			fields:        []string{},
			expectedQuery: "SELECT * FROM candles WHERE s = 'AAPL'",
		},
		{
			name:          "Empty table name",
			table:         "",
			filters:       map[string]interface{}{"s": "'AAPL'"},
			operators:     map[string]string{},
			fields:        []string{"s", "o"},
			expectedQuery: "SELECT s, o FROM  WHERE s = 'AAPL'",
		},
		{
			name:          "Special characters in filters",
			table:         "candles",
			filters:       map[string]interface{}{"name": "'O'Reilly'"},
			operators:     map[string]string{},
			fields:        []string{"name"},
			expectedQuery: "SELECT name FROM candles WHERE name = 'O'Reilly'",
		},
		{
			name:          "Query with interval filters",
			table:         "candles",
			filters:       map[string]interface{}{"time": "now() - interval '15 minutes'"},
			operators:     map[string]string{"time": ">="},
			fields:        []string{"*"},
			expectedQuery: "SELECT * FROM candles WHERE time >= now() - interval '15 minutes'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := BuildQueryWithParams(tt.table, tt.filters, tt.operators, tt.fields)

			if query != tt.expectedQuery {
				t.Errorf("expected query %s, got %s", tt.expectedQuery, query)
			}
			//t.Log(query)
		})
	}
}
