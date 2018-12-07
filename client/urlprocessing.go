package urlprocessing

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

func FindURLs(body string) ([]string, error) {

	var returnedLinks []string
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		returnedLinks = append(returnedLinks, link)

	})
	return returnedLinks, err
}

func ParseUrl(urlPage string, ctxUrl string) (string, error) {

	u, err := url.Parse(ctxUrl)
	if err != nil {
		return "", err
	}

	base, err := url.Parse(urlPage)
	if err != nil {
		return "", err
	}
	fullLink := (base.ResolveReference(u)).String()

	// Checking for http in first fourth characters in url
	if strings.Contains(fullLink[0:4], "http") != true {
		return "", fmt.Errorf("Ссылка не поддерживает протокол http или https")
	}

	extension := filepath.Ext(fullLink)
	if len(extension) == 0 {
		return fullLink, err
	}
	if (extension == "htm") || (extension == "html") {
		return fullLink, err
	}

	resp, err := http.Head(fullLink)
	if err != nil {
		return "", fmt.Errorf("Страница не доступна")
	}
	contentType := resp.Header.Get("Content-type")
	if strings.Contains(contentType, "text/html") {
		return fullLink, err
	}
	return "", fmt.Errorf("Расширение " + extension + " не входит в область поиска")

}
