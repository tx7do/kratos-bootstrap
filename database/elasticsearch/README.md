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

```bash
docker pull bitnami/elasticsearch:latest

docker run -itd \
    --name elasticsearch \
    -p 9200:9200 \
    -p 9300:9300 \
    -e ELASTICSEARCH_USERNAME=elastic \
    -e ELASTICSEARCH_PASSWORD=elastic \
    -e ELASTICSEARCH_NODE_NAME=elasticsearch-node-1 \
    -e ELASTICSEARCH_CLUSTER_NAME=elasticsearch-cluster \
    bitnami/elasticsearch:latest
```
