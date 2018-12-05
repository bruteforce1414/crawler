package urlprocessing

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/url"
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

	})
	return returnedLinks
}

func ParseUrl(urlPage string, ctxUrl string) string {
	var fullLink string
	u, err := url.Parse(ctxUrl)
	if err != nil {
		log.Fatal(err)
	}
	base, err := url.Parse(urlPage)
	if err != nil {
		log.Fatal(err)
	}
	fullLink = (base.ResolveReference(u)).String()
	return fullLink
}
