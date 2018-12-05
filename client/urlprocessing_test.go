package urlprocessing_test

import (
	"fmt"
	"github.com/bruteforce1414/crawler/client"
	"github.com/bruteforce1414/simple-web-client-with-test/client"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

)

func TestHttpClient_Get(t *testing.T) {
	a := assert.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()
	tURLS := []string{
		"http://google.com", "http://usatu.com", server.URL + "/teachers/283208/",
	}
	mux.HandleFunc("/teachers/283208/", func(w http.ResponseWriter, r *http.Request) {
		var links string
		var s string
		var link string
		for _, link = range tURLS {
			links = links + "<a href=" + strconv.Quote(link) + ">" + link + "</a>"
		}
		s = fmt.Sprintf("<html><body>%s</body></html>", links)
		w.Write([]byte(s))
	})

	clientTest:=client.NewHttpClient(server.URL)
		resp, err := clientTest.Get("/teachers/283208/")
	if err != nil {
		panic(err)
	}
	fmt.Println("Ответ сервера: ", string(resp))
	for i, link := range urlprocessing.FindURLs(string(resp)) {
		fmt.Println("Ссылка №",i+1,": ", link)
	}
	a.Equal(tURLS, urlprocessing.FindURLs(string(resp)))
}




func TestHttpClient_Get2(t *testing.T) {
	a := assert.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	tURLS := []string{
		server.URL ,"/teachers/283208/","info",
	}
	mux.HandleFunc("/teachers/283208/", func(w http.ResponseWriter, r *http.Request) {
		var links string
		var s string
		var link string
		for _, link = range tURLS {
			links = links + "<a href=" + strconv.Quote(link) + ">" + link + "</a>"
		}
		s = fmt.Sprintf("<html><body>%s</body></html>", links)
		w.Write([]byte(s))
	})

	mux.HandleFunc("/teachers/283208/info", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello!!!"))
	})

	reqGetText:="/teachers/283208/"
	clientTest := client.NewHttpClient(server.URL)
	resp, err := clientTest.Get(reqGetText)
	if err != nil {
		panic(err)
	}
	fmt.Println("Ответ сервера: ", string(resp))
	for i, link := range urlprocessing.FindURLs(string(resp)) {
		fmt.Println("Ссылка №",i+1,": ", link)
		link=urlprocessing.ParseUrl(server.URL+reqGetText,link)
		fmt.Println("Полный путь по ссылке №",i+1,": ", link)
	}
	a.Equal(tURLS, urlprocessing.FindURLs(string(resp)))
	time.Sleep(100000)

}

