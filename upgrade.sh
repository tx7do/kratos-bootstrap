#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(pwd)"

dirs=(
  "api"
  "cache/redis"
  "tracer"
  "logger"
  "logger/aliyun"
  "logger/fluent"
  "logger/logrus"
  "logger/tencent"
  "logger/zap"
  "logger/zerolog"
  "config"
  "config/apollo"
  "config/consul"
  "config/etcd"
  "config/kubernetes"
  "config/nacos"
  "config/polaris"
  "registry"
  "registry/consul"
  "registry/etcd"
  "registry/eureka"
  "registry/kubernetes"
  "registry/nacos"
  "registry/polaris"
  "registry/servicecomb"
  "registry/zookeeper"
  "oss/minio"
  "database/cassandra"
  "database/clickhouse"
  "database/elasticsearch"
  "database/ent"
  "database/gorm"
  "database/influxdb"
  "database/mongodb"
  "rpc"
  "bootstrap"
)

for d in "${dirs[@]}"; do
  target="$ROOT_DIR/$d"
  if [ -d "$target" ]; then
    printf "=> [%s] running go get and go mod tidy\n" "$d"
    (
      cd "$target"
      # 获取依赖（递归模块）并整理 go.mod
      go get -v ./...
      go mod tidy
    )
  else
    printf "-> skip `%s` (not found)\n" "$d"
  fi
done

# 最后在根目录再执行一次
printf "=> [%s] final go get and go mod tidy\n" "$ROOT_DIR"
(
  cd "$ROOT_DIR"
  go get -v ./...
  go mod tidy
)

printf "done\n"