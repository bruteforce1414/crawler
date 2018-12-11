package queue

import (
	"github.com/streadway/amqp"
	"log"
)

type Public interface {
	Messages(linkPublishing string)
	Close()
}

func (s *public) Close() {
	s.connAMQP.Close()
}

type public struct {
	connAMQP *amqp.Connection // тут ставишь нужный тип
	msg      amqp.Queue
	channel  channelPublic
}

type channelPublic struct {
	channel *amqp.Channel
}

func NewPublic(queueName string) (Public, error) {
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
	failOnError(err, "Failed to declare an exchange")
	publicObject := public{
		connAMQP: conn,
		msg:      q,
		channel:  channelPublic{channel: ch},
	}
	return &publicObject, err
}

func (c *public) Messages(linkPublishing string) {

	err := c.channel.channel.Publish(
		"",         // exchange
		c.msg.Name, // routing key
		false,      // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(string(linkPublishing)),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", string(linkPublishing))

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
