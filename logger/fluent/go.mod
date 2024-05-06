module github.com/tx7do/kratos-bootstrap/logger/fluent

go 1.22.0

toolchain go1.22.1

replace github.com/tx7do/kratos-bootstrap/api => ../../api

require (
	github.com/go-kratos/kratos/contrib/log/fluent/v2 v2.0.0-20240504101732-d0d5761f9ca8
	github.com/go-kratos/kratos/v2 v2.7.3
	github.com/tx7do/kratos-bootstrap/api v0.0.2
)

require (
	github.com/fluent/fluent-logger-golang v1.9.0 // indirect
	github.com/philhofer/fwd v1.1.2 // indirect
	github.com/tinylib/msgp v1.1.9 // indirect
	google.golang.org/protobuf v1.34.0 // indirect
)
