# mysql

通过docker新建一个mysql容器并不复杂,这边我们秒杀的服务器的基础配置是。
```shell script
$ docker run --name miaosha-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 -d mysql
```

1. 新建miaosha数据库,CHARSET=utf8mb4 COLLATE=utf8_unicode_ci
2. 在同目录下的scheme下的sql文件导入到miaosha数据库中。

上面是简单的单实例mysql启动方式。
因为在秒杀中，mysql的瓶颈一般在消化订单的写入部分，所以这边就不附上索引等等相关的内容了。而在秒杀中需要保证的是mysql的高可用，所以，这边对mysql做的是主从复制的配置。这边因为机器缺乏，所以还是用docker-compose的方式配置。

# 正文

## 主从复制的方式

`MySQL 5.6` 开始主从复制有两种方式：**基于日志**（`binlog`）和 **基于** `GTID`（**全局事务标示符**）。

本文只涉及基于日志 `binlog` 的 **主从配置**。

## 主从复制的流程



![img](https://user-gold-cdn.xitu.io/2018/7/2/1645b204db420d67?imageView2/0/w/1280/h/960/format/webp/ignore-error/1)



`MySQL` 同步操作通过 `3` 个线程实现，其基本步骤如下：

1. **主服务器** 将数据的更新记录到 **二进制日志**（`Binary log`）中，用于记录二进制日志事件，这一步由 **主库线程** 完成；
2. **从库** 将 **主库** 的 **二进制日志** 复制到本地的 **中继日志**（`Relay log`），这一步由 **从库** `I/O` **线程** 完成；
3. **从库** 读取 **中继日志** 中的 **事件**，将其重放到数据中，这一步由 **从库** `SQL` **线程** 完成。

## 主从模式的优点

#### 1. 负载均衡

通常情况下，会使用 **主服务器** 对数据进行 **更新**、**删除** 和 **新建** 等操作，而将 **查询** 工作落到 **从库** 头上。

#### 2. 异地容灾备份

可以将主服务器上的数据同步到 **异地从服务器** 上，极大地提高了 **数据安全性**。

#### 3. 高可用

数据库的复制功能实现了 **主服务器** 与 **从服务器间** 的数据同步，一旦主服务器出了 **故障**，从服务器立即担当起主服务器的角色，保障系统持续稳定运作。

#### 4. 高扩展性

**主从复制** 模式支持 `2` 种扩展方式:

- **scale-up**

向上扩展或者 **纵向扩展**，主要是提供比现在服务器 **性能更好** 的服务器，比如 **增加** `CPU` 和 **内存** 以及 **磁盘阵列**等，因为有多台服务器，所以可扩展性比单台更大。

- **scale-out**

向外扩展或者 **横向扩展**，是指增加 **服务器数量** 的扩展，这样主要能分散各个服务器的压力。

## 主从模式的缺点

#### 1. 成本增加

搭建主从肯定会增加成本，毕竟一台服务器和两台服务器的成本完全不同，另外由于主从必须要开启 **二进制日志**，所以也会造成额外的 **性能消耗**。

#### 2. 数据延迟

**从库** 从 **主库** 复制数据肯定是会有一定的 **数据延迟** 的。所以当刚插入就出现查询的情况，可能查询不出来。当然如果是插入者自己查询，那么可以直接从 **主库** 中查询出来，当然这个也是需要用代码来控制的。

#### 3. 写入更慢

**主从复制** 主要是针对 **读远大于写** 或者对 **数据备份实时性** 要求较高的系统中。因为 **主服务器** 在写中需要更多操作，而且 **只有一台** 可以写入的 **主库**，所以写入的压力并不能被分散。

## 主从复制的前提条件

1. 主从服务器 **操作系统版本** 和 **位数** 一致。
2. 主数据库和从数据库的 **版本** 要一致。
3. 主数据库和从数据库中的 **数据** 要一致。
4. **主数据库** 开启 **二进制日志**，主数据库和从数据库的 `server_id` 在局域网内必须 **唯一**。

## 具体配置

### 1. 环境准备

| 名称           | 版本号     |
| :------------- | :--------- |
| Docker         | 18.03.1-ce |
| Docker Compose | 1.21.1     |
| MySQL          | 5.7.17     |

### 2. 配置docker-compose.yml

docker-compose.yml

```
version: '2'
services:
  mysql-master:
    build:
      context: ./
      dockerfile: master/Dockerfile
    environment:
      - "MYSQL_ROOT_PASSWORD=123456"
      - "MYSQL_DATABASE=replicas_db"
    links:
      - mysql-slave
    ports:
      - "33065:3306"
    restart: always
    hostname: mysql-master
  mysql-slave:
    build:
      context: ./
      dockerfile: slave/Dockerfile
    environment:
      - "MYSQL_ROOT_PASSWORD=123456"
      - "MYSQL_DATABASE=replicas_db"
    ports:
      - "33066:3306"
    restart: always
    hostname: mysql-slave
复制代码
```

### 3. 主数据库配置

#### 3.1. 配置Dockerfile

Dockerfile

```
FROM mysql:5.7.17
MAINTAINER harrison
ADD ./master/my.cnf /etc/mysql/my.cnf
```

#### 3.2. 配置my.cnf文件

my.cnf

```
[mysqld]
## 设置server_id，一般设置为IP，注意要唯一
server_id=100  
## 复制过滤：也就是指定哪个数据库不用同步（mysql库一般不同步）
binlog-ignore-db=mysql  
## 开启二进制日志功能，可以随便取，最好有含义（关键就是这里了）
log-bin=replicas-mysql-bin  
## 为每个session分配的内存，在事务过程中用来存储二进制日志的缓存
binlog_cache_size=1M  
## 主从复制的格式（mixed,statement,row，默认格式是statement）
binlog_format=mixed  
## 二进制日志自动删除/过期的天数。默认值为0，表示不自动删除。
expire_logs_days=7  
## 跳过主从复制中遇到的所有错误或指定类型的错误，避免slave端复制中断。
## 如：1062错误是指一些主键重复，1032错误是因为主从数据库数据不一致
slave_skip_errors=1062
```

### 4. 从数据库配置

#### 4.1. 配置Dockerfile

Dockerfile

```
FROM mysql:5.7.17
MAINTAINER harrison
ADD ./slave/my.cnf /etc/mysql/my.cnf
```

#### 4.2. 配置my.cnf文件

```
[mysqld]
## 设置server_id，一般设置为IP，注意要唯一
server_id=101  
## 复制过滤：也就是指定哪个数据库不用同步（mysql库一般不同步）
binlog-ignore-db=mysql  
## 开启二进制日志功能，以备Slave作为其它Slave的Master时使用
log-bin=replicas-mysql-slave1-bin  
## 为每个session 分配的内存，在事务过程中用来存储二进制日志的缓存
binlog_cache_size=1M  
## 主从复制的格式（mixed,statement,row，默认格式是statement）
binlog_format=mixed  
## 二进制日志自动删除/过期的天数。默认值为0，表示不自动删除。
expire_logs_days=7  
## 跳过主从复制中遇到的所有错误或指定类型的错误，避免slave端复制中断。
## 如：1062错误是指一些主键重复，1032错误是因为主从数据库数据不一致
slave_skip_errors=1062  
## relay_log配置中继日志
relay_log=replicas-mysql-relay-bin  
## log_slave_updates表示slave将复制事件写进自己的二进制日志
log_slave_updates=1  
## 防止改变数据(除了特殊的线程)
read_only=1  
复制代码
```

## MySQL的复制类型

### 基于语句的复制

主服务器上面执行的语句在从服务器上面再执行一遍，在 `MySQL-3.23` 版本以后支持。

> 问题：时间上可能不完全同步造成偏差，执行语句的用户也可能是不同一个用户。

### 基于行的复制

把主服务器上面改变后的内容直接复制过去，而不关心到底改变该内容是由哪条语句引发的，在 `MySQL-5.0` 版本以后引入。

> 问题：比如一个工资表中有一万个用户，我们把每个用户的工资+1000，那么基于行的复制则要复制一万行的内容，由此造成的开销比较大，而基于语句的复制仅仅一条语句就可以了。

### 混合类型的复制

`MySQL` 默认使用 **基于语句的复制**，当 **基于语句的复制** 会引发问题的时候就会使用 **基于行的复制**，`MySQL` 会自动进行选择。

