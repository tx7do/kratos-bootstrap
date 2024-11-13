module github.com/tx7do/kratos-bootstrap/logger/logrus

go 1.22.0

toolchain go1.22.1

replace github.com/tx7do/kratos-bootstrap/api => ../../api

require (
	github.com/go-kratos/kratos/contrib/log/logrus/v2 v2.0.0-20241105072421-f8b97f675b32
	github.com/go-kratos/kratos/v2 v2.8.2
	github.com/sirupsen/logrus v1.9.3
	github.com/tx7do/kratos-bootstrap/api v0.0.5
)

require (
	golang.org/x/sys v0.27.0 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
)
