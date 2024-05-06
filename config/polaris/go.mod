module github.com/tx7do/kratos-bootstrap/config/polaris

go 1.22.0

toolchain go1.22.1

replace github.com/tx7do/kratos-bootstrap/api => ../../api

require (
	github.com/go-kratos/kratos/v2 v2.7.3
	github.com/tx7do/kratos-bootstrap/api v0.0.0-00010101000000-000000000000
)

require (
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/kr/text v0.2.0 // indirect
	google.golang.org/protobuf v1.34.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
