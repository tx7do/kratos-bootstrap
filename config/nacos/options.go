package nacos

const (
	DefaultGroup  = "DEFAULT_GROUP"
	DefaultDataID = "bootstrap.yaml"
)

type Option func(*options)

type options struct {
	group  string
	dataID string
}

// WithGroup With nacos config group.
func WithGroup(group string) Option {
	return func(o *options) {
		o.group = group
	}
}

// WithDataID With nacos config data id.
func WithDataID(dataID string) Option {
	return func(o *options) {
		o.dataID = dataID
	}
}
