一个简单的博客系统

## 结构

```text
browser -> gateway -> visitor -> services
                   -> admin   ->
```

1. gateway 模块负责鉴权
2. visitor 模块负责普通用户页面逻辑
3. admin 模块负责管理后台页面逻辑
4. services 负责向上述模块提供数据库、日志和配置项等服务

services 模块包括：

1. users 用户账号体系
2. logs 日志服务
3. files 文件服务
4. db 数据库服务
5. posts/tags/comments 文章、标签和评论服务

## 使用 Docker 启动

### 创建网络

```shell
docker network create blog
```

### 启动 mysql

```shell
mkdir -p blog/mysql/conf

echo -e [client]\\n\
default-character-set=utf8mb4\\n\
[mysql]\\n\
default-character-set=utf8mb4\\n\
[mysqld]\\n\
collation-server=utf8mb4_unicode_ci\\n\
character-set-server=utf8mb4\\n\
max_allowed_packet=256M\\n\
default-time-zone=+08:00 > blog/mysql/conf/my.cnf

docker run -d \
-p 3306:3306 \
-e MYSQL_ROOT_PASSWORD=Aa123456 \
-e MYSQL_DATABASE=blog \
-v $PWD/blog/mysql/conf:/etc/mysql/conf.d \
--name mysql \
--network blog \
mysql:5.7
```

### 启动博客

```shell
mkdir -p blog/data

docker run -d \
-p 8080:8080 \
-e BLOG_DSN='root:Aa123456@tcp(mysql)/blog?charset=utf8mb4&parseTime=True&loc=Local' \
-v $PWD/blog/data:/app/data \
--name blog \
--network blog \
dujigui/blog:latest
```

通过 `docker logs -f blog` 命令可以看到容器已经正常启动并且监听 8080 端口。

### 初始化

启动后访问 `http://localhost:8080`，第一次启动时，会自动重定向至 `http://localhost:8080/init`，完成初始化即可。

网站首页为 `http://localhost:8080`，后台管理页面为 `http://localhost:8080/admin`。

## 直接编译启动

首先准备好 mysql，按照[go-sql-driver 文档说明](https://github.com/go-sql-driver/mysql#dsn-data-source-name) 将 dsn 设置为环境变量，例如 `BLOG_DSN="root:Aa123456@/blog?charset=utf8mb4&parseTime=True&loc=Local"``。

除此之外，开发或者测试时需要添加环境变量 `BLOG_DEBUG=true`。