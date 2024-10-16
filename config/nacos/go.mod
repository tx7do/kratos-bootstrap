module github.com/tx7do/kratos-bootstrap/config/nacos

go 1.22.0

toolchain go1.22.1

replace github.com/tx7do/kratos-bootstrap/api => ../../api

require (
	github.com/go-kratos/kratos/contrib/config/nacos/v2 v2.0.0-20241014124617-8b8dc4b0f8be
	github.com/go-kratos/kratos/v2 v2.8.1
	github.com/nacos-group/nacos-sdk-go v1.1.5
	github.com/tx7do/kratos-bootstrap/api v0.0.4
)

require (
	dario.cat/mergo v1.0.1 // indirect
	github.com/aliyun/alibaba-cloud-sdk-go v1.63.30 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/go-errors/errors v1.5.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/opentracing/opentracing-go v1.2.1-0.20220228012449-10b1cf09e00b // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
