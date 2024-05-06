module github.com/tx7do/kratos-bootstrap/logger/logrus

go 1.22.0

toolchain go1.22.1

replace github.com/tx7do/kratos-bootstrap/api => ../../api

require (
	github.com/go-kratos/kratos/contrib/log/logrus/v2 v2.0.0-20240504101732-d0d5761f9ca8
	github.com/go-kratos/kratos/v2 v2.7.3
	github.com/sirupsen/logrus v1.9.3
	github.com/tx7do/kratos-bootstrap/api v0.0.2
)

require (
	golang.org/x/sys v0.20.0 // indirect
	google.golang.org/protobuf v1.34.0 // indirect
)
