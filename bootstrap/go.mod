module github.com/tx7do/kratos-bootstrap/bootstrap

go 1.25.3

replace (
	github.com/armon/go-metrics => github.com/hashicorp/go-metrics v0.4.1

	github.com/tx7do/kratos-bootstrap/api => ../api
	github.com/tx7do/kratos-bootstrap/logger => ../logger

	github.com/tx7do/kratos-bootstrap/registry => ../registry
	github.com/tx7do/kratos-bootstrap/registry/consul => ../registry/consul
	github.com/tx7do/kratos-bootstrap/registry/etcd => ../registry/etcd
	github.com/tx7do/kratos-bootstrap/registry/eureka => ../registry/eureka
	github.com/tx7do/kratos-bootstrap/registry/kubernetes => ../registry/kubernetes
	github.com/tx7do/kratos-bootstrap/registry/nacos => ../registry/nacos
	github.com/tx7do/kratos-bootstrap/registry/polaris => ../registry/polaris
	github.com/tx7do/kratos-bootstrap/registry/servicecomb => ../registry/servicecomb
	github.com/tx7do/kratos-bootstrap/registry/zookeeper => ../registry/zookeeper

	github.com/tx7do/kratos-bootstrap/remoteconfig/apollo => ../remoteconfig/apollo
	github.com/tx7do/kratos-bootstrap/remoteconfig/consul => ../remoteconfig/consul
	github.com/tx7do/kratos-bootstrap/remoteconfig/etcd => ../remoteconfig/etcd
	github.com/tx7do/kratos-bootstrap/remoteconfig/kubernetes => ../remoteconfig/kubernetes
	github.com/tx7do/kratos-bootstrap/remoteconfig/nacos => ../remoteconfig/nacos
	github.com/tx7do/kratos-bootstrap/remoteconfig/polaris => ../remoteconfig/polaris
	github.com/tx7do/kratos-bootstrap/tracer => ../tracer
	github.com/tx7do/kratos-bootstrap/utils => ../utils
)

require (
	github.com/go-kratos/kratos/v2 v2.9.2
	github.com/google/subcommands v1.2.0
	github.com/olekukonko/tablewriter v1.1.2
	github.com/spf13/cobra v1.10.2
	github.com/tx7do/kratos-bootstrap/api v0.0.29
	github.com/tx7do/kratos-bootstrap/logger v0.0.12
	github.com/tx7do/kratos-bootstrap/registry v0.1.0
	github.com/tx7do/kratos-bootstrap/registry/consul v0.1.2
	github.com/tx7do/kratos-bootstrap/registry/etcd v0.1.2
	github.com/tx7do/kratos-bootstrap/registry/eureka v0.1.2
	github.com/tx7do/kratos-bootstrap/registry/nacos v0.1.2
	github.com/tx7do/kratos-bootstrap/registry/servicecomb v0.1.2
	github.com/tx7do/kratos-bootstrap/registry/zookeeper v0.1.2
	github.com/tx7do/kratos-bootstrap/remoteconfig/apollo v0.1.2
	github.com/tx7do/kratos-bootstrap/remoteconfig/consul v0.1.2
	github.com/tx7do/kratos-bootstrap/remoteconfig/etcd v0.1.2
	github.com/tx7do/kratos-bootstrap/remoteconfig/nacos v0.1.3
	github.com/tx7do/kratos-bootstrap/remoteconfig/polaris v0.1.2
	github.com/tx7do/kratos-bootstrap/tracer v0.0.13
	github.com/tx7do/kratos-bootstrap/utils v0.1.8
	golang.org/x/tools v0.40.0
)

