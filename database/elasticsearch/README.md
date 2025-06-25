# ElasticSearch

## 概念对比

| ES存储结构   | RDBMS存储结构 |
|----------|-----------|
| Index    | 表         |
| Document | 行         |
| Field    | 表字段       |
| Mapping  | 表结构定义     |

## mapping

- 动态映射（dynamic mapping）
- 显式映射（explicit mapping）
- 严格映射（strict mappings）

## Docker部署

### 拉取镜像

```bash
docker pull bitnami/elasticsearch:latest
```

### 启动容器

```bash
docker run -itd \
    --name elasticsearch \
    -p 9200:9200 \
    -p 9300:9300 \
    -e ELASTICSEARCH_USERNAME=elastic \
    -e ELASTICSEARCH_PASSWORD=elastic \
    -e xpack.security.enabled=true \
    -e discovery.type=single-node \
    -e http.cors.enabled=true \
    -e http.cors.allow-origin=http://localhost:13580,http://127.0.0.1:13580 \
    -e http.cors.allow-headers=X-Requested-With,X-Auth-Token,Content-Type,Content-Length,Authorization \
    -e http.cors.allow-credentials=true \
    bitnami/elasticsearch:latest
```

安装管理工具：

```bash
docker pull appbaseio/dejavu:latest

docker run -itd \
    --name dejavu-test \
    -p 13580:1358 \
    appbaseio/dejavu:latest
```

访问管理工具：<http://localhost:13580/>
