package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	site_to_visit := []string{
		"https://en.wikipedia.org/wiki/Albert_Stanley,_1st_Baron_Ashfield",
		// "https://en.wikipedia.org/wiki/Horsecar",
	}

	crawl(site_to_visit)

	fmt.Print("Done Crawling !!!!!")
}

func crawl(site_to_visit []string) {
	collector := colly.NewCollector()

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	collector.OnResponse(func(r *colly.Response) {
		fmt.Println("Got a response from", r.Request.URL)
	})

	collector.OnHTML("a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.Contains(link, "https") {
			fmt.Println(link)
		}

		e.Request.Visit(e.Attr("href"))
	})

	collector.OnError(func(r *colly.Response, e error) {
		fmt.Println("Blimey, an error occurred!:", e)
	})

	for index, url := range site_to_visit {
		fmt.Printf("Website %d\n", index+1)
		collector.Visit(url)
		fmt.Print("----------------------------------")
	}
}
