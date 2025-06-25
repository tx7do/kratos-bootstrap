package elasticsearch

import "github.com/go-kratos/kratos/v2/errors"

var (
	// ErrRequestFailed is returned when a request to Elasticsearch fails.
	ErrRequestFailed = errors.InternalServer("REQUEST_FAILED", "request failed")

	// ErrIndexNotFound is returned when the specified index does not exist.
	ErrIndexNotFound = errors.InternalServer("INDEX_NOT_FOUND", "index not found")

	// ErrIndexAlreadyExists is returned when trying to create an index that already exists.
	ErrIndexAlreadyExists = errors.InternalServer("INDEX_ALREADY_EXISTS", "index already exists")

	ErrCreateIndex = errors.InternalServer("CREATE_INDEX_FAILED", "failed to create index")

	ErrDeleteIndex = errors.InternalServer("DELETE_INDEX_FAILED", "failed to delete index")

	// ErrDocumentNotFound is returned when a document is not found in the index.
	ErrDocumentNotFound = errors.InternalServer("DOCUMENT_NOT_FOUND", "document not found")

	// ErrDocumentAlreadyExists is returned when trying to create a document that already exists.
	ErrDocumentAlreadyExists = errors.InternalServer("DOCUMENT_ALREADY_EXISTS", "document already exists")

	// ErrInvalidQuery is returned when the query provided to Elasticsearch is invalid.
	ErrInvalidQuery = errors.InternalServer("INVALID_QUERY", "invalid query")

	// ErrUnmarshalResponse is returned when the response from Elasticsearch cannot be unmarshalled.
	ErrUnmarshalResponse = errors.InternalServer("UNMARSHAL_RESPONSE_FAILED", "failed to unmarshal response")

	ErrInsertDocument = errors.InternalServer("INSERT_DOCUMENT_FAILED", "failed to insert document")

	ErrBatchInsertDocument = errors.InternalServer("BATCH_INSERT_DOCUMENT_FAILED", "failed to batch insert documents")

	ErrGetDocument = errors.InternalServer("GET_DOCUMENT_FAILED", "failed to get document")

	ErrSearchDocument = errors.InternalServer("SEARCH_DOCUMENT_FAILED", "failed to search document")
)
