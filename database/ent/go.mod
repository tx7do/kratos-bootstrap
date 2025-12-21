module github.com/tx7do/kratos-bootstrap/database/ent

go 1.24.6

replace github.com/tx7do/kratos-bootstrap/api => ../../api

require (
	entgo.io/ent v0.14.5
	github.com/go-kratos/kratos/v2 v2.9.2
	github.com/tx7do/go-crud/entgo v0.0.15
	github.com/tx7do/kratos-bootstrap/api v0.0.33
)

require (
	github.com/XSAM/otelsql v0.41.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/gnostic v0.7.1 // indirect
	github.com/google/gnostic-models v0.7.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/copier v0.4.0 // indirect
	github.com/tx7do/go-crud v0.0.6 // indirect
	github.com/tx7do/go-utils v1.1.34 // indirect
	github.com/tx7do/go-utils/mapper v0.0.3 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.39.0 // indirect
	go.opentelemetry.io/otel/metric v1.39.0 // indirect
	go.opentelemetry.io/otel/trace v1.39.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
