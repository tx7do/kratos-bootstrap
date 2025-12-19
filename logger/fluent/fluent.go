package fluent

import (
	"fmt"
	"net"
	"net/url"
	"strconv"

	"github.com/fluent/fluent-logger-golang/fluent"

	"github.com/go-kratos/kratos/v2/log"
)

var _ log.Logger = (*Logger)(nil)

// Logger is fluent logger sdk.
type Logger struct {
	opts options
	log  *fluent.Fluent
}

// NewFluentLogger new a std logger with options.
// target:
//
//	tcp://127.0.0.1:24224
//	unix://var/run/fluent/fluent.sock
func NewFluentLogger(addr string, opts ...Option) (*Logger, error) {
	option := options{}
	for _, o := range opts {
		o(&option)
	}
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	c := fluent.Config{
		Timeout:            option.timeout,
		WriteTimeout:       option.writeTimeout,
		BufferLimit:        option.bufferLimit,
		RetryWait:          option.retryWait,
		MaxRetry:           option.maxRetry,
		MaxRetryWait:       option.maxRetryWait,
		TagPrefix:          option.tagPrefix,
		Async:              option.async,
		ForceStopAsyncSend: option.forceStopAsyncSend,
	}
	switch u.Scheme {
	case "tcp":
		host, port, err2 := net.SplitHostPort(u.Host)
		if err2 != nil {
			return nil, err2
		}
		if c.FluentPort, err = strconv.Atoi(port); err != nil {
			return nil, err
		}
		c.FluentNetwork = u.Scheme
		c.FluentHost = host
	case "unix":
		c.FluentNetwork = u.Scheme
		c.FluentSocketPath = u.Path
	default:
		return nil, fmt.Errorf("unknown network: %s", u.Scheme)
	}
	fl, err := fluent.New(c)
	if err != nil {
		return nil, err
	}
	return &Logger{
		opts: option,
		log:  fl,
	}, nil
}

// Log print the kv pairs log.
func (l *Logger) Log(level log.Level, keyvals ...any) error {
	if len(keyvals) == 0 {
		return nil
	}
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "KEYVALS UNPAIRED")
	}

	data := make(map[string]string, len(keyvals)/2+1)

	for i := 0; i < len(keyvals); i += 2 {
		data[fmt.Sprint(keyvals[i])] = fmt.Sprint(keyvals[i+1])
	}

	if err := l.log.Post(level.String(), data); err != nil {
		println(err)
	}
	return nil
}

// Close the logger.
func (l *Logger) Close() error {
	return l.log.Close()
}
