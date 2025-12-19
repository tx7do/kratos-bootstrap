::指定起始文件夹
set DIR=%cd%

cd %DIR%\api
go get all
go mod tidy

cd %DIR%\cache\redis
go get all
go mod tidy


cd %DIR%\tracer
go get all
go mod tidy


cd %DIR%\logger
go get all
go mod tidy

cd %DIR%\logger\aliyun
go get all
go mod tidy

cd %DIR%\logger\fluent
go get all
go mod tidy

cd %DIR%\logger\logrus
go get all
go mod tidy

cd %DIR%\logger\tencent
go get all
go mod tidy

cd %DIR%\logger\zap
go get all
go mod tidy

cd %DIR%\logger\zerolog
go get all
go mod tidy


cd %DIR%\config
go get all
go mod tidy

cd %DIR%\config\apollo
go get all
go mod tidy

cd %DIR%\config\consul
go get all
go mod tidy

cd %DIR%\config\etcd
go get all
go mod tidy

cd %DIR%\config\kubernetes
go get all
go mod tidy

cd %DIR%\config\nacos
go get all
go mod tidy

cd %DIR%\config\polaris
go get all
go mod tidy


cd %DIR%\registry
go get all
go mod tidy

cd %DIR%\registry\consul
go get all
go mod tidy

cd %DIR%\registry\etcd
go get all
go mod tidy

cd %DIR%\registry\eureka
go get all
go mod tidy

cd %DIR%\registry\kubernetes
go get all
go mod tidy

cd %DIR%\registry\nacos
go get all
go mod tidy

cd %DIR%\registry\polaris
go get all
go mod tidy

cd %DIR%\registry\servicecomb
go get all
go mod tidy

cd %DIR%\registry\zookeeper
go get all
go mod tidy


cd %DIR%\oss\minio
go get all
go mod tidy


cd %DIR%\database\cassandra
go get all
go mod tidy

cd %DIR%\database\clickhouse
go get all
go mod tidy

cd %DIR%\database\elasticsearch
go get all
go mod tidy

cd %DIR%\database\ent
go get all
go mod tidy

cd %DIR%\database\gorm
go get all
go mod tidy

cd %DIR%\database\influxdb
go get all
go mod tidy

cd %DIR%\database\mongodb
go get all
go mod tidy


cd %DIR%\rpc
go get all
go mod tidy

cd %DIR%\bootstrap
go get all
go mod tidy

cd %DIR%
go get all
go mod tidy
