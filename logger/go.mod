module github.com/tx7do/kratos-bootstrap/logger

go 1.23

toolchain go1.23.2

replace (
	github.com/armon/go-metrics => github.com/hashicorp/go-metrics v0.4.1
	github.com/imdario/mergo => dario.cat/mergo v0.3.16

	github.com/tx7do/kratos-bootstrap/api => ../api
)
