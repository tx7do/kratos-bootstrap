package s3

import "errors"

var (
	ErrNilClient      = errors.New("s3 client is nil")
	ErrEmptyBucket    = errors.New("s3 bucket is empty")
	ErrEmptyObjectKey = errors.New("s3 object key is empty")
	ErrNilObjectBody  = errors.New("s3 object body is nil")
)