require (
	dario.cat/mergo v1.0.2 // indirect
	github.com/alibabacloud-go/alibabacloud-gateway-pop v0.1.3 // indirect
	github.com/alibabacloud-go/alibabacloud-gateway-spi v0.0.5 // indirect
	github.com/alibabacloud-go/darabonba-array v0.1.0 // indirect
	github.com/alibabacloud-go/darabonba-encode-util v0.0.2 // indirect
	github.com/alibabacloud-go/darabonba-map v0.0.2 // indirect
	github.com/alibabacloud-go/darabonba-openapi/v2 v2.1.13 // indirect
	github.com/alibabacloud-go/darabonba-signature-util v0.0.7 // indirect
	github.com/alibabacloud-go/darabonba-string v1.0.2 // indirect
	github.com/alibabacloud-go/debug v1.0.1 // indirect
	github.com/alibabacloud-go/endpoint-util v1.1.1 // indirect
	github.com/alibabacloud-go/kms-20160120/v3 v3.4.0 // indirect
	github.com/alibabacloud-go/openapi-util v0.1.1 // indirect
	github.com/alibabacloud-go/tea v1.3.14 // indirect
	github.com/alibabacloud-go/tea-utils/v2 v2.0.9 // indirect
	github.com/alibabacloud-go/tea-xml v1.1.3 // indirect
	github.com/aliyun/alibaba-cloud-sdk-go v1.63.107 // indirect
	github.com/aliyun/alibabacloud-dkms-gcs-go-sdk v0.5.1 // indirect
	github.com/aliyun/alibabacloud-dkms-transfer-go-sdk v0.1.9 // indirect
	github.com/aliyun/aliyun-secretsmanager-client-go v1.1.5 // indirect
	github.com/aliyun/credentials-go v1.4.9 // indirect
	github.com/apolloconfig/agollo/v4 v4.4.0 // indirect
	github.com/armon/go-metrics v0.5.4 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/clbanning/mxj/v2 v2.7.0 // indirect
	github.com/clipperhouse/displaywidth v0.6.1 // indirect
	github.com/clipperhouse/stringish v0.1.1 // indirect
	github.com/clipperhouse/uax29/v2 v2.3.0 // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd/v22 v22.6.0 // indirect
	github.com/deckarep/golang-set v1.8.0 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/fluent/fluent-logger-golang v1.10.1 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-chassis/cari v0.9.0 // indirect
	github.com/go-chassis/foundation v0.4.0 // indirect
	github.com/go-chassis/openlog v1.1.3 // indirect
	github.com/go-chassis/sc-client v0.7.0 // indirect
	github.com/go-kratos/aegis v0.2.0 // indirect
	github.com/go-kratos/kratos/contrib/log/fluent/v2 v2.0.0-20251205160234-b9fab9a5a5ab // indirect
	github.com/go-kratos/kratos/contrib/log/logrus/v2 v2.0.0-20251205160234-b9fab9a5a5ab // indirect
	github.com/go-kratos/kratos/contrib/log/tencent/v2 v2.0.0-20251205160234-b9fab9a5a5ab // indirect
	github.com/go-kratos/kratos/contrib/log/zap/v2 v2.0.0-20251205160234-b9fab9a5a5ab // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/form/v4 v4.3.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/go-zookeeper/zk v1.0.4 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/gorilla/websocket v1.5.4-0.20250319132907-e064f32e3674 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3 // indirect
	github.com/hashicorp/consul/api v1.33.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-metrics v0.5.4 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/serf v0.10.2 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/karlseguin/ccache/v2 v2.0.8 // indirect
	github.com/klauspost/compress v1.18.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.19 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nacos-group/nacos-sdk-go/v2 v2.3.5 // indirect
	github.com/olekukonko/cat v0.0.0-20250911104152-50322a0618f6 // indirect
	github.com/olekukonko/errors v1.1.0 // indirect
	github.com/olekukonko/ll v0.1.3 // indirect
	github.com/opentracing/opentracing-go v1.2.1-0.20220228012449-10b1cf09e00b // indirect
	github.com/openzipkin/zipkin-go v0.4.3 // indirect
	github.com/orcaman/concurrent-map v1.0.0 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.23.2 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.67.4 // indirect
	github.com/prometheus/procfs v0.19.2 // indirect
	github.com/sagikazarmark/locafero v0.12.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/spf13/viper v1.21.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/tencentcloud/tencentcloud-cls-sdk-go v1.0.14 // indirect
	github.com/tinylib/msgp v1.6.1 // indirect
	github.com/tjfoc/gmsm v1.4.1 // indirect
	go.etcd.io/etcd/api/v3 v3.6.6 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.6.6 // indirect
	go.etcd.io/etcd/client/v3 v3.6.6 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.39.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.39.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.39.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.39.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.39.0 // indirect
	go.opentelemetry.io/otel/exporters/zipkin v1.39.0 // indirect
	go.opentelemetry.io/otel/metric v1.39.0 // indirect
	go.opentelemetry.io/otel/sdk v1.39.0 // indirect
	go.opentelemetry.io/otel/trace v1.39.0 // indirect
	go.opentelemetry.io/proto/otlp v1.9.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.1 // indirect
	go.yaml.in/yaml/v2 v2.4.3 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/crypto v0.46.0 // indirect
	golang.org/x/exp v0.0.0-20251125195548-87e1e737ad39 // indirect
	golang.org/x/mod v0.31.0 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	golang.org/x/time v0.14.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/grpc v1.77.0 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
