module github.com/tx7do/kratos-bootstrap/logger

go 1.24.6

replace (
	github.com/armon/go-metrics => github.com/hashicorp/go-metrics v0.4.1
	github.com/imdario/mergo => dario.cat/mergo v0.3.16

	github.com/tx7do/kratos-bootstrap/api => ../api
	github.com/tx7do/kratos-bootstrap/utils => ../utils
)

require (
	github.com/go-kratos/kratos/contrib/log/aliyun/v2 v2.0.0-20251205160234-b9fab9a5a5ab
	github.com/go-kratos/kratos/contrib/log/fluent/v2 v2.0.0-20251205160234-b9fab9a5a5ab
	github.com/go-kratos/kratos/contrib/log/logrus/v2 v2.0.0-20251205160234-b9fab9a5a5ab
	github.com/go-kratos/kratos/contrib/log/tencent/v2 v2.0.0-20251205160234-b9fab9a5a5ab
	github.com/go-kratos/kratos/contrib/log/zap/v2 v2.0.0-20251205160234-b9fab9a5a5ab
	github.com/go-kratos/kratos/v2 v2.9.2
	github.com/sirupsen/logrus v1.9.3
	github.com/tx7do/kratos-bootstrap/api v0.0.29
	github.com/tx7do/kratos-bootstrap/utils v0.1.8
	go.uber.org/zap v1.27.1
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require (
	github.com/aliyun/aliyun-log-go-sdk v0.1.112 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dennwc/varint v1.0.0 // indirect
	github.com/fluent/fluent-logger-golang v1.10.1 // indirect
	github.com/go-kit/kit v0.13.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-kratos/aegis v0.2.0 // indirect
	github.com/go-logfmt/logfmt v0.6.1 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/form/v4 v4.3.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/grafana/regexp v0.0.0-20250905093917-f7b3be9d1853 // indirect
	github.com/klauspost/compress v1.18.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.23.2 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.67.4 // indirect
	github.com/prometheus/procfs v0.19.2 // indirect
	github.com/prometheus/prometheus v0.308.0 // indirect
	github.com/tencentcloud/tencentcloud-cls-sdk-go v1.0.14 // indirect
	github.com/tinylib/msgp v1.6.1 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.39.0 // indirect
	go.opentelemetry.io/otel/metric v1.39.0 // indirect
	go.opentelemetry.io/otel/trace v1.39.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.yaml.in/yaml/v2 v2.4.3 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/grpc v1.77.0 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
