package fluent

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"

	"github.com/fluent/fluent-logger-golang/fluent"

	bLogger "github.com/tx7do/kratos-bootstrap/logger"
)

// Compile-time assertion: Logger implements bLogger.Logger.
var _ bLogger.Logger = (*Logger)(nil)

// Logger is fluent logger sdk.
type Logger struct {
	opts  options
	log   *fluent.Fluent
	extra []any
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

// Debug 输出 DEBUG 级别日志。
func (l *Logger) Debug(_ context.Context, msg string, keyvals ...any) {
	l.post("DEBUG", msg, keyvals)
}

// Info 输出 INFO 级别日志。
func (l *Logger) Info(_ context.Context, msg string, keyvals ...any) {
	l.post("INFO", msg, keyvals)
}

// Warn 输出 WARN 级别日志。
func (l *Logger) Warn(_ context.Context, msg string, keyvals ...any) {
	l.post("WARN", msg, keyvals)
}

// Error 输出 ERROR 级别日志。
func (l *Logger) Error(_ context.Context, msg string, keyvals ...any) {
	l.post("ERROR", msg, keyvals)
}

// With 返回附加了指定 key-value 对的新 Logger 实例。
func (l *Logger) With(keyvals ...any) bLogger.Logger {
	return &Logger{
		opts:  l.opts,
		log:   l.log,
		extra: append(append([]any{}, l.extra...), keyvals...),
	}
}

// post 发送日志到 fluentd 服务。
func (l *Logger) post(tag, msg string, keyvals []any) {
	data := make(map[string]string, 3+len(l.extra)/2+len(keyvals)/2)
	data["level"] = tag
	data["msg"] = msg
	all := append(append([]any{}, l.extra...), keyvals...)
	for i := 0; i+1 < len(all); i += 2 {
		data[fmt.Sprint(all[i])] = fmt.Sprint(all[i+1])
	}
	if err := l.log.Post(tag, data); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

// Close the logger.
func (l *Logger) Close() error {
	return l.log.Close()
}
