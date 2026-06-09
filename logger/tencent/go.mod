module github.com/tx7do/kratos-bootstrap/logger/tencent

go 1.25.0

replace (
	github.com/tx7do/kratos-bootstrap/api => ../../api
	github.com/tx7do/kratos-bootstrap/logger => ../
)

require (
	github.com/tencentcloud/tencentcloud-cls-sdk-go v1.0.14
	github.com/tx7do/kratos-bootstrap/api v0.0.33
	github.com/tx7do/kratos-bootstrap/logger v0.1.2
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-kratos/kratos/v2 v2.9.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/klauspost/compress v1.18.5 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	go.opentelemetry.io/otel v1.43.0 // indirect
	go.opentelemetry.io/otel/trace v1.43.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
)
