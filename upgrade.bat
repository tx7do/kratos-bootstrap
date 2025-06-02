cd api
go get all
go mod tidy

cd ../utils
go get all
go mod tidy

cd ../cache/redis
go get all
go mod tidy

cd ../../logger
go get all
go mod tidy


cd ../tracer
go get all
go mod tidy

cd ../config
go get all
go mod tidy


cd ../registry
go get all
go mod tidy

cd consul
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


cd ../../../oss/minio
go get all
go mod tidy

cd ../../database/cassandra
go get all
go mod tidy

cd ../clickhouse
go get all
go mod tidy

cd ../ent
go get all
go mod tidy

cd ../gorm
go get all
go mod tidy

cd ../influxdb
go get all
go mod tidy

cd ../mongodb
go get all
go mod tidy


cd ../../rpc
go get all
go mod tidy

cd ../bootstrap
go get all
go mod tidy

cd ../
go get all
go mod tidy
