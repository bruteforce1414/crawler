package queue

import (
	"github.com/streadway/amqp"
)

type Subcribe interface {
	Messages() <-chan amqp.Delivery
	Close()
}

type subcribe struct {
	connAMQP *amqp.Connection
	msg      <-chan amqp.Delivery
}

func (s *subcribe) Close() {
	s.connAMQP.Close()
}

func (s *subcribe) Messages() <-chan amqp.Delivery {
	return s.msg
}

func NewSubscribe(queueName string) (Subcribe, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	failOnError(err, "Failed to declare a queue")
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	subscribeObject := subcribe{connAMQP: conn, msg: msgs}
	return &subscribeObject, err
}
