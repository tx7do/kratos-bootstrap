module github.com/tx7do/kratos-bootstrap/script_engine/source/database

go 1.24.6

replace (
	github.com/tx7do/kratos-bootstrap/api => ../../../api
	github.com/tx7do/kratos-bootstrap/script_engine => ../../
)

require (
	github.com/tx7do/go-scripts v0.0.6
	github.com/tx7do/kratos-bootstrap/api v0.0.33
	github.com/tx7do/kratos-bootstrap/script_engine v0.0.0-00010101000000-000000000000
)

require (
	github.com/go-kratos/kratos/v2 v2.9.2 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

// TODO: 当 go-scripts/source/database 发布后，添加以下依赖:
// require github.com/tx7do/go-scripts/source/database v0.0.0
