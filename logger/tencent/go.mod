module github.com/tx7do/kratos-bootstrap/logger/tencent

go 1.24.6

replace (
	github.com/tx7do/kratos-bootstrap/api => ../../api
	github.com/tx7do/kratos-bootstrap/logger => ../
)

require (
	github.com/go-kratos/kratos/v2 v2.9.2
	github.com/tencentcloud/tencentcloud-cls-sdk-go v1.0.14
	github.com/tx7do/kratos-bootstrap/api v0.0.32
	github.com/tx7do/kratos-bootstrap/logger v0.1.1
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-kratos/aegis v0.2.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/form/v4 v4.3.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/klauspost/compress v1.18.2 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.39.0 // indirect
	go.opentelemetry.io/otel/metric v1.39.0 // indirect
	go.opentelemetry.io/otel/trace v1.39.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251213004720-97cd9d5aeac2 // indirect
	google.golang.org/grpc v1.77.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
