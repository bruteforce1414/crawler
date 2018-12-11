package urlprocessing_test

import (
	"fmt"
	"github.com/bruteforce1414/crawler/client"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
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
	reqGetText := server.URL + "/teachers/283208/"
	clientTest := server.Client()
	resp, err := clientTest.Get(reqGetText)
	if err != nil {
		panic(err)
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	t.Log("Ответ сервера: ", bodyString)
	findLinks, err := urlprocessing.FindURLs(bodyString)
	for i, link := range findLinks {
		t.Log("Ссылка №", i+1, ": ", link)
	}
	a.Equal(tURLS, findLinks)
}

func TestHttpClient_Get2(t *testing.T) {
	a := assert.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	tURLS := []string{
		server.URL, "/teachers/283208/", "info", "1.doc", "skype:live:1d6db6b30d8c9a1e?call",
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

	reqGetText := server.URL + "/teachers/283208/"

	clientTest := server.Client()

	resp, err := clientTest.Get(reqGetText)
	if err != nil {
		panic(err)
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	t.Log("Ответ сервера: ", bodyString)
	findLinks, err := urlprocessing.FindURLs(bodyString)
	for i, link := range findLinks {
		t.Log("Ссылка №", i+1, ": ", link)
		link, err = urlprocessing.ParseUrl(reqGetText, link)
		t.Log("Полный путь по ссылке №", i+1, ": ", link)
	}
	a.Equal(tURLS, findLinks)
	time.Sleep(100000)

}

func TestHttpClient_Get3(t *testing.T) {
	a := assert.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	tURLS := []string{
		server.URL, "/teachers/283208/", "info", "skype:live:1d6db6b30d8c9a1e?call",
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

	reqGetText := server.URL + "/teachers/283208/"

	clientTest := server.Client()

	resp, err := clientTest.Get(reqGetText)
	if err != nil {
		panic(err)
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	t.Log("Ответ сервера: ", bodyString)

	findLinks, err := urlprocessing.FindURLs(bodyString)
	parsedLinks := make([]string, len(findLinks))
	parsedErrors := make([]error, len(findLinks))
	for i, link := range findLinks {
		t.Log("Ссылка №", i+1, ": ", link)
		parsedLinks[i], parsedErrors[i] = urlprocessing.ParseUrl(reqGetText, link)
		t.Log("Полный путь по ссылке №", i+1, ": ", parsedLinks[i])
		t.Log("Ошибка для пути №", i+1, ": ", parsedErrors[i])
	}
	a.Equal(tURLS, findLinks)
	// нужно проверить спарсенные ссылки!!!
	a.Equal("Ссылка не поддерживает протокол http или https", parsedErrors[3].Error())
	time.Sleep(100000)

}

func TestHttpClient_Get4(t *testing.T) {
	a := assert.New(t)

	link, _ := urlprocessing.ParseUrl("https://www.deviantart.com/", "michaelmaddox222/?offset=10#comments")

	a.Equal("https://www.deviantart.com/michaelmaddox222/?offset=10", link)
	// нужно проверить спарсенные ссылки!!!

}
