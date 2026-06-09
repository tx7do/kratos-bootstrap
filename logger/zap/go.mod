module github.com/tx7do/kratos-bootstrap/logger/zap

go 1.25.0

replace (
	github.com/tx7do/kratos-bootstrap/api => ../../api
	github.com/tx7do/kratos-bootstrap/logger => ../
)

require (
	github.com/tx7do/kratos-bootstrap/api v0.0.33
	github.com/tx7do/kratos-bootstrap/logger v0.1.2
	go.uber.org/zap v1.27.1
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-kratos/kratos/v2 v2.9.2 // indirect
	go.opentelemetry.io/otel v1.43.0 // indirect
	go.opentelemetry.io/otel/trace v1.43.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
