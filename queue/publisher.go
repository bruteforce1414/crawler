package queue

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type Consumer interface {
	Messages(linksPublishing []string)
}

type consumer struct {
	connection *amqp.Connection // тут ставишь нужный тип
}

func NewConsumer() (Consumer, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	_, err = ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")
	consumerObject := consumer{conn}
	return &consumerObject, err
}

func (c *consumer) Messages(linksPublishing []string) {
	fmt.Println(linksPublishing)

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
