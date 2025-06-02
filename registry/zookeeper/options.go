package zookeeper

// Option is etcd registry option.
type Option func(o *options)

type options struct {
	namespace string
	user      string
	password  string
}

// WithRootPath with registry root path.
func WithRootPath(path string) Option {
	return func(o *options) { o.namespace = path }
}

// WithDigestACL with registry password.
func WithDigestACL(user string, password string) Option {
	return func(o *options) {
		o.user = user
		o.password = password
	}
}
