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

const MQURL = "amap://admin:admin@127.0.0.1:5672/miaosha"

type RabbitMQ struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	QueueName string
	Exchange  string
	Key       string
	MQurl     string
	sync.Mutex
}

func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	return &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, MQurl: MQURL}
}

func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalln("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

func NewRabbitMQSimple(queueName string) *RabbitMQ {
	rabbitmq := NewRabbitMQ(queueName, "", "")
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.MQurl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq!")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel!")
	return rabbitmq
}

func (r *RabbitMQ) PublishSimple(message string) error {
	r.Lock()
	defer r.Unlock()
	_, err := r.channel.QueueDeclare(r.QueueName, false, false, false, false, nil)
	if err != nil {
		return err
	}
	r.channel.Publish(r.Exchange, r.QueueName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	})
	return nil
}

func (r *RabbitMQ) ComSumeSimple(orderSerive services.IOrderService, productService services.IproductSerive) {
	q, err := r.channel.QueueDeclare(r.QueueName, false, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
	}
	r.channel.Qos(1, 0, false)
	msgs, err := r.channel.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
	}
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received ad message :%s", d.Body)
			message := &models.Message{}
			err := json.Unmarshal([]byte(d.Body), message)
			if err != nil {
				fmt.Println(err)
			}
			_, err = orderSerive.InsertOrderByMessage(message)
			if err != nil {
				fmt.Println(err)
			}

			d.Ack(false)

		}
	}()

	log.Printf("[*] Waiting for messages.To exit press CTRL+C")
	<-forever
}
