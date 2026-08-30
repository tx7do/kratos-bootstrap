module github.com/tx7do/kratos-bootstrap/registry

go 1.25.0

replace github.com/tx7do/kratos-bootstrap/api => ../api

require (
	github.com/go-kratos/kratos/v2 v2.9.2
	github.com/tx7do/go-utils v1.1.40
	github.com/tx7do/kratos-bootstrap/api v0.0.45
)

require google.golang.org/protobuf v1.36.12 // indirect
