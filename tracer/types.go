package tracer

type Type string

const (
	Std      Type = "std"
	Zipkin   Type = "zipkin"
	Jaeger   Type = "jaeger"
	OtlpHttp Type = "otlp-http"
	OtlpGrpc Type = "otlp-grpc"
	Aliyun   Type = "aliyun"
	Tencent  Type = "tencent"
)
