package polaris

import "github.com/polarismesh/polaris-go"

// Option is polaris config option.
type Option func(o *options)

type options struct {
	namespace  string
	fileGroup  string
	fileName   string
	configFile polaris.ConfigFile
}

// WithNamespace with polaris config namespace
func WithNamespace(namespace string) Option {
	return func(o *options) {
		o.namespace = namespace
	}
}

// WithFileGroup with polaris config fileGroup
func WithFileGroup(fileGroup string) Option {
	return func(o *options) {
		o.fileGroup = fileGroup
	}
}

// WithFileName with polaris config fileName
func WithFileName(fileName string) Option {
	return func(o *options) {
		o.fileName = fileName
	}
}
