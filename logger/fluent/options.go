package fluent

import "time"

// Option is fluentd logger option.
type Option func(*options)

type options struct {
	timeout            time.Duration
	writeTimeout       time.Duration
	bufferLimit        int
	retryWait          int
	maxRetry           int
	maxRetryWait       int
	tagPrefix          string
	async              bool
	forceStopAsyncSend bool
}

// WithTimeout with config Timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(opts *options) {
		opts.timeout = timeout
	}
}

// WithWriteTimeout with config WriteTimeout.
func WithWriteTimeout(writeTimeout time.Duration) Option {
	return func(opts *options) {
		opts.writeTimeout = writeTimeout
	}
}

// WithBufferLimit with config BufferLimit.
func WithBufferLimit(bufferLimit int) Option {
	return func(opts *options) {
		opts.bufferLimit = bufferLimit
	}
}

// WithRetryWait with config RetryWait.
func WithRetryWait(retryWait int) Option {
	return func(opts *options) {
		opts.retryWait = retryWait
	}
}

// WithMaxRetry with config MaxRetry.
func WithMaxRetry(maxRetry int) Option {
	return func(opts *options) {
		opts.maxRetry = maxRetry
	}
}

// WithMaxRetryWait with config MaxRetryWait.
func WithMaxRetryWait(maxRetryWait int) Option {
	return func(opts *options) {
		opts.maxRetryWait = maxRetryWait
	}
}

// WithTagPrefix with config TagPrefix.
func WithTagPrefix(tagPrefix string) Option {
	return func(opts *options) {
		opts.tagPrefix = tagPrefix
	}
}

// WithAsync with config Async.
func WithAsync(async bool) Option {
	return func(opts *options) {
		opts.async = async
	}
}

// WithForceStopAsyncSend with config ForceStopAsyncSend.
func WithForceStopAsyncSend(forceStopAsyncSend bool) Option {
	return func(opts *options) {
		opts.forceStopAsyncSend = forceStopAsyncSend
	}
}
