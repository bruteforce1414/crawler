package queue

import (
	"github.com/streadway/amqp"
	"log"
)

type Publisher interface {
	Messages(linksPublishing []string)
}

type publisher struct {
	connection *amqp.Connection // тут ставишь нужный тип
}

func NewPublisher() (Publisher, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	_, err = ch.QueueDeclare(
		"publisher_queue", // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	failOnError(err, "Failed to declare an exchange")
	publisherObject := publisher{conn}
	return &publisherObject, err
}

func (c *publisher) Messages(linksPublishing []string) {
	//fmt.Println(linksPublishing)

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
