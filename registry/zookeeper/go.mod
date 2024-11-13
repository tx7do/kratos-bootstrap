module github.com/tx7do/kratos-bootstrap/registry/zookeeper

go 1.22.0

toolchain go1.22.1

replace github.com/tx7do/kratos-bootstrap/api => ../../api

require (
	github.com/go-kratos/kratos/contrib/registry/zookeeper/v2 v2.0.0-20241105072421-f8b97f675b32
	github.com/go-kratos/kratos/v2 v2.8.2
	github.com/go-zookeeper/zk v1.0.4
	github.com/stretchr/testify v1.9.0
	github.com/tx7do/kratos-bootstrap/api v0.0.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sync v0.9.0 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
