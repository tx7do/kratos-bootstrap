module github.com/tx7do/kratos-bootstrap/config/polaris

go 1.22.0

toolchain go1.22.1

replace github.com/tx7do/kratos-bootstrap/api => ../../api

require (
	github.com/go-kratos/kratos/v2 v2.8.1
	github.com/tx7do/kratos-bootstrap/api v0.0.4
)

require (
	dario.cat/mergo v1.0.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
