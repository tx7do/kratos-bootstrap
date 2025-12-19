module github.com/tx7do/kratos-bootstrap/config

go 1.24.6

replace github.com/tx7do/kratos-bootstrap/api => ../api

require (
	github.com/go-kratos/kratos/v2 v2.9.2
	github.com/tx7do/kratos-bootstrap/api v0.0.32
	google.golang.org/protobuf v1.36.11
)

require (
	dario.cat/mergo v1.0.2 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
