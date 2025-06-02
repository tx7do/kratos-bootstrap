module github.com/tx7do/kratos-bootstrap/database/ent

go 1.23.0

toolchain go1.23.3

replace github.com/tx7do/kratos-bootstrap/api => ../../api

require (
	entgo.io/ent v0.14.4
	github.com/XSAM/otelsql v0.38.0
	github.com/go-kratos/kratos/v2 v2.8.4
	github.com/go-sql-driver/mysql v1.9.2
	github.com/jackc/pgx/v4 v4.18.3
	github.com/lib/pq v1.10.9
	github.com/tx7do/kratos-bootstrap/api v0.0.18
	go.opentelemetry.io/otel v1.36.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgtype v1.14.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/metric v1.36.0 // indirect
	go.opentelemetry.io/otel/trace v1.36.0 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
