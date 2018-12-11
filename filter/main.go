package main

import (
	"fmt"
	"github.com/bruteforce1414/crawler/queue"
	"github.com/streadway/amqp"
	"github.com/willf/bloom"
	"os"
	"time"
)

var (
	counter  int
	counter2 int
)

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Println("%s: %s", msg, err)
		os.Exit(0)
	}
}

func main() {

	msgs, err := queue.NewSubscribe("publisher_queue")

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	public, err := queue.NewPublic("subscriber_queue")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	filter := bloom.New(10000000, 5)
	go func(c queue.Public, ch <-chan amqp.Delivery) {
		for d := range msgs.Messages() {
			counter2 = counter2 + 1
			if !filter.Test([]byte(string(d.Body))) {
				filter.Add([]byte(string(d.Body)))
				c.Messages(string(d.Body))
				counter = counter + 1
			}
			d.Ack(false)
		}
	}(public, msgs.Messages())

	time.Sleep(60 * time.Second)
	msgs.Close()
	fmt.Println("Количество уникальных ссылок: ", counter)
	fmt.Println("Количество обработанных ссылок: ", counter2)
}
