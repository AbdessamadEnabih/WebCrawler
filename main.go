package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

var (
	visited         = make(map[string]bool)
	maxDepthReached int32
	mu              sync.RWMutex
)

type LogLevel string

const (
	INFO  LogLevel = "INFO"
	ERROR LogLevel = "ERROR"
	APP   LogLevel = "APP"
)

func logger(level LogLevel, msg string) {
	logFile := "log.log"
	if level == APP {
		logFile = "WebCrawler.log"
	}
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	log.SetOutput(writer)

	switch level {
	case INFO:
		log.Println(msg)
	case ERROR:
		log.Printf("[ERROR] %s", msg)
	case APP:
		log.Printf("[CRAWL] %s", msg)
	}
}

func crawl(url string, depth int, wg *sync.WaitGroup, c *colly.Collector) {
	defer wg.Done()

	mu.RLock()
	if visited[url] {
		mu.RUnlock()
		return
	}
	mu.RUnlock()

	mu.Lock()
	visited[url] = true
	mu.Unlock()

	err := c.Visit(url)
	if err != nil {
		logger(ERROR, fmt.Sprintf("Error visiting %s: %v", url, err))
		return
	}

	logger(APP, "Crawled to "+url)

	if depth > 0 && int(maxDepthReached) < 2 {
		links := extractLinks(c)
		for i, link := range links {
			logger(APP, fmt.Sprintf("Sublink %d of %s", i+1, url))
			maxDepthReached++
			wg.Add(1)
			go crawl(link, depth-1, wg, c.Clone())
		}
	}
}

func extractLinks(c *colly.Collector) []string {
	var links []string
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.HasPrefix(link, "https://") {
			links = append(links, link)
		}
	})
	return links
}

func main() {
	var wg sync.WaitGroup
	seedURLs := []string{
		"https://en.wikipedia.org/wiki/Albert_Stanley,_1st_Baron_Ashfield",
		// "https://en.wikipedia.org/wiki/Horsecar",
	}

	c := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(2),
	)

	for _, url := range seedURLs {
		maxDepthReached = 0
		logger(APP, "First Seed URL: "+url)
		wg.Add(1)
		go crawl(url, 2, &wg, c.Clone())
	}

	wg.Wait()
	fmt.Println("Crawling completed!")
}