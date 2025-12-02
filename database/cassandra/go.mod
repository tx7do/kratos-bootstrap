module github.com/tx7do/kratos-bootstrap/database/cassandra

go 1.24.6

replace github.com/tx7do/kratos-bootstrap/api => ../../api

require (
	github.com/go-kratos/kratos/v2 v2.9.1
	github.com/gocql/gocql v1.7.0
	github.com/tx7do/kratos-bootstrap/api v0.0.28
	github.com/tx7do/kratos-bootstrap/utils v0.1.5
)

require (
	github.com/golang/snappy v1.0.0 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/kr/text v0.2.0 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
)
