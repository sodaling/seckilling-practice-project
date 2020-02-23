package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"seckilling-practice-project/models"
	"seckilling-practice-project/services"
	"sync"
)

// 连接信息
const _MQUrl = "amqp://soda:soda@127.0.0.1:5672/miaosha"

type RabbitMq struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// 队列名称
	QueueName string
	// 交换路由名称
	Exchange string
	// bind key 名称
	key   string
	MqUrl string
	sync.Mutex
}

// 创建结构实体
func NewRabbitMq(queueName string, exchange string, key string) *RabbitMq {
	return &RabbitMq{QueueName: queueName, Exchange: exchange, key: key, MqUrl: _MQUrl}
}

// 断开channel和connection
func (r *RabbitMq) Destroy() {
	r.channel.Close()
	r.conn.Close()
}

// 错误处理函数
func (r *RabbitMq) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("#{%s}:{%s}", message, err)
		panic(fmt.Sprintf("#{%s}:{%s}", message, err))
	}
}

// 简单模式下RabbitMQ实例
func NewRabbitMQSimple(queueName string) *RabbitMq {
	// 空字符串表示模式或者匿名的转发器。
	// 消息通过队列的routingKey路由到指定的队列中去，如果存在的话。
	// 这边默认转发给对应queue
	rabbitmq := NewRabbitMq(queueName, "", "")
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.MqUrl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

// 简单模式队列生产
func (r *RabbitMq) PublishSimple(message string) error {
	// channel并不是线程安全，这边加锁来使用
	r.Lock()
	defer r.Unlock()
	_, err := r.channel.QueueDeclare(r.QueueName, false, false, false, false, nil)
	if err != nil {
		return err
	}

	// 这边exchange是空
	r.channel.Publish(r.Exchange, r.QueueName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	})
	return nil
}

// simple模式下消费者
func (r *RabbitMq) ConsumeSimple(orderService services.IOrderService, productService services.IProductService) {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		false,
		false,
		//是否具有排他性
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	//消费者流控
	r.channel.Qos(
		1,
		0,
		false,
	)

	//接收消息
	msgs, err := r.channel.Consume(
		q.Name, // queue
		"",     // consumer
		//是否自动应答
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	//启用协程处理消息
	go func() {
		for d := range msgs {
			//消息逻辑处理，可以自行设计逻辑
			log.Printf("Received a message: %s", d.Body)
			message := &models.Message{}
			err := json.Unmarshal([]byte(d.Body), message)
			if err != nil {
				fmt.Println(err)
				return
			}
			//插入订单
			_, err = orderService.InsertOrderByMessage(message)
			if err != nil {
				fmt.Println(err)
				return
			}

			//扣除商品数量
			err = productService.SubNumberOne(message.ProductID)
			if err != nil {
				fmt.Println(err)
			}
			//如果为true表示确认所有未确认的消息，
			//为false表示确认当前消息
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
