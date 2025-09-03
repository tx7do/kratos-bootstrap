module github.com/tx7do/kratos-bootstrap/remoteconfig/polaris

go 1.24.6

replace (
	github.com/armon/go-metrics => github.com/hashicorp/go-metrics v0.4.1

	github.com/tx7do/kratos-bootstrap/api => ../../api
	github.com/tx7do/kratos-bootstrap/utils => ../../utils
)

require (
	github.com/go-kratos/kratos/v2 v2.8.4
	github.com/tx7do/kratos-bootstrap/api v0.0.27
)

require (
	dario.cat/mergo v1.0.2 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	golang.org/x/sync v0.14.0 // indirect
	google.golang.org/protobuf v1.36.8 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
