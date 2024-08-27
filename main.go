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

var visited = make(map[string]bool)
var maxDepthReached int = 0

// Declare a mutex to protect the `visited` map from concurrent access, ensuring thread safety.
var mu sync.Mutex

type LogLevel string

const (
	INFO  LogLevel = "INFO"
	ERROR LogLevel = "ERROR"
	APP   LogLevel = "APP"
)

func logger(level LogLevel, msg string) {
	LogFile := "log.log"

	if level == APP {
		LogFile = "WebCrawler.log"
	}

	file, err := os.Create(LogFile)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		return
	}

	defer file.Close()

	writer := bufio.NewWriter(file)
	log.SetOutput(writer)

	// Log the message
	switch level {
	case INFO:
		log.Println(msg)
	case ERROR:
		log.Printf("[ERROR] %s", msg)
	case APP:
		log.Printf("[CRAWL] %s", msg)
	}

	writer.Flush()
}

// `crawl` function initiates the web crawling process for a given URL. It first checks if the URL has already been visited;
// if not, it marks the URL as visited and then uses Colly to visit the URL and extract links. New URLs found are then
// passed to the `crawl` function recursively to continue the crawling process.
func crawl(url string, depth int, wg *sync.WaitGroup) {
	// Defer the decrement of the WaitGroup counter to ensure it is called upon function completion,
	// effectively signaling that one less goroutine is active.
	defer wg.Done()
	collector := colly.NewCollector()

	// Lock the mutex to protect the `visited` map from concurrent access.
	mu.Lock()
	if visited[url] {
		mu.Unlock()
		return
	}
	visited[url] = true
	// Unlock the mutex, allowing other goroutines to access the `visited` map.
	mu.Unlock()

	collector.Visit(url)
	logger(APP, "Crawled to "+url)

	if depth > 0 && maxDepthReached < 2 {
		Links := extractLinks(collector)
		for index, extractLink := range Links {
			fmt.Print("In")
			logger(APP, "Sublink "+string(index+1)+"of "+url)
			maxDepthReached++
			go crawl(extractLink, depth-1, wg)
		}
	}
}

// The function `extractLinks` uses Colly to extract and return all links with "https" from a webpage.
func extractLinks(c *colly.Collector) []string {
	var links []string

	// Set up an HTML element handler to find all anchor tags ('<a>') and extract their href attributes.
	// Only links starting with "https" are collected to filter out relative or invalid URLs.
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
	// Initialize a WaitGroup to ensure the main thread waits for all goroutines to complete their execution,
	// preventing the program from exiting prematurely and ensuring all web pages are fully crawled.
	var wg sync.WaitGroup

	seedUrls := []string{
		"https://en.wikipedia.org/wiki/Albert_Stanley,_1st_Baron_Ashfield",
		// "https://en.wikipedia.org/wiki/Horsecar",
	}

	for _, url := range seedUrls {
		maxDepthReached = 0
		logger(APP, "First Seed URL"+url)
		wg.Add(1)
		go crawl(url, 2, &wg)
	}

	// Wait for all goroutines to complete their execution. This ensures that the program does not exit
	// prematurely, allowing all web pages to be fully crawled and processed.
	wg.Wait()

	fmt.Print("Done complete !!!!!")
}
