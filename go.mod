module github.com/tx7do/kratos-bootstrap

go 1.20

require (
	github.com/go-chassis/sc-client v0.7.0
	github.com/go-kratos/aegis v0.2.0
	github.com/go-kratos/kratos/contrib/config/apollo/v2 v2.0.0-20230519061918-96480c11ee42
	github.com/go-kratos/kratos/contrib/config/consul/v2 v2.0.0-20230519061918-96480c11ee42
	github.com/go-kratos/kratos/contrib/config/etcd/v2 v2.0.0-20230519061918-96480c11ee42
	github.com/go-kratos/kratos/contrib/config/kubernetes/v2 v2.0.0-20230519061918-96480c11ee42
	github.com/go-kratos/kratos/contrib/config/nacos/v2 v2.0.0-20230519061918-96480c11ee42
	github.com/go-kratos/kratos/contrib/config/polaris/v2 v2.0.0-20230519061918-96480c11ee42
	github.com/go-kratos/kratos/contrib/log/aliyun/v2 v2.0.0-20230519061918-96480c11ee42
	github.com/go-kratos/kratos/contrib/log/fluent/v2 v2.0.0-20230519061918-96480c11ee42
	github.com/go-kratos/kratos/contrib/log/logrus/v2 v2.0.0-20230519061918-96480c11ee42
	github.com/go-kratos/kratos/contrib/log/tencent/v2 v2.0.0-20230519061918-96480c11ee42
	github.com/go-kratos/kratos/contrib/log/zap/v2 v2.0.0-20230516054017-1d50f502622a
	github.com/go-kratos/kratos/contrib/registry/consul/v2 v2.0.0-20230519061918-96480c11ee42
	github.com/go-kratos/kratos/contrib/registry/etcd/v2 v2.0.0-20230519061918-96480c11ee42
	github.com/go-kratos/kratos/contrib/registry/eureka/v2 v2.0.0-20230519061918-96480c11ee42
	github.com/go-kratos/kratos/contrib/registry/kubernetes/v2 v2.0.0-20230519061918-96480c11ee42
	github.com/go-kratos/kratos/contrib/registry/nacos/v2 v2.0.0-20230519061918-96480c11ee42
	github.com/go-kratos/kratos/contrib/registry/polaris/v2 v2.0.0-20230519061918-96480c11ee42
	github.com/go-kratos/kratos/contrib/registry/servicecomb/v2 v2.0.0-20230519061918-96480c11ee42
	github.com/go-kratos/kratos/contrib/registry/zookeeper/v2 v2.0.0-20230519061918-96480c11ee42
	github.com/go-kratos/kratos/v2 v2.6.2
	github.com/go-zookeeper/zk v1.0.3
	github.com/google/subcommands v1.2.0
	github.com/gorilla/handlers v1.5.1
	github.com/hashicorp/consul/api v1.20.0
	github.com/minio/minio-go/v7 v7.0.53
	github.com/nacos-group/nacos-sdk-go v1.1.4
	github.com/olekukonko/tablewriter v0.0.5
	github.com/polarismesh/polaris-go v1.4.3
	github.com/sirupsen/logrus v1.9.2
	github.com/spf13/cobra v1.7.0
	github.com/stretchr/testify v1.8.3
	go.etcd.io/etcd/client/v3 v3.5.9
	go.opentelemetry.io/otel v1.15.1
	go.opentelemetry.io/otel/exporters/jaeger v1.15.1
	go.opentelemetry.io/otel/exporters/zipkin v1.15.1
	go.opentelemetry.io/otel/sdk v1.15.1
	go.uber.org/zap v1.24.0
	golang.org/x/tools v0.7.0
	google.golang.org/grpc v1.54.0
	google.golang.org/protobuf v1.30.0
	k8s.io/client-go v0.27.2
)

require (
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.18 // indirect
	github.com/aliyun/aliyun-log-go-sdk v0.1.44 // indirect
	github.com/apolloconfig/agollo/v4 v4.3.0 // indirect
	github.com/armon/go-metrics v0.3.10 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cenkalti/backoff/v4 v4.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dlclark/regexp2 v1.7.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/emicklei/go-restful/v3 v3.9.0 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/felixge/httpsnoop v1.0.1 // indirect
	github.com/fluent/fluent-logger-golang v1.9.0 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/go-chassis/cari v0.6.0 // indirect
	github.com/go-chassis/foundation v0.4.0 // indirect
	github.com/go-chassis/openlog v1.1.3 // indirect
	github.com/go-errors/errors v1.0.1 // indirect
	github.com/go-kit/kit v0.10.0 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.1 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/go-playground/form/v4 v4.2.0 // indirect
	github.com/go-redis/redis/extra/rediscmd/v8 v8.11.5 // indirect
	github.com/go-redis/redis/extra/redisotel/v8 v8.11.5 // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/gofrs/uuid v4.2.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/gnostic v0.5.7-v3refs // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gorilla/websocket v1.4.3-0.20210424162022-e8629af678b7 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.2.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/karlseguin/ccache/v2 v2.0.8 // indirect
	github.com/klauspost/compress v1.16.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/lufia/plan9stats v0.0.0-20230110061619-bbe2e5e100de // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.4.3 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible // indirect
	github.com/openzipkin/zipkin-go v0.4.1 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/pelletier/go-toml/v2 v2.0.0-beta.8 // indirect
	github.com/philhofer/fwd v1.1.1 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/polarismesh/specification v1.2.1 // indirect
	github.com/power-devops/perfstat v0.0.0-20221212215047-62379fc7944b // indirect
	github.com/prometheus/client_golang v1.12.2 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.35.0 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/rs/xid v1.4.0 // indirect
	github.com/shirou/gopsutil/v3 v3.23.2 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/spf13/afero v1.9.2 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.11.0 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/tencentcloud/tencentcloud-cls-sdk-go v1.0.2 // indirect
	github.com/tinylib/msgp v1.1.6 // indirect
	github.com/tklauser/go-sysconf v0.3.11 // indirect
	github.com/tklauser/numcpus v0.6.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	go.etcd.io/etcd/api/v3 v3.5.9 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.9 // indirect
	go.opentelemetry.io/otel/trace v1.15.1 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/crypto v0.6.0 // indirect
	golang.org/x/mod v0.9.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/oauth2 v0.4.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/term v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	golang.org/x/time v0.1.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/api v0.27.2 // indirect
	k8s.io/apimachinery v0.27.2 // indirect
	k8s.io/klog/v2 v2.90.1 // indirect
	k8s.io/kube-openapi v0.0.0-20230501164219-8b0f38b5fd1f // indirect
	k8s.io/utils v0.0.0-20230209194617-a36077c30491 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)
