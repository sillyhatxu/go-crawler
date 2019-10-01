package main

import (
	"fmt"
	"github.com/gocolly/colly"
)

const (
	//host = "https://sillyhatxu.com/"
	//url  = "https://www.ldoceonline.com/dictionary/breeze"
	//url  = "https://www.ldoceonline.com/dictionary/reserved"
	url = "https://sillyhatxu.com/"
)

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("sillyhatxu.com"),
		//colly.AllowedDomains("ldoceonline.com", "www.ldoceonline.com"),
	)
	//在请求之前调用
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	//如果在请求期间发生错误则调用
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	//收到回复后调用
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	//OnResponse如果接收到的内容是HTML ，则在之后调用
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		//e.Request.Visit(e.Attr("href"))
		fmt.Println("---- OnHTML a[href] -----")
		fmt.Println(e.Text)
	})
	c.OnHTML("tr td:nth-of-type(1)", func(e *colly.HTMLElement) {
		fmt.Println("---- OnHTML -----")
		fmt.Println("---- OnHTML tr td:nth-of-type(1) -----")
		fmt.Println(e.Text)
	})

	//OnHTML如果接收到的内容是HTML或XML ，则在之后调用
	c.OnXML("//h1", func(e *colly.XMLElement) {
		fmt.Println(e.Text)
	})

	//在OnXML回调之后调用
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})
	c.Visit(url)
}
