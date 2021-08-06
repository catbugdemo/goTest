package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
)

func main() {
	// 创建 collector
	c := colly.NewCollector()

	// 事件监听
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Printf("Link found: %q -> %s\n",e.Text,link)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting",r.URL)
	})

	c.Visit("http://go-colly.org/")
}
