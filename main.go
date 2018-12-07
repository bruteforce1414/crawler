package main

import (
	"fmt"
	"github.com/bruteforce1414/crawler/client"
	"github.com/bruteforce1414/crawler/queue"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Println("%s: %s", msg, err)
		os.Exit(0)
	}
}

func main() {

	go worker()
	fmt.Scanln()
}

func worker() {
	queue.Publicher()
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
				fmt.Println("Ошибка:", err)
				os.Exit(0)
			}
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			bodyString := string(bodyBytes)

			linksAll, err := urlprocessing.FindURLs(bodyString)
			if err != nil {
				log.Fatal(err)
			}
			for i, link := range linksAll {
				fmt.Println("Cсылка №", i+1, ": ", link)
				link, err = urlprocessing.ParseUrl(string(d.Body), link)
				if link != "" {
					// здесь будет добавление в очередь
					fmt.Println("Полный путь по ссылке №", i+1, ": ", link)
					queue.Subscriber(link)
				}
				if err != nil {
					fmt.Println("Ошибка для ссылки №", i+1, ": ", err)
				}

			}
			log.Printf("Done")
			d.Ack(false)

		}

	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
