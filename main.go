package main

import (
	"fmt"
	"github.com/bruteforce1414/crawler/client"
	"github.com/bruteforce1414/crawler/queue"
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
	forever := make(chan bool)

	msgs, err := queue.Subscriber()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	consumer, err := queue.NewConsumer()
	go func(с queue.Consumer) {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			go worker(string(d.Body), с)
			d.Ack(false)
		}
	}(consumer)

	<-forever
}

func worker(url string, consumer queue.Consumer) {

	resp, err := http.Get("https://vedica.ru")
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
	for i, link := range linksAll {
		fmt.Println("Cсылка №", i+1, ": ", link)
		link, err = urlprocessing.ParseUrl(string(url), link)
		if link != "" {
			// здесь будет добавление в очередь
			fmt.Println("Полный путь по ссылке №", i+1, ": ", link)
			linksPublishing = append(linksPublishing, link)
		}
		if err != nil {
			fmt.Println("Ошибка для ссылки №", i+1, ": ", err)
		}
		consumer.Messages(linksPublishing)
	}

}
