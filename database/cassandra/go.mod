module github.com/tx7do/kratos-bootstrap/database/cassandra

go 1.23.0

toolchain go1.23.3

replace github.com/tx7do/kratos-bootstrap/api => ../../api

require (
	github.com/go-kratos/kratos/v2 v2.8.2
	github.com/gocql/gocql v1.7.0
	github.com/tx7do/kratos-bootstrap/api v0.0.5
)

require (
	github.com/golang/snappy v0.0.4 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/kr/text v0.2.0 // indirect
	golang.org/x/sync v0.9.0 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
)
