package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

const (
	host = "https://www.ldoceonline.com"
	url  = "https://www.ldoceonline.com/dictionary/breeze"
)

func main() {
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("ldoceonline.com", "www.ldoceonline.com"),
	)

	//doc, _ := htmlquery.Parse(strings.NewReader(htmlPage))
	// On every a element which has href attribute call callback
	c.OnHTML("div[class=dictionary]", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
		for _, node := range e.DOM.Find("./div[class=dictentry]").Nodes {
			fmt.Println(node.Data)
		}
		fmt.Println("-------------")
		//if link, found := e.DOM.Find("/div[class=dictentry]").Attr("href"); found {
		//	e.Request.Visit(link)
		//}

		//div
		//class = "dictionary"
		//link := e.Attr("href")
		// Print link
		//fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		//c.Visit(e.Request.AbsoluteURL(link))
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit(url)
}
