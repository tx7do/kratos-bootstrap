module github.com/tx7do/kratos-bootstrap/logger/tencent

go 1.22.0

toolchain go1.22.1

replace github.com/tx7do/kratos-bootstrap/api => ../../api

require (
	github.com/go-kratos/kratos/contrib/log/tencent/v2 v2.0.0-20241014124617-8b8dc4b0f8be
	github.com/go-kratos/kratos/v2 v2.8.1
	github.com/tx7do/kratos-bootstrap/api v0.0.4
)

require (
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/tencentcloud/tencentcloud-cls-sdk-go v1.0.11 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
)
