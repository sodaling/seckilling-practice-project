# rabbitmq

docker相关命令：

你可以通过 `http://container-ip:15672` 在容器里面访问。当然，如果你需要的话，可以把这个端口映射出来，在主机上的`8080`端口通过浏览器访问管理界面。

```shell
$ docker run -d --hostname my-rabbit --name some-rabbit -p 8080:15672 rabbitmq:3-management
```

启动后就可以在浏览器访问RabbitMQ的管理界面了。

1. 在`http://localhost:8080/#/vhosts`添加先的vhost。
2. 在`http://localhost:8080/#/users`添加新的用户。
3. 把2中新建的用户，添加到1新建的vhost的权限中。

