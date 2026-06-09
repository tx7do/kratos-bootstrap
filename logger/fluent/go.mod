module github.com/tx7do/kratos-bootstrap/logger/fluent

go 1.25.0

replace (
	github.com/tx7do/kratos-bootstrap/api => ../../api
	github.com/tx7do/kratos-bootstrap/logger => ../
)

require (
	github.com/fluent/fluent-logger-golang v1.10.1
	github.com/tx7do/kratos-bootstrap/api v0.0.33
	github.com/tx7do/kratos-bootstrap/logger v0.1.2
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-kratos/kratos/v2 v2.9.2 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/tinylib/msgp v1.6.4 // indirect
	go.opentelemetry.io/otel v1.43.0 // indirect
	go.opentelemetry.io/otel/trace v1.43.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
