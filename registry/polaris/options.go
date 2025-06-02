package polaris

import "time"

type options struct {
	// required, namespace in polaris
	Namespace string

	// required, service access token
	ServiceToken string

	// optional, protocol in polaris. The Default value is nil, it means use protocol config in service
	Protocol *string

	// service weight in polaris. The Default value is 100, 0 <= weight <= 10000
	Weight int

	// service priority. Default value is 0. The smaller the value, the lower the priority
	Priority int

	// To show service is healthy or not. The Default value is True.
	Healthy bool

	// Heartbeat enable .Not in polaris. The Default value is True.
	Heartbeat bool

	// To show service is isolated or not. Default value is False.
	Isolate bool

	// TTL timeout. if the node needs to use heartbeat to report, required. If not set,server will throw ErrorCode-400141
	TTL int

	// optional, Timeout for a single query. Default value is global config
	// Total is (1+RetryCount) * Timeout
	Timeout time.Duration

	// optional, retry count. Default value is global config
	RetryCount int
}

// Option is polaris option.
type Option func(o *options)

// WithNamespace with a Namespace option.
func WithNamespace(namespace string) Option {
	return func(o *options) { o.Namespace = namespace }
}

// WithServiceToken with ServiceToken option.
func WithServiceToken(serviceToken string) Option {
	return func(o *options) { o.ServiceToken = serviceToken }
}

// WithProtocol with a Protocol option.
func WithProtocol(protocol string) Option {
	return func(o *options) { o.Protocol = &protocol }
}

// WithWeight with a Weight option.
func WithWeight(weight int) Option {
	return func(o *options) { o.Weight = weight }
}

// WithHealthy with a Healthy option.
func WithHealthy(healthy bool) Option {
	return func(o *options) { o.Healthy = healthy }
}

// WithIsolate with an Isolate option.
func WithIsolate(isolate bool) Option {
	return func(o *options) { o.Isolate = isolate }
}

// WithTTL with TTL option.
func WithTTL(TTL int) Option {
	return func(o *options) { o.TTL = TTL }
}

// WithTimeout with a Timeout option.
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) { o.Timeout = timeout }
}

// WithRetryCount with RetryCount option.
func WithRetryCount(retryCount int) Option {
	return func(o *options) { o.RetryCount = retryCount }
}

// WithHeartbeat with a Heartbeat option.
func WithHeartbeat(heartbeat bool) Option {
	return func(o *options) { o.Heartbeat = heartbeat }
}
