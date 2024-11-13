cd api
go get all
go mod tidy

cd ../cache/redis
go get all
go mod tidy

cd ../../config/apollo
go get all
go mod tidy

cd ../consul
go get all
go mod tidy

cd ../etcd
go get all
go mod tidy

cd ../kubernetes
go get all
go mod tidy

cd ../nacos
go get all
go mod tidy

cd ../polaris
go get all
go mod tidy

cd ../../logger/aliyun
go get all
go mod tidy

cd ../fluent
go get all
go mod tidy

cd ../logrus
go get all
go mod tidy

cd ../tencent
go get all
go mod tidy

cd ../zap
go get all
go mod tidy

cd ../../oss/minio
go get all
go mod tidy

cd ../../registry/consul
go get all
go mod tidy

cd ../etcd
go get all
go mod tidy

cd ../eureka
go get all
go mod tidy

cd ../kubernetes
go get all
go mod tidy

cd ../nacos
go get all
go mod tidy

cd ../polaris
go get all
go mod tidy

cd ../servicecomb
go get all
go mod tidy

cd ../zookeeper
go get all
go mod tidy

cd ../../
go get all
go mod tidy
