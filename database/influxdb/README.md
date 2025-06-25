# InfluxDB

### Docker部署

pull image

```bash
docker pull bitnami/influxdb:latest
```

#### 2.x

```bash
docker run -itd \
    --name influxdb2-server \
    -p 8086:8086 \
    -e INFLUXDB_HTTP_AUTH_ENABLED=true \
    -e INFLUXDB_ADMIN_USER=admin \
    -e INFLUXDB_ADMIN_USER_PASSWORD=123456789 \
    -e INFLUXDB_ADMIN_USER_TOKEN=admintoken123 \
    -e INFLUXDB_DB=my_database \
    bitnami/influxdb:2.7.11
```

create admin user sql script:

```sql
create user "admin" with password '123456789' with all privileges
```

管理后台: <http://localhost:8086/>

#### 3.x

```bash
docker run -itd \
    --name influxdb3-server \
    -p 8181:8181 \
    -e INFLUXDB_NODE_ID=0 \
    -e INFLUXDB_HTTP_PORT_NUMBER=8181 \
    -e INFLUXDB_HTTP_AUTH_ENABLED=true \
    -e INFLUXDB_CREATE_ADMIN_TOKEN=yes \
    -e INFLUXDB_DB=my_database \
    bitnami/influxdb:latest

docker run -itd \
  --name influxdb3-explorer \
  -p 8888:80 \
  -p 8889:8888 \
  quay.io/influxdb/influxdb3-explorer:latest \
  --mode=admin
```

这个版本分离出来一个管理后台 InfluxDB Explorer：<http://localhost:8888/>

在管理后台填写：`http://host.docker.internal:8181`
