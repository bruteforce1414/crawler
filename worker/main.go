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
	"time"
)

var (
	counter int
)

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Println("%s: %s", msg, err)
		os.Exit(0)
	}
}

func main() {

	msgs, err := queue.NewSubscribe("subscriber_queue")

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	public, err := queue.NewPublic("publisher_queue")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	go func(c queue.Public, ch <-chan amqp.Delivery) {
		for d := range msgs.Messages() {
			counter = counter + 1
			go worker(string(d.Body), c)
			log.Printf("Received a message: %s", d.Body)

			d.Ack(false)
		}
	}(public, msgs.Messages())

	time.Sleep(60 * time.Second)
	msgs.Close()
	fmt.Println("Общее количество обработанных ссылок ", counter)
}

func worker(url string, public queue.Public) {

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("адрес не доступен")
		return
	}
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	linksAll, err := urlprocessing.FindURLs(bodyString)

	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("LinksAll length:", len(linksAll))
	for i, link := range linksAll {

		link, err = urlprocessing.ParseUrl(string(url), link)
		fmt.Println("Cсылка №", i+1, ": ", link)
		if link != "" {
			public.Messages(link)
		}

	}

}
