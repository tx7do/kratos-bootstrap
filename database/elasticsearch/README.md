# ElasticSearch

## 什么是 ElasticSearch

ElasticSearch（简称 ES）是一个**基于 Lucene 构建的开源、分布式、高可扩展、近实时的搜索引擎与数据分析引擎**。

它不仅可以做**全文检索**，还支持结构化查询、聚合分析、地理位置查询、日志分析等能力，是目前业界最主流的搜索与大数据分析引擎。

- 开发语言：Java
- 协议：RESTful API（HTTP/JSON）
- 定位：分布式文档存储 + 搜索引擎 + 数据分析引擎
- 生态：ELK Stack（ElasticSearch、Logstash、Kibana）核心组件

## ElasticSearch 核心概念（与关系型数据库对比）

Elasticsearch 是**分布式**、**RESTful 风格**的搜索引擎，采用**文档型存储**，与传统关系型数据库（MySQL/Oracle）结构对比如下：

| Elasticsearch  存储结构	 | 关系型数据库（RDBMS）	 | 说明                 |
|----------------------|----------------|--------------------|
| Cluster（集群）          | 	数据库实例	        | 一个 ES 集群包含多个节点     |
| Node（节点）	            | 数据库服务	         | 单个 ES 服务进程         |
| Index（索引）            | 	表（Table）      | 	一类相似数据的集合         |
| Type（类型）	            | 表子分类           | 	ES 7.x 已移除，仅做历史了解 |
| Document（文档）	        | 行（Row）	        | 一条数据，JSON 格式存储     |
| Field（字段）	           | 表字段（Column）	   | 文档中的属性             |
| Mapping	             | 表结构定义（Schema）	 | 定义字段类型、分词器、是否索引等   |
| Shard（分片）	           | 分表 / 分区	       | 索引水平拆分，提升并发与存储     |
| Replica（副本）	         | 主从备份	          | 分片备份，保证高可用         |

## Mapping 详解（字段映射规则）

Mapping 是 ES 中定义索引结构的规则，决定字段如何存储、检索、分词。

### 1. 三大映射类型

#### （1）动态映射（Dynamic Mapping）

- **特点**：ES 自动识别数据类型，无需手动定义
- **适用场景**：测试环境、临时数据、结构不确定的数据
- 自动类型规则：
    - 字符串 → `text` + `keyword`
    - 数字 → `long`/`integer`/`float`
    - 布尔 → `boolean`
    - 日期 → `date`

#### （2）显式映射（Explicit Mapping）

- **特点**：手动创建索引并定义字段、类型、分词器
- **适用场景**：生产环境、需要精准控制的业务数据
- **优势**：性能更高、查询更精准、避免自动映射错误

#### （3）严格映射（Strict Mapping）

- **特点**：仅允许使用定义好的字段，插入未定义字段直接报错
- **适用场景**：严格数据规范、防止脏数据写入
- **配置**：`"dynamic": "strict"`

### 2. 常用字段数据类型

- 字符串：`text`（支持分词检索）、`keyword`（精确匹配）
- 数值：`long`、`integer`、`short`、`byte`、`double`、`float`
- 布尔：`boolean`
- 日期：`date`
- 对象：`object`、`nested`（嵌套数组对象）
- 地理：`geo_point`（经纬度）

## Docker 部署 Elasticsearch

### 1. 镜像选择（稳定版）

```bash
# 长期支持版（推荐生产）
docker pull elasticsearch:8.19.14

# 最新版
docker pull elasticsearch:9.3.3
```

### 2. 单节点部署命令（开发 / 测试）

```bash
docker run -d \
  --name elasticsearch \
  --restart=always \
  -p 9200:9200 \
  -p 9300:9300 \
  -e "discovery.type=single-node" \
  -e "ES_JAVA_OPTS=-Xms512m -Xmx512m" \
  -e "xpack.security.enabled=false" \
  -e "xpack.security.enrollment.enabled=false" \
  elasticsearch:8.19.14
```

### 4. 验证部署

```bash
# 访问服务
curl http://localhost:9200

# 出现以下信息即部署成功
{
  "name" : "xxxx",
  "cluster_name" : "elasticsearch",
  "cluster_uuid" : "xxxx",
  "version" : {
    "number" : "8.19.14"
  },
  "tagline" : "You Know, for Search"
}
```

## Elasticsearch 基础操作（RESTful API）

所有操作基于 **HTTP 请求**，支持 GET/PUT/POST/DELETE 方法。

### 1. 索引操作

#### （1）创建索引

```bash
PUT /my_index
{
  "settings": {
    "number_of_shards": 1,   # 主分片数
    "number_of_replicas": 0  # 副本数
  }
}
```

#### （2）查看索引

```bash
GET /my_index
GET /_all  # 查看所有索引
```

#### （3）删除索引

```bash
DELETE /my_index
```

### 2. 文档操作

#### （1）创建 / 更新文档

```bash
# 指定文档ID
PUT /my_index/_doc/1
{
  "name": "张三",
  "age": 25,
  "city": "北京",
  "create_time": "2025-01-01"
}

# 自动生成ID
POST /my_index/_doc
{...}
```

#### （2）查询文档

