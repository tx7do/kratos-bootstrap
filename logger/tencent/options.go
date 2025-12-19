package tencent

type options struct {
	topicID      string
	accessKey    string
	accessSecret string
	endpoint     string
}

func defaultOptions() *options {
	return &options{}
}

func WithEndpoint(endpoint string) Option {
	return func(cls *options) {
		cls.endpoint = endpoint
	}
}

func WithTopicID(topicID string) Option {
	return func(cls *options) {
		cls.topicID = topicID
	}
}

func WithAccessKey(ak string) Option {
	return func(cls *options) {
		cls.accessKey = ak
	}
}

func WithAccessSecret(as string) Option {
	return func(cls *options) {
		cls.accessSecret = as
	}
}

type Option func(cls *options)
