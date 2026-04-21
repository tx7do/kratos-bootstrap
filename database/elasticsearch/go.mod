module github.com/tx7do/kratos-bootstrap/database/elasticsearch

go 1.25.0

replace github.com/tx7do/kratos-bootstrap/api => ../../api

require (
	github.com/go-kratos/kratos/v2 v2.9.2
	github.com/tx7do/go-crud/elasticsearch v0.0.10
	github.com/tx7do/kratos-bootstrap/api v0.0.39
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/elastic/elastic-transport-go/v8 v8.11.0 // indirect
	github.com/elastic/go-elasticsearch/v9 v9.3.2 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/gnostic v0.7.1 // indirect
	github.com/google/gnostic-models v0.7.1 // indirect
	github.com/tx7do/go-crud/api v0.0.7 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.43.0 // indirect
	go.opentelemetry.io/otel/metric v1.43.0 // indirect
	go.opentelemetry.io/otel/trace v1.43.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/sys v0.43.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260420184626-e10c466a9529 // indirect
	google.golang.org/grpc v1.80.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
