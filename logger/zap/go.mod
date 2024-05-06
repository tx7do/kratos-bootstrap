module github.com/tx7do/kratos-bootstrap/logger/zap

go 1.22.0

toolchain go1.22.1

replace github.com/tx7do/kratos-bootstrap/api => ../../api

require (
	github.com/go-kratos/kratos/contrib/log/zap/v2 v2.0.0-20240504101732-d0d5761f9ca8
	github.com/go-kratos/kratos/v2 v2.7.3
	github.com/tx7do/kratos-bootstrap/api v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.27.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require (
	go.uber.org/multierr v1.11.0 // indirect
	google.golang.org/protobuf v1.34.0 // indirect
)
