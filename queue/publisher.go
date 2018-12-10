package queue

import (
	"github.com/streadway/amqp"
	"log"
)

type Public interface {
	Messages(linksPublishing []string)
	Close()
}

func (s *public) Close() {
	s.connAMQP.Close()
}

type public struct {
	connAMQP *amqp.Connection // тут ставишь нужный тип
}

func NewPublic(queueName string) (Public, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare an exchange")
	publicObject := public{conn}
	return &publicObject, err
}

func (c *public) Messages(linksPublishing []string) {
	//fmt.Println(linksPublishing)

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
