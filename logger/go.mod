module github.com/tx7do/kratos-bootstrap/logger

go 1.23.0

toolchain go1.24.3

replace (
	github.com/armon/go-metrics => github.com/hashicorp/go-metrics v0.4.1
	github.com/imdario/mergo => dario.cat/mergo v0.3.16

	github.com/tx7do/kratos-bootstrap/api => ../api
	github.com/tx7do/kratos-bootstrap/utils => ../utils
)

require (
	github.com/go-kratos/kratos/contrib/log/aliyun/v2 v2.0.0-20250429074618-c82f7957223f
	github.com/go-kratos/kratos/contrib/log/fluent/v2 v2.0.0-20250429074618-c82f7957223f
	github.com/go-kratos/kratos/contrib/log/logrus/v2 v2.0.0-20250429074618-c82f7957223f
	github.com/go-kratos/kratos/contrib/log/tencent/v2 v2.0.0-20250429074618-c82f7957223f
	github.com/go-kratos/kratos/contrib/log/zap/v2 v2.0.0-20250429074618-c82f7957223f
	github.com/go-kratos/kratos/v2 v2.8.4
	github.com/sirupsen/logrus v1.9.3
	github.com/tx7do/kratos-bootstrap/api v0.0.13
	github.com/tx7do/kratos-bootstrap/utils v0.1.2
	go.uber.org/zap v1.27.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require (
	github.com/aliyun/aliyun-log-go-sdk v0.1.100 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fluent/fluent-logger-golang v1.9.0 // indirect
	github.com/go-kit/kit v0.13.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-kratos/aegis v0.2.0 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/form/v4 v4.2.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/philhofer/fwd v1.1.3-0.20240916144458-20a13a1f6b7c // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/tencentcloud/tencentcloud-cls-sdk-go v1.0.11 // indirect
	github.com/tinylib/msgp v1.2.5 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.35.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250505200425-f936aa4a68b2 // indirect
	google.golang.org/grpc v1.72.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/ini.v1 v1.66.4 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
