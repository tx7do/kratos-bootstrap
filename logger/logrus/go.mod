module github.com/tx7do/kratos-bootstrap/logger/logrus

go 1.22.0

toolchain go1.22.1

replace github.com/tx7do/kratos-bootstrap/api => ../../api

require (
	github.com/go-kratos/kratos/contrib/log/logrus/v2 v2.0.0-20241014124617-8b8dc4b0f8be
	github.com/go-kratos/kratos/v2 v2.8.1
	github.com/sirupsen/logrus v1.9.3
	github.com/tx7do/kratos-bootstrap/api v0.0.4
)

require (
	golang.org/x/sys v0.26.0 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
)
