package main

import (
	"fmt"
	"github.com/bruteforce1414/crawler/client"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"net/http"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	Publicher()
	go worker()
	fmt.Scanln()
}

func worker() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
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
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {

		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			resp, err := http.Get(string(d.Body))
			if err != nil {

				panic(err)
			}
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			bodyString := string(bodyBytes)

			linksAll := urlprocessing.FindURLs(bodyString)
			for i, link := range linksAll {

				fmt.Println("Cсылка №", i+1, ": ", link)
				link = urlprocessing.ParseUrl(string(d.Body), link)
				if link != "" {
					// здесь будет добавление в очередь
					fmt.Println("Полный путь по ссылке №", i+1, ": ", link)
				}
				if link == "" {
					fmt.Println("Ссылка №", i+1, "не относится к формату text/html", link)
				}
				if link == "error" {
					fmt.Println("Ссылка на ресурс не работает")
				}
			}
			log.Printf("Done")
			d.Ack(false)

		}

	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func Publicher() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := "http://shkola114.ru/index.php?option=com_content&view=article&id=431&Itemid=195"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)
}
