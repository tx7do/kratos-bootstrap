# ClickHouse

## Docker部署

```bash
docker pull bitnami/clickhouse:latest

docker run -itd \
    --name clickhouse-server \
    --network=app-tier \
    -p 8123:8123 \
    -p 9000:9000 \
    -p 9004:9004 \
    -e ALLOW_EMPTY_PASSWORD=no \
    -e CLICKHOUSE_ADMIN_USER=default \
    -e CLICKHOUSE_ADMIN_PASSWORD=123456 \
    bitnami/clickhouse:latest
```
