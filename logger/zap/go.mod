module github.com/tx7do/kratos-bootstrap/logger/zap

go 1.22.0

toolchain go1.22.1

replace github.com/tx7do/kratos-bootstrap/api => ../../api

require (
	github.com/go-kratos/kratos/contrib/log/zap/v2 v2.0.0-20241014124617-8b8dc4b0f8be
	github.com/go-kratos/kratos/v2 v2.8.1
	github.com/tx7do/kratos-bootstrap/api v0.0.4
	go.uber.org/zap v1.27.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require (
	go.uber.org/multierr v1.11.0 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
)
