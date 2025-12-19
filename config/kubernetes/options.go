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

// Namespace with kubernetes namespace.
func Namespace(ns string) Option {
	return func(o *options) {
		o.Namespace = ns
	}
}

// LabelSelector with kubernetes label selector.
func LabelSelector(label string) Option {
	return func(o *options) {
		o.LabelSelector = label
	}
}

// FieldSelector with kubernetes field selector.
func FieldSelector(field string) Option {
	return func(o *options) {
		o.FieldSelector = field
	}
}

// KubeConfig with kubernetes config.
func KubeConfig(config string) Option {
	return func(o *options) {
		o.KubeConfig = config
	}
}

// Master with kubernetes master.
func Master(master string) Option {
	return func(o *options) {
		o.Master = master
	}
}
