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
	for i, link := range urlprocessing.FindURLs(bodyString) {
		t.Log("Ссылка №", i+1, ": ", link)
	}
	a.Equal(tURLS, urlprocessing.FindURLs(bodyString))
}

func TestHttpClient_Get2(t *testing.T) {
	a := assert.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	tURLS := []string{
		server.URL, "/teachers/283208/", "info",
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
	for i, link := range urlprocessing.FindURLs(bodyString) {
		t.Log("Ссылка №", i+1, ": ", link)
		link = urlprocessing.ParseUrl(reqGetText, link)
		t.Log("Полный путь по ссылке №", i+1, ": ", link)
	}
	a.Equal(tURLS, urlprocessing.FindURLs(bodyString))
	time.Sleep(100000)

}
