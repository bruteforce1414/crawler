package urlprocessing

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

func FindURLs(body string) []string {
	var returnedLinks []string
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		panic(err)
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		returnedLinks = append(returnedLinks, link)

		fmt.Println("Нашлась ссылка:№", i+1, link)
	})
	return returnedLinks
}

func ParseUrl(urlPage string, ctxUrl string) string {
	var fullLink string
	u, err := url.Parse(ctxUrl)
	if err != nil {
		return "error"
	}
	base, err := url.Parse(urlPage)
	if err != nil {
		return "error"
	}
	fullLink = (base.ResolveReference(u)).String()

	// Checking for http in first fourth characters in url
	if strings.Contains(fullLink[0:4], "http") != true {
		fmt.Println("fullLink[0:4]", fullLink[0:4])
		return "Не http или https"
	}

	ext := filepath.Ext(fullLink)
	fmt.Println("Расширение файла:", ext)
	if len(ext) == 0 {
		return fullLink
	}
	if (ext == "htm") || (ext == "html") {
		return fullLink
	}

	resp, err := http.Head(fullLink)
	if err != nil {
		//	fmt.Println("Доступ к странице отсутствует, возвращена ошибка:", err)
		return "error"
	}
	contentType := resp.Header.Get("Content-type")
	if strings.Contains(contentType, "text/html") {
		fmt.Println("contentType", contentType)
		return fullLink
	}
	return ""

}
