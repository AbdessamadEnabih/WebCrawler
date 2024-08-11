package main

import (
	"fmt"
	"strings"
	"sync"
	"github.com/gocolly/colly"
)

func crawl(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	collector := colly.NewCollector()

	extractLinks(collector)

	collector.Visit(url)

}

func extractLinks(c *colly.Collector) []string {
	var links []string
	c.OnHTML("a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.Contains(link, "https") {
			fmt.Println(link)
			links = append(links, link)
		}
	})
	return links
}

func main() {
	var wg sync.WaitGroup

	seedUrls := []string{
		"https://en.wikipedia.org/wiki/Albert_Stanley,_1st_Baron_Ashfield",
		// "https://en.wikipedia.org/wiki/Horsecar",
	}

	for _, url := range seedUrls {
		wg.Add(1)
		go crawl(url, &wg)
	}

	wg.Wait()

	fmt.Print("Done complete !!!!!")
}
