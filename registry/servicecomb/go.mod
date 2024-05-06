module github.com/tx7do/kratos-bootstrap/registry/servicecomb

go 1.22.0

toolchain go1.22.1

replace github.com/tx7do/kratos-bootstrap/api => ../../api

require (
	github.com/go-chassis/sc-client v0.7.0
	github.com/go-kratos/kratos/contrib/registry/servicecomb/v2 v2.0.0-20240504101732-d0d5761f9ca8
	github.com/go-kratos/kratos/v2 v2.7.3
	github.com/tx7do/kratos-bootstrap/api v0.0.2
)

require (
	github.com/cenkalti/backoff/v4 v4.1.1 // indirect
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/go-chassis/cari v0.6.0 // indirect
	github.com/go-chassis/foundation v0.4.0 // indirect
	github.com/go-chassis/openlog v1.1.3 // indirect
	github.com/gofrs/uuid v4.2.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/gorilla/websocket v1.4.3-0.20210424162022-e8629af678b7 // indirect
	github.com/karlseguin/ccache/v2 v2.0.8 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	google.golang.org/protobuf v1.34.0 // indirect
)
