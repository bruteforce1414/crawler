package main

import (
	"fmt"
	"github.com/bruteforce1414/crawler/client"
	"github.com/bruteforce1414/crawler/queue"
	"github.com/streadway/amqp"
	"github.com/willf/bloom"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var filter = bloom.New(1000000, 5)

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

	go func(с queue.Public, ch <-chan amqp.Delivery) {
		for d := range msgs.Messages() {
			if !filter.Test([]byte(string(d.Body))) {
				filter.Add([]byte(string(d.Body)))
				go worker(string(d.Body), с)
				log.Printf("Received a message: %s", d.Body)

			}
			d.Ack(false)
		}
	}(public, msgs.Messages())

	time.Sleep(30 * time.Second)
	msgs.Close()
}

func worker(url string, public queue.Public) {

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("адрес не доступен")

	}
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	linksAll, err := urlprocessing.FindURLs(bodyString)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("LinksAll length:", len(linksAll))
	for i, link := range linksAll {
		fmt.Println("Cсылка №", i+1, ": ", link)
		link, err = urlprocessing.ParseUrl(string(url), link)
		if link != "" {
			if !filter.Test([]byte(string(link))) {
				filter.Add([]byte(string(link)))
				public.Messages(link)
			}
		}
		if err != nil {
		}
	}

}
