# rabbitmq

docker相关命令：

你可以通过 `http://container-ip:15672` 在容器里面访问。当然，如果你需要的话，可以把这个端口映射出来，在主机上的`8080`端口通过浏览器访问管理界面。

```shell
$ docker run -d --hostname my-rabbit --name my-rabbit -p 5672:5672 -p 8080:15672 rabbitmq:3-management
```

启动后就可以在浏览器访问RabbitMQ的管理界面了,同时RabbitMQ的5672端口也会映射出来。

1. 在`http://localhost:8080/#/vhosts`添加先的vhost。
2. 在`http://localhost:8080/#/users`添加新的用户。
3. 把2中新建的用户，添加到1新建的vhost的权限中。

## 本项目的RabbitMQ配置

这边因为只随便模拟单商品秒杀情况，所以选择了简单的默认exchange（也就是exchange留空），这样在push消息的时候，选择的key就会自动推送到对应的名字的queue。

#### 怎么保证RabbitMQ在保证数据不丢失了？

这边主要利用的RabbitMQ的消息持久化：

#### queue的持久化

queue的持久化是通过durable=true来实现的。在代码中，声明队列的代码：

```go
	_, err := r.channel.QueueDeclare(r.QueueName, true, false, false, false, nil)
// 其中第二个参数true，就是希望保证队列的持久化
```

其中源码中，对这个参数解释如下：

> Durable and Non-Auto-Deleted queues will survive server restarts and remain
> when there are no remaining consumers or bindings.  Persistent publishings will
> be restored in this queue on server restart.  These queues are only able to be
> bound to durable exchanges.

#### 消息的持久化
如过将queue的持久化标识durable设置为true,则代表是一个持久的队列，那么在服务重启之后，也会存在，因为服务会把持久化的queue存放在硬盘上，当服务重启的时候，会重新什么之前被持久化的queue。队列是可以被持久化，但是里面的消息是否为持久化那还要看消息的持久化设置。也就是说，重启之前那个queue里面还没有发出去的消息的话，重启之后那队列里面是不是还存在原来的消息，这个就要取决于发生着在发送消息时对消息的设置了。
如果要在重启后保持消息的持久化必须设置消息是持久化的标识。

设置消息的持久化的话，我们可以需要在推送的消息结构里面控制：

```go
type Publishing struct {
...
	DeliveryMode    uint8     // Transient (0 or 1) or Persistent (2)
...
}
```

所以，这边相关推送消息的代码是：

```go
	r.channel.Publish(r.Exchange, r.QueueName, false, false, amqp.Publishing{
		ContentType:  "text/plain",
		Body:         []byte(message),
		DeliveryMode: amqp.Persistent,//这个常量的值是2
	})
```

#### exchange的持久化
上面阐述了队列的持久化和消息的持久化，如果不设置exchange的持久化对消息的可靠性来说没有什么影响，但是同样如果exchange不设置持久化，那么当broker服务重启之后，exchange将不复存在，那么既而发送方rabbitmq producer就无法正常发送消息。这里建议，同样设置exchange的持久化。exchange的持久化设置也特别简单。但是因为我们这边是使用的默认exchange，就不需要了。

### 进一步讨论
1. 将queue，exchange, message等都设置了持久化之后就能保证100%保证数据不丢失了嚒？
   答案是否定的。
   首先，从consumer端来说，如果这时autoAck=true，那么当consumer接收到相关消息之后，还没来得及处理就crash掉了，那么这样也算数据丢失，这种情况也好处理，只需将autoAck设置为false(方法定义如下)，然后在正确处理完消息之后进行手动ack。

   其次，关键的问题是消息在正确存入RabbitMQ之后，还需要有一段时间（这个时间很短，但不可忽视）才能存入磁盘之中，RabbitMQ并不是为每条消息都做fsync的处理，可能仅仅保存到cache中而不是物理磁盘上，在这段时间内RabbitMQ broker发生crash, 消息保存到cache但是还没来得及落盘，那么这些消息将会丢失。那么这个怎么解决呢？首先可以引入RabbitMQ的mirrored-queue即镜像队列，这个相当于配置了副本，当master在此特殊时间内crash掉，可以自动切换到slave，这样有效的保障了HA, 除非整个集群都挂掉，这样也不能完全的100%保障RabbitMQ不丢消息，但比没有mirrored-queue的要好很多，很多现实生产环境下都是配置了mirrored-queue的。还有要在producer引入事务机制或者Confirm机制来确保消息已经正确的发送至broker端，有关RabbitMQ的事务机制或者Confirm机制可以参考：RabbitMQ之消息确认机制（事务+Confirm）。后面会记下代码怎么在consumer端设置Ack。

2. 消息什么时候刷到磁盘？
   写入文件前会有一个Buffer,大小为1M,数据在写入文件时，首先会写入到这个Buffer，如果Buffer已满，则会将Buffer写入到文件（未必刷到磁盘）。
   有个固定的刷盘时间：25ms,也就是不管Buffer满不满，每个25ms，Buffer里的数据及未刷新到磁盘的文件内容必定会刷到磁盘。
   每次消息写入后，如果没有后续写入请求，则会直接将已写入的消息刷到磁盘：使用Erlang的receive x after 0实现，只要进程的信箱里没有消息，则产生一个timeout消息，而timeout会触发刷盘操作。

### consumer端设置Ack

执行一个任务可能需要花费几秒钟，你可能会担心如果一个消费者在执行任务过程中挂掉了。基于现在的代码，一旦RabbitMQ将消息分发给了消费者，就会从内存中删除。在这种情况下，如果杀死正在执行任务的消费者，会丢失正在处理的消息，也会丢失已经分发给这个消费者但尚未处理的消息。

但是，我们不想丢失任何任务，如果有一个消费者挂掉了，那么我们应该将分发给它的任务交付给另一个消费者去处理。

为了确保消息不会丢失，RabbitMQ支持消息应答。消费者发送一个消息应答，告诉RabbitMQ这个消息已经接收并且处理完毕了。RabbitMQ可以删除它了。

如果一个消费者挂掉却没有发送应答，RabbitMQ会理解为这个消息没有处理完全，然后交给另一个消费者去重新处理。这样，你就可以确认即使消费者偶尔挂掉也不会不丢失任何消息了。

没有任何消息超时限制；只有当消费者挂掉时，RabbitMQ才会重新投递。即使处理一条消息会花费很长的时间。

同时，RabbitMQ的任务分发机制以下2种：

- Round-robin（轮询分发）:在默认情况下，RabbitMQ将逐个发送消息到在序列中的下一个消费者(而不考虑每个任务的时长等等，且是提前一次性分配，并非一个一个分配)。平均每个消费者获得相同数量的消息。
- Fair dispatch（公平分发）:为限制RabbitMQ只发不超过1条的消息给同一个消费者。当消息处理完毕后，有了反馈，才会进行第二次发送。

而第二种机制，就需要按照前面所说的，必须关闭自动应答，改为手动应答。现在简单说一下代码里面是怎么设置的：

```go
	//消费者流控
	r.channel.Qos(
		1,
		0,
		false,
	)
```

这段代码限制RabbitMQ只发不超过1条的消息给同一个消费者。当消息处理完毕后，有了反馈，才会进行第二次发送。

```go
	//接收消息
	msgs, err := r.channel.Consume(
		q.Name, // queue
		"",     // consumer
		//下面的false指的是是否自动应答，需要关闭。
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
```

关闭自动应答后需要注意的是，一定得在消费完消息后应答对应消息：

```go
		for d := range msgs {
...
			d.Ack(false)
		}
```

这边注意，返回的msgs是chan，所以需要用range读取。