module github.com/tx7do/kratos-bootstrap/database/clickhouse

go 1.25.3

replace (
	github.com/tx7do/kratos-bootstrap/api => ../../api
	github.com/tx7do/kratos-bootstrap/utils => ../../utils
)

require (
	github.com/ClickHouse/clickhouse-go/v2 v2.41.0
	github.com/go-kratos/kratos/v2 v2.9.2
	github.com/tx7do/go-crud/clickhouse v0.0.3
	github.com/tx7do/kratos-bootstrap/api v0.0.29
	github.com/tx7do/kratos-bootstrap/utils v0.1.7
)

require (
	github.com/ClickHouse/ch-go v0.69.0 // indirect
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/google/gnostic v0.7.1 // indirect
	github.com/google/gnostic-models v0.7.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/copier v0.4.0 // indirect
	github.com/klauspost/compress v1.18.2 // indirect
	github.com/paulmach/orb v0.12.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/tx7do/go-crud v0.0.5 // indirect
	github.com/tx7do/go-utils v1.1.34 // indirect
	github.com/tx7do/go-utils/mapper v0.0.3 // indirect
	go.opentelemetry.io/otel v1.39.0 // indirect
	go.opentelemetry.io/otel/trace v1.39.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/sys v0.39.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/grpc v1.77.0 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)
