package consul

import "context"

// Option is consul config option.
type Option func(o *options)

type options struct {
	ctx  context.Context
	path string
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