```bash
# 根据ID查询
GET /my_index/_doc/1

# 全量查询
GET /my_index/_search
{
  "query": {
    "match_all": {}
  }
}
```

#### （3）删除文档

```bash
DELETE /my_index/_doc/1
```

### 3. 条件查询（常用）

```bash
# 关键词匹配查询
GET /my_index/_search
{
  "query": {
    "match": {
      "name": "张三"
    }
  }
}

# 精确匹配
GET /my_index/_search
{
  "query": {
    "term": {
      "city.keyword": "北京"
    }
  }
}

# 范围查询
GET /my_index/_search
{
  "query": {
    "range": {
      "age": {
        "gte": 20,
        "lte": 30
      }
    }
  }
}
```

## ES SQL

- ES 从 **6.3+ 版本原生支持 SQL**
- 支持 **类 MySQL 标准 SQL** 语法
- 底层自动把 SQL 翻译成 ES DSL 查询
- 适合：快速查询、报表、数据分析、临时排查

### 核心概念对应（ES SQL ↔ MySQL）

| ES SQL        | 	MySQL     |
|---------------|------------|
| index（索引）     | 	table（表）  |
| document（文档）	 | row（行）     |
| field（字段）     | 	column（列） |
| _id	          | 主键         |

### 最常用 ES SQL 语法

#### 1. 查询全部数据

````sql
SELECT * FROM my_index;
````

#### 2. 条件查询 WHERE

````sql
SELECT name, age, city 
FROM my_index 
WHERE age > 20 AND city = '北京';
````

#### 3. 精确匹配（等价 term）

````sql
SELECT * FROM my_index 
WHERE city.keyword = '北京';
````

#### 4. 模糊查询 LIKE

````sql
SELECT * FROM my_index 
WHERE name LIKE '%张%';
````

#### 5. 分页 LIMIT

````sql
SELECT * FROM my_index 
LIMIT 10;

SELECT * FROM my_index 
LIMIT 20,10;  -- 第2页，每页10条
````

#### 6. 排序 ORDER BY

````sql
SELECT * FROM my_index 
ORDER BY age DESC, create_time ASC;
````

#### 7. 聚合查询 GROUP BY

````sql
-- 按城市统计人数
SELECT city.keyword, COUNT(*) AS cnt
FROM my_index
GROUP BY city.keyword;

-- 平均年龄
SELECT AVG(age) FROM my_index;

-- 最大/最小/求和
SELECT MAX(age), MIN(age), SUM(age) FROM my_index;
````

#### 8. HAVING 过滤聚合结果

````sql
SELECT city.keyword, COUNT(*) AS cnt
FROM my_index
GROUP BY city.keyword
HAVING cnt > 5;
````

#### 9. 时间范围查询

````sql
SELECT * FROM my_index
WHERE create_time >= '2024-01-01' 
  AND create_time < '2025-01-01';
````

#### 10. IN 查询

````sql
SELECT * FROM my_index
WHERE age IN (20,25,30);
````

#### 11. IS NULL / IS NOT NULL

````sql
SELECT * FROM my_index
WHERE name IS NOT NULL;
````

#### 12. AS 别名

````sql
SELECT name AS username, age AS user_age
FROM my_index;
````

### 执行 ES SQL 的 3 种方式

#### 1. REST API（最常用）

```bash
POST /_sql
{
  "query": "SELECT * FROM my_index WHERE age > 20"
}
```

#### 2. Kibana 开发控制台

```bash
POST /_sql?format=txt
{
  "query": "SELECT * FROM my_index"
}
```

#### 3. curl 命令行

```bash
curl -X POST -H "Content-Type: application/json" \
  http://localhost:9200/_sql \
  -d '{
    "query": "SELECT name,age FROM my_index LIMIT 5"
  }'
```

### 高级用法

#### 1. 查看 SQL 翻译成的 DSL

```bash
POST /_sql/translate
{
  "query": "SELECT * FROM my_index WHERE age > 20"
}
```

可以看到 ES 底层真正执行的查询语句。

#### 2. 分页深度查询（超过 10000 条）

```sql
SELECT * FROM my_index
ORDER BY _id
LIMIT 10000, 10;
```

#### 3. 多索引联合查询

```sql
SELECT * FROM index1, index2 
WHERE age > 20;
```

### ES SQL 支持的常用函数

- 数值：`AVG()`, `SUM()`, `MAX()`, `MIN()`, `COUNT()`
- 字符串：`LENGTH()`, `UPPER()`, `LOWER()`, `CONCAT()`
- 日期：`YEAR()`, `MONTH()`, `DAY()`, `NOW()`
- 逻辑：`CASE WHEN ... THEN ... END`

示例：

```sql
SELECT 
  name,
  YEAR(create_time) AS year,
  CASE WHEN age>30 THEN '中年' ELSE '青年' END AS age_type
FROM my_index;
```

### 注意事项

1. **text 字段必须加 .keyword 才能精确匹配、分组、排序**
    ```sql
    WHERE city.keyword = '北京'
    GROUP BY city.keyword
    ```
2. **不支持事务、JOIN、子查询有限支持**
3. **默认只返回 1000 条数据**，需要用 LIMIT 翻页
4. **聚合查询性能非常快**，比 MySQL 强很多
