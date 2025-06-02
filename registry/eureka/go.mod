module github.com/tx7do/kratos-bootstrap/registry/eureka

go 1.24.0

toolchain go1.24.3

replace (
	github.com/armon/go-metrics => github.com/hashicorp/go-metrics v0.4.1

	github.com/tx7do/kratos-bootstrap/api => ../../api
	github.com/tx7do/kratos-bootstrap/registry => ../
)

require (
	github.com/go-kratos/kratos/v2 v2.8.4
	github.com/stretchr/testify v1.10.0
	github.com/tx7do/kratos-bootstrap/api v0.0.21
	github.com/tx7do/kratos-bootstrap/registry v0.1.0
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	golang.org/x/sync v0.14.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
