module github.com/tx7do/kratos-bootstrap/database/doris

go 1.25.3

replace github.com/tx7do/kratos-bootstrap/api => ../../api

require (
	github.com/go-kratos/kratos/v2 v2.9.2
	github.com/tx7do/go-crud/doris v0.0.1
	github.com/tx7do/kratos-bootstrap/api v0.0.36
)

require (
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/jmoiron/sqlx v1.4.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
