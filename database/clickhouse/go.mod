module github.com/tx7do/kratos-bootstrap/database/clickhouse

go 1.25.3

replace (
	github.com/tx7do/kratos-bootstrap/api => ../../api
	github.com/tx7do/kratos-bootstrap/utils => ../../utils
)

require (
	github.com/ClickHouse/clickhouse-go/v2 v2.41.0
	github.com/go-kratos/kratos/v2 v2.9.1
	github.com/stretchr/testify v1.11.1
	github.com/tx7do/kratos-bootstrap/api v0.0.28
	github.com/tx7do/kratos-bootstrap/utils v0.1.5
	google.golang.org/protobuf v1.36.10
)

require (
	github.com/ClickHouse/ch-go v0.69.0 // indirect
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.18.1 // indirect
	github.com/paulmach/orb v0.12.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	go.opentelemetry.io/otel v1.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/sys v0.38.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251124214823-79d6a2a48846 // indirect
	google.golang.org/grpc v1.77.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
