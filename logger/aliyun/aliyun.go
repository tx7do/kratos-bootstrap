package aliyun

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"google.golang.org/protobuf/proto"

	bLogger "github.com/tx7do/kratos-bootstrap/logger"
)

// Logger 扩展了项目 Logger 接口，增加 Aliyun SLS 特有的方法。
type Logger interface {
	bLogger.Logger

	GetProducer() *producer.Producer
	Close() error
}

type aliyunLog struct {
	producer *producer.Producer
	opts     *options
	extra    []any
}

func (a *aliyunLog) GetProducer() *producer.Producer {
	return a.producer
}

func (a *aliyunLog) Close() error {
	return a.producer.Close(5000)
}

// Debug 输出 DEBUG 级别日志。
func (a *aliyunLog) Debug(_ context.Context, msg string, keyvals ...any) {
	a.post("DEBUG", msg, keyvals)
}

// Info 输出 INFO 级别日志。
func (a *aliyunLog) Info(_ context.Context, msg string, keyvals ...any) {
	a.post("INFO", msg, keyvals)
}

// Warn 输出 WARN 级别日志。
func (a *aliyunLog) Warn(_ context.Context, msg string, keyvals ...any) {
	a.post("WARN", msg, keyvals)
}

// Error 输出 ERROR 级别日志。
func (a *aliyunLog) Error(_ context.Context, msg string, keyvals ...any) {
	a.post("ERROR", msg, keyvals)
}

// With 返回附加了指定 key-value 对的新 Logger 实例。
func (a *aliyunLog) With(keyvals ...any) bLogger.Logger {
	return &aliyunLog{
		producer: a.producer,
		opts:     a.opts,
		extra:    append(append([]any{}, a.extra...), keyvals...),
	}
}

// post 发送日志到 Aliyun SLS。
func (a *aliyunLog) post(level, msg string, keyvals []any) {
	contents := make([]*sls.LogContent, 0, 3+len(a.extra)/2+len(keyvals)/2)
	contents = append(contents, &sls.LogContent{
		Key:   newString("level"),
		Value: newString(level),
	})
	if msg != "" {
		contents = append(contents, &sls.LogContent{
			Key:   newString("msg"),
			Value: newString(msg),
		})
	}
	all := append(append([]any{}, a.extra...), keyvals...)
	for i := 0; i+1 < len(all); i += 2 {
		contents = append(contents, &sls.LogContent{
			Key:   newString(toString(all[i])),
			Value: newString(toString(all[i+1])),
		})
	}

	logInst := &sls.Log{
		Time:     proto.Uint32(uint32(time.Now().Unix())),
		Contents: contents,
	}
	_ = a.producer.SendLog(a.opts.project, a.opts.logstore, "", "", logInst)
}

// NewAliyunLogger 创建 Aliyun SLS 日志记录器。
func NewAliyunLogger(options ...Option) (Logger, error) {
	opts := defaultOptions()
	for _, o := range options {
		o(opts)
	}

	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.Endpoint = opts.endpoint

	//producerConfig.AccessKeyID = opts.accessKey
	//producerConfig.AccessKeySecret = opts.accessSecret
	producerConfig.CredentialsProvider = sls.NewStaticCredentialsProvider(opts.accessKey, opts.accessSecret, opts.securityToken)

	producerInst, err := producer.NewProducer(producerConfig)
	if err != nil {
		return nil, err
	}
	producerInst.Start()

	return &aliyunLog{
		opts:     opts,
		producer: producerInst,
	}, nil
}

// newString string convert to *string
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
