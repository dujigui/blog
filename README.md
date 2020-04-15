一个简单的博客系统

## 特色



## 架构

```text
browser -> gateway -> visitor -> services
                   -> admin   ->
```

### gateway

鉴权，日志，访问控制（频率，黑名单）

### visitor

游客，普通用户

### admin

管理员

### services

1. user
2. file
3. backup
4. comment
5. ip/location
6. log
7. post
8. tag
9. notice
10. stats

1. mysql
2. redis

## 安装

### mysql

```
// install mysql
docker pull mysql:5.7
docker run --name blog -e MYSQL_ROOT_PASSWORD=Aa123456 -d mysql:5.7

// create mysql user blog
mysql -u root -p
create user blog identified by 'Aa123456';
grant all privileges on * . * to 'blog'@'%';
flush privileges;
exit;

// create databases
mysql -u blog -p
create database blog character set utf8mb4 collate utf8mb4_unicode_ci;
exit;
```

2. 启动

首先准备好mysql，按照[go-sql-driver 文档说明](https://github.com/go-sql-driver/mysql#dsn-data-source-name) 将 dsn 设置为环境变量，例如 BLOG_DSN="blog:Aa123456@/blog?charset=utf8mb4&parseTime=True&loc=Local"

有两种运行方式，第一种使用打包好的 docker 镜像运行

第二种是直接编译运行，开发或者测试时需要添加环境变量 BLOG_DEBUG=true