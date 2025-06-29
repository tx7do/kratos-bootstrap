package clickhouse

import "github.com/go-kratos/kratos/v2/errors"

var (
	// ErrInvalidColumnName is returned when an invalid column name is used.
	ErrInvalidColumnName = errors.InternalServer("INVALID_COLUMN_NAME", "invalid column name")

	// ErrInvalidTableName is returned when an invalid table name is used.
	ErrInvalidTableName = errors.InternalServer("INVALID_TABLE_NAME", "invalid table name")

	// ErrInvalidCondition is returned when an invalid condition is used in a query.
	ErrInvalidCondition = errors.InternalServer("INVALID_CONDITION", "invalid condition in query")

	// ErrQueryExecutionFailed is returned when a query execution fails.
	ErrQueryExecutionFailed = errors.InternalServer("QUERY_EXECUTION_FAILED", "query execution failed")

	// ErrExecutionFailed is returned when a general execution fails.
	ErrExecutionFailed = errors.InternalServer("EXECUTION_FAILED", "execution failed")

	// ErrAsyncInsertFailed is returned when an asynchronous insert operation fails.
	ErrAsyncInsertFailed = errors.InternalServer("ASYNC_INSERT_FAILED", "async insert operation failed")

	// ErrRowScanFailed is returned when scanning rows from a query result fails.
	ErrRowScanFailed = errors.InternalServer("ROW_SCAN_FAILED", "row scan failed")

	// ErrRowsIterationError is returned when there is an error iterating over rows.
	ErrRowsIterationError = errors.InternalServer("ROWS_ITERATION_ERROR", "rows iteration error")

	// ErrRowNotFound is returned when a specific row is not found in the result set.
	ErrRowNotFound = errors.InternalServer("ROW_NOT_FOUND", "row not found")

	// ErrConnectionFailed is returned when the connection to ClickHouse fails.
	ErrConnectionFailed = errors.InternalServer("CONNECTION_FAILED", "failed to connect to ClickHouse")

	// ErrDatabaseNotFound is returned when the specified database is not found.
	ErrDatabaseNotFound = errors.InternalServer("DATABASE_NOT_FOUND", "specified database not found")

	// ErrTableNotFound is returned when the specified table is not found.
	ErrTableNotFound = errors.InternalServer("TABLE_NOT_FOUND", "specified table not found")

	// ErrInsertFailed is returned when an insert operation fails.
	ErrInsertFailed = errors.InternalServer("INSERT_FAILED", "insert operation failed")

	// ErrUpdateFailed is returned when an update operation fails.
	ErrUpdateFailed = errors.InternalServer("UPDATE_FAILED", "update operation failed")

	// ErrDeleteFailed is returned when a delete operation fails.
	ErrDeleteFailed = errors.InternalServer("DELETE_FAILED", "delete operation failed")

	// ErrTransactionFailed is returned when a transaction fails.
	ErrTransactionFailed = errors.InternalServer("TRANSACTION_FAILED", "transaction failed")

	// ErrClientNotInitialized is returned when the ClickHouse client is not initialized.
	ErrClientNotInitialized = errors.InternalServer("CLIENT_NOT_INITIALIZED", "clickhouse client not initialized")

	// ErrGetServerVersionFailed is returned when getting the server version fails.
	ErrGetServerVersionFailed = errors.InternalServer("GET_SERVER_VERSION_FAILED", "failed to get server version")

	// ErrPingFailed is returned when a ping to the ClickHouse server fails.
	ErrPingFailed = errors.InternalServer("PING_FAILED", "ping to ClickHouse server failed")

	// ErrCreatorFunctionNil is returned when the creator function is nil.
	ErrCreatorFunctionNil = errors.InternalServer("CREATOR_FUNCTION_NIL", "creator function cannot be nil")

	// ErrBatchPrepareFailed is returned when a batch prepare operation fails.
	ErrBatchPrepareFailed = errors.InternalServer("BATCH_PREPARE_FAILED", "batch prepare operation failed")

	// ErrBatchSendFailed is returned when a batch send operation fails.
	ErrBatchSendFailed = errors.InternalServer("BATCH_SEND_FAILED", "batch send operation failed")

	// ErrBatchAppendFailed is returned when appending to a batch fails.
	ErrBatchAppendFailed = errors.InternalServer("BATCH_APPEND_FAILED", "batch append operation failed")

	// ErrBatchInsertFailed is returned when a batch insert operation fails.
	ErrBatchInsertFailed = errors.InternalServer("BATCH_INSERT_FAILED", "batch insert operation failed")

	// ErrInvalidDSN is returned when the data source name (DSN) is invalid.
	ErrInvalidDSN = errors.InternalServer("INVALID_DSN", "invalid data source name")

	// ErrInvalidProxyURL is returned when the proxy URL is invalid.
	ErrInvalidProxyURL = errors.InternalServer("INVALID_PROXY_URL", "invalid proxy URL")

	// ErrPrepareInsertDataFailed is returned when preparing insert data fails.
	ErrPrepareInsertDataFailed = errors.InternalServer("PREPARE_INSERT_DATA_FAILED", "failed to prepare insert data")

	// ErrInvalidColumnData is returned when the column data type is invalid.
	ErrInvalidColumnData = errors.InternalServer("INVALID_COLUMN_DATA", "invalid column data type")
)
