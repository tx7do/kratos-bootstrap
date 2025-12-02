module github.com/tx7do/kratos-bootstrap/database/ent

go 1.24.6

replace github.com/tx7do/kratos-bootstrap/api => ../../api

require (
	entgo.io/ent v0.14.5
	github.com/XSAM/otelsql v0.40.0
	github.com/go-kratos/kratos/v2 v2.9.1
	github.com/go-sql-driver/mysql v1.9.3
	github.com/jackc/pgx/v5 v5.7.6
	github.com/lib/pq v1.10.9
	github.com/tx7do/kratos-bootstrap/api v0.0.28
	go.opentelemetry.io/otel v1.38.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel/metric v1.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
	golang.org/x/crypto v0.45.0 // indirect
	golang.org/x/sync v0.18.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)
