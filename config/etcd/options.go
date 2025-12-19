package etcd

import "context"

// Option is etcd config option.
type Option func(o *options)

type options struct {
	ctx    context.Context
	path   string
	prefix bool
}

// WithContext with registry context.
func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

// WithPath is config path
func WithPath(p string) Option {
	return func(o *options) {
		o.path = p
	}
}

// WithPrefix is config prefix
func WithPrefix(prefix bool) Option {
	return func(o *options) {
		o.prefix = prefix
	}
}
