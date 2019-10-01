package main

import (
	"fmt"
	"github.com/gocolly/colly"
)

const (
	url = "https://sillyhatxu.com/"
)

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("sillyhatxu.com"),
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
	c.OnHTML("div[class=post-preview]", func(e *colly.HTMLElement) {
		href := e.ChildAttr("a", "href")
		title := e.ChildText("a")
		fmt.Println("---- OnHTML div -----")
		fmt.Println(fmt.Sprintf("title : %s ; href : %s", title, href))
	})

	//在OnXML回调之后调用
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})
	c.Visit(url)
}
