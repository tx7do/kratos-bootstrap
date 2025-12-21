module github.com/tx7do/kratos-bootstrap/registry

go 1.24.6

replace github.com/tx7do/kratos-bootstrap/api => ../api

require (
	github.com/go-kratos/kratos/v2 v2.9.2
	github.com/tx7do/kratos-bootstrap/api v0.0.33
)

require google.golang.org/protobuf v1.36.11 // indirect
