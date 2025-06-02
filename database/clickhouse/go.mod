module github.com/tx7do/kratos-bootstrap/database/clickhouse

go 1.23.0

toolchain go1.23.3

replace github.com/tx7do/kratos-bootstrap/api => ../../api

require (
	github.com/ClickHouse/clickhouse-go/v2 v2.35.0
	github.com/go-kratos/kratos/v2 v2.8.4
	github.com/tx7do/kratos-bootstrap/api v0.0.20
	github.com/tx7do/kratos-bootstrap/utils v0.1.3
)

require (
	github.com/ClickHouse/ch-go v0.66.0 // indirect
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/paulmach/orb v0.11.1 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	go.opentelemetry.io/otel v1.36.0 // indirect
	go.opentelemetry.io/otel/trace v1.36.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
