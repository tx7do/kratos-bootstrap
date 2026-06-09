package tencent

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	cls "github.com/tencentcloud/tencentcloud-cls-sdk-go"
	"google.golang.org/protobuf/proto"

	bLogger "github.com/tx7do/kratos-bootstrap/logger"
)

type Logger interface {
	bLogger.Logger

	GetProducer() *cls.AsyncProducerClient
	Close() error
}

type tencentLog struct {
	producer *cls.AsyncProducerClient
	opts     *options
	extra    []any
}

var _ bLogger.Logger = (*tencentLog)(nil)

func (l *tencentLog) GetProducer() *cls.AsyncProducerClient {
	return l.producer
}

func (l *tencentLog) Close() error {
	return l.producer.Close(5000)
}

// Debug 输出 DEBUG 级别日志。
func (l *tencentLog) Debug(_ context.Context, msg string, keyvals ...any) {
	l.post("DEBUG", msg, keyvals)
}

// Info 输出 INFO 级别日志。
func (l *tencentLog) Info(_ context.Context, msg string, keyvals ...any) {
	l.post("INFO", msg, keyvals)
}

// Warn 输出 WARN 级别日志。
func (l *tencentLog) Warn(_ context.Context, msg string, keyvals ...any) {
	l.post("WARN", msg, keyvals)
}

// Error 输出 ERROR 级别日志。
func (l *tencentLog) Error(_ context.Context, msg string, keyvals ...any) {
	l.post("ERROR", msg, keyvals)
}

// With 返回附加了指定 key-value 对的新 Logger 实例。
func (l *tencentLog) With(keyvals ...any) bLogger.Logger {
	return &tencentLog{
		producer: l.producer,
		opts:     l.opts,
		extra:    append(append([]any{}, l.extra...), keyvals...),
	}
}

// post 发送日志到 Tencent CLS。
func (l *tencentLog) post(level, msg string, keyvals []any) {
	contents := make([]*cls.Log_Content, 0, 3+len(l.extra)/2+len(keyvals)/2)
	contents = append(contents, &cls.Log_Content{
		Key:   newString("level"),
		Value: newString(level),
	})
	if msg != "" {
		contents = append(contents, &cls.Log_Content{
			Key:   newString("msg"),
			Value: newString(msg),
		})
	}
	all := append(append([]any{}, l.extra...), keyvals...)
	for i := 0; i+1 < len(all); i += 2 {
		contents = append(contents, &cls.Log_Content{
			Key:   newString(toString(all[i])),
			Value: newString(toString(all[i+1])),
		})
	}

	logInst := &cls.Log{
		Time:     proto.Int64(time.Now().Unix()),
		Contents: contents,
	}
	_ = l.producer.SendLog(l.opts.topicID, logInst, nil)
}

func NewTencentLogger(options ...Option) (Logger, error) {
	opts := defaultOptions()
	for _, o := range options {
		o(opts)
	}
	producerConfig := cls.GetDefaultAsyncProducerClientConfig()
	producerConfig.AccessKeyID = opts.accessKey
	producerConfig.AccessKeySecret = opts.accessSecret
	producerConfig.Endpoint = opts.endpoint
	producerInst, err := cls.NewAsyncProducerClient(producerConfig)
	if err != nil {
		return nil, err
	}
	producerInst.Start()
	return &tencentLog{
		producer: producerInst,
		opts:     opts,
	}, nil
}

func newString(s string) *string {
	return &s
}

// toString convert any type to string
func toString(v any) string {
	var key string
	if v == nil {
		return key
	}
	switch v := v.(type) {
	case float64:
		key = strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		key = strconv.FormatFloat(float64(v), 'f', -1, 32)
	case int:
		key = strconv.Itoa(v)
	case uint:
		key = strconv.FormatUint(uint64(v), 10)
	case int8:
		key = strconv.Itoa(int(v))
	case uint8:
		key = strconv.FormatUint(uint64(v), 10)
	case int16:
		key = strconv.Itoa(int(v))
	case uint16:
		key = strconv.FormatUint(uint64(v), 10)
	case int32:
		key = strconv.Itoa(int(v))
	case uint32:
		key = strconv.FormatUint(uint64(v), 10)
	case int64:
		key = strconv.FormatInt(v, 10)
	case uint64:
		key = strconv.FormatUint(v, 10)
	case string:
		key = v
	case bool:
		key = strconv.FormatBool(v)
	case []byte:
		key = string(v)
	case fmt.Stringer:
		key = v.String()
	default:
		newValue, _ := json.Marshal(v)
		key = string(newValue)
	}
	return key
}
