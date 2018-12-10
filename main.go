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
	forever := make(chan bool)

	msgs, err := queue.Subscriber()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	publisher, err := queue.NewPublisher()

	go func(с queue.Publisher, ch <-chan amqp.Delivery) {
		//	fmt.Println("Полученные сообщения", msgs)
		for d := range msgs {
			//			log.Printf("Received a message: %s", d.Body)
			go worker(string(d.Body), с)
			d.Ack(false)
		}
	}(publisher, msgs)
	time.Sleep(30 * time.Second)
	defer close()
	//<-forever
}

func worker(url string, publisher queue.Publisher) {

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
	for _, link := range linksAll {
		//	fmt.Println("Cсылка №", i+1, ": ", link)
		link, err = urlprocessing.ParseUrl(string(url), link)
		if link != "" {
			// здесь будет добавление в очередь
			//	fmt.Println("Полный путь по ссылке №", i+1, ": ", link)
			linksPublishing = append(linksPublishing, link)
		}
		if err != nil {
			//		fmt.Println("Ошибка для ссылки №", i+1, ": ", err)
		}
		publisher.Messages(linksPublishing)
	}

}
