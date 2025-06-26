# MongoDB

## 概念对比

| MongoDB存储结构 | RDBMS存储结构   |
|-------------|-------------|
| database    | database    |
| collection  | table       |
| document    | row         |
| field       | column      |
| index       | 索引          |
| primary key | primary key |

## Docker部署

下载镜像：

```bash
docker pull bitnami/mongodb:latest
docker pull bitnami/mongodb-exporter:latest
```

带密码安装：

```bash
docker run -itd \
    --name mongodb-server \
    -p 27017:27017 \
    -e MONGODB_ROOT_USER=root \
    -e MONGODB_ROOT_PASSWORD=123456 \
    -e MONGODB_USERNAME=test \
    -e MONGODB_PASSWORD=123456 \
    -e MONGODB_DATABASE=finances \
    bitnami/mongodb:latest
```

不带密码安装：

```bash
docker run -itd \
    --name mongodb-server \
    -p 27017:27017 \
    -e ALLOW_EMPTY_PASSWORD=yes \
    bitnami/mongodb:latest
```

有两点需要注意：

1. 如果需要映射数据卷，需要把本地路径的所有权改到1001：`sudo chown -R 1001:1001 data/db`，否则会报错：
   `‘mkdir: cannot create directory ‘/bitnami/mongodb’: Permission denied’`；
2. 从MongoDB 5.0开始，有些机器运行会报错：`Illegal instruction`，这是因为机器硬件不支持 **AVX 指令集** 的缘故，没办法，MongoDB降级吧。
