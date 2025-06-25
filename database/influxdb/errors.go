package influxdb

import "github.com/go-kratos/kratos/v2/errors"

var (
	ErrInfluxDBClientNotInitialized = errors.InternalServer("INFLUXDB_CLIENT_NOT_INITIALIZED", "client not initialized")

	ErrInfluxDBConnectFailed = errors.InternalServer("INFLUXDB_CONNECT_FAILED", "connect failed")

	ErrInfluxDBCreateDatabaseFailed = errors.InternalServer("INFLUXDB_CREATE_DATABASE_FAILED", "database create failed")

	ErrInfluxDBQueryFailed = errors.InternalServer("INFLUXDB_QUERY_FAILED", "query failed")

	ErrClientNotConnected = errors.InternalServer("INFLUXDB_CLIENT_NOT_CONNECTED", "client not connected")

	ErrInvalidPoint = errors.InternalServer("INFLUXDB_INVALID_POINT", "invalid point")

	ErrNoPointsToInsert = errors.InternalServer("INFLUXDB_NO_POINTS_TO_INSERT", "no points to insert")

	ErrEmptyData = errors.InternalServer("INFLUXDB_EMPTY_DATA", "empty data")

	ErrBatchInsertFailed = errors.InternalServer("INFLUXDB_BATCH_INSERT_FAILED", "batch insert failed")

	ErrInsertFailed = errors.InternalServer("INFLUXDB_INSERT_FAILED", "insert failed")
)
