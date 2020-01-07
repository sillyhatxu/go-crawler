package springfieldspringfield

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

func CrawlerScript(url string) {
	c := colly.NewCollector(
		colly.Debugger(&debug.LogDebugger{}),
		colly.AllowedDomains("springfieldspringfield.co.uk", "www.springfieldspringfield.co.uk"),
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
	//<div id="content_container">
	c.OnHTML("#content_container", func(e *colly.HTMLElement) {
		//fmt.Println(e.Text)
		//<div class="main-content-left">
		e.ForEach("div[class=main-content-left]", func(i int, eLeft *colly.HTMLElement) {
			//<h1>This is title</h1>
			title := eLeft.ChildText("h1")
			fmt.Println("title:", title)
			eLeft.ForEach("div[class=season-episodes]", func(i int, eSeason *colly.HTMLElement) {
				season := eSeason.ChildText("h3")
				fmt.Println("season:", season)
				hrefs := eSeason.ChildAttrs("a", "href")
				for _, href := range hrefs {
					fmt.Println("href:", href)
				}
			})
		})
	})
	c.Visit(url)
}
