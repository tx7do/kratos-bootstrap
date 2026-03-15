# Doris

## Docker部署

```bash
docker pull apache/doris:be-4.0.4
docker pull apache/doris:fe-4.0.4

docker network create --driver bridge --subnet=172.20.80.0/24 doris-network

mkdir -p /data/fe/{doris-meta,conf,log}
mkdir -p /data/be/{storage,conf,log}
chmod -R 777 /data/fe /data/be

docker run -itd \
    --name=doris-fe \
    --env FE_SERVERS="fe1:172.20.80.2:9010" \
    --env FE_ID=1 \
    -p 8030:8030 \
    -p 9030:9030 \
    -p 9010:9010 \
    -v /data/fe/doris-meta:/opt/apache-doris/fe/doris-meta \
    -v /data/fe/conf:/opt/apache-doris/fe/conf \
    -v /data/fe/log:/opt/apache-doris/fe/log \
    --network=doris-network \
    --ip=172.20.80.2 \
    apache/doris:fe-4.0.4

docker run -itd \
    --name=doris-be \
    --env FE_SERVERS="fe1:172.20.80.2:9010" \
    --env BE_ADDR="172.20.80.3:9050" \
    -p 8040:8040 \
    -p 9050:9050 \
    -v /data/be/storage:/opt/apache-doris/be/storage \
    -v /data/be/conf:/opt/apache-doris/be/conf \
    -v /data/be/log:/opt/apache-doris/be/log \
    --network=doris-network \
    --ip=172.20.80.3 \
    apache/doris:be-4.0.4
```
