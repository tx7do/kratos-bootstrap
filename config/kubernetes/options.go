package kubernetes

import "time"

// Option is kubernetes option.
type Option func(*options)

type options struct {
	// kubernetes namespace
	Namespace string
	// kubernetes labelSelector example `app=test`
	LabelSelector string
	// kubernetes fieldSelector example `app=test`
	FieldSelector string
	// set KubeConfig out-of-cluster Use outside cluster
	KubeConfig string
	// set master url
	Master string
	// client rate limit queries per second, default 5
	QPS int
	// client token bucket burst, default 10
	Burst int
	// single request timeout
	Timeout time.Duration
}

// WithNamespace with kubernetes namespace.
func WithNamespace(ns string) Option {
	return func(o *options) {
		o.Namespace = ns
	}
}

// WithLabelSelector with kubernetes label selector.
func WithLabelSelector(label string) Option {
	return func(o *options) {
		o.LabelSelector = label
	}
}

// WithFieldSelector with kubernetes field selector.
func WithFieldSelector(field string) Option {
	return func(o *options) {
		o.FieldSelector = field
	}
}

// WithKubeConfig with kubernetes config.
func WithKubeConfig(config string) Option {
	return func(o *options) {
		o.KubeConfig = config
	}
}

// WithMaster with kubernetes master.
func WithQPS(qps int) Option {
	return func(op *options) {
		op.QPS = qps
	}
}

func WithBurst(burst int) Option {
	return func(op *options) {
		op.Burst = burst
	}
}

func WithRequestTimeout(timeout time.Duration) Option {
	return func(op *options) {
		op.Timeout = timeout
	}
}

func WithMaster(master string) Option {
	return func(o *options) {
		o.Master = master
	}
}
