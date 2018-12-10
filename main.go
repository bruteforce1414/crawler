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

	go func(с queue.Public, ch <-chan amqp.Delivery) {
		for d := range msgs.Messages() {
			log.Printf("Received a message: %s", d.Body)
			go worker(string(d.Body), с)
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
	var linksPublishing []string
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("LinksAll length:", len(linksAll))
	for i, link := range linksAll {
		fmt.Println("Cсылка №", i+1, ": ", link)
		link, err = urlprocessing.ParseUrl(string(url), link)
		if link != "" {
			linksPublishing = append(linksPublishing, link)
		}
		if err != nil {
		}
		public.Messages(linksPublishing)
	}

}
