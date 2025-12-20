module github.com/tx7do/kratos-bootstrap/script_engine

go 1.24.6

replace github.com/tx7do/kratos-bootstrap/api => ../api

require (
	github.com/tx7do/go-scripts v0.0.2
	github.com/tx7do/kratos-bootstrap/api v0.0.32
)

require google.golang.org/protobuf v1.36.11 // indirect
