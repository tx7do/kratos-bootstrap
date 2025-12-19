package kubernetes

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
func WithMaster(master string) Option {
	return func(o *options) {
		o.Master = master
	}
}
