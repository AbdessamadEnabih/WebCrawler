package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
	site_to_visit := []string{
		"https://en.wikipedia.org/wiki/Albert_Stanley,_1st_Baron_Ashfield",
		"https://en.wikipedia.org/wiki/Horsecar",
	}

	for _, url := range site_to_visit {
		collector := colly.NewCollector()

		// whenever the collector is about to make a new request
		collector.OnRequest(func(r *colly.Request) {
			// print the url of that request
			fmt.Println("Visiting", r.URL)
		})
		collector.OnResponse(func(r *colly.Response) {
			fmt.Println("Got a response from", r.Request.URL)
		})
		collector.OnError(func(r *colly.Response, e error) {
			fmt.Println("Blimey, an error occurred!:", e)
		})
		collector.Visit(url)
	}

}
