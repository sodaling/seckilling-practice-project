# mysql

通过docker新建一个mysql容器并不复杂,这边我们秒杀的服务器的基础配置是。
```shell script
$ docker run --name miaosha-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 -d mysql
```

1. 新建miaosha数据库,CHARSET=utf8mb4 COLLATE=utf8_unicode_ci
2. 在同目录下的scheme下的sql文件导入到miaosha数据库中。

