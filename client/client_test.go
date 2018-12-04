package client_test

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/bruteforce1414/simple-web-client-with-test/client"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"strconv"
	"testing"
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

	clientTest := client.NewHttpClient(server.URL)
	resp, err := clientTest.Get("/teachers/283208/")
	if err != nil {
		panic(err)
	}
	fmt.Println("Ответ сервера: ", string(resp))
	for i, link := range FindURLs(string(resp)) {
		fmt.Println("Ссылка №", i+1, ": ", link)
	}
	a.Equal(tURLS, FindURLs(string(resp)))
}

func FindURLs(body string) []string {
	var returnedLinks []string
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		panic(err)
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		returnedLinks = append(returnedLinks, link)

	})
	return returnedLinks
}
