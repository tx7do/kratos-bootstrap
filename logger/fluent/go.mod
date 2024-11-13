module github.com/tx7do/kratos-bootstrap/logger/fluent

go 1.22.0

toolchain go1.22.1

replace github.com/tx7do/kratos-bootstrap/api => ../../api

require (
	github.com/go-kratos/kratos/contrib/log/fluent/v2 v2.0.0-20241105072421-f8b97f675b32
	github.com/go-kratos/kratos/v2 v2.8.2
	github.com/tx7do/kratos-bootstrap/api v0.0.4
)

require (
	github.com/fluent/fluent-logger-golang v1.9.0 // indirect
	github.com/philhofer/fwd v1.1.3-0.20240916144458-20a13a1f6b7c // indirect
	github.com/tinylib/msgp v1.2.4 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
)
