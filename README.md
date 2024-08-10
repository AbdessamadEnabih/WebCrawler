# Concurrent Web Crawler in Go

This is a simple concurrent web crawler written in Go. The web crawler starts from a list of seed URLs, fetches their content, extracts links, and then recursively crawls those links. The project leverages Go's goroutines and channels to efficiently handle multiple web pages concurrently.

## Features

- **Concurrency with Goroutines**: Fetch and process multiple web pages simultaneously.
- **Synchronization with Channels**: Manage a queue of URLs to be visited and avoid race conditions.
- **Link Extraction**: Extract and follow links from the HTML content of web pages.
- **URL Deduplication**: Ensure each URL is only crawled once to avoid redundant work.
- **Depth-Limited Crawling**: Configure how deep the crawler should go from the initial seed URLs.
- **Rate Limiting**: Control the number of requests per second to avoid overwhelming servers.

## Requirements

- Go 1.16 or higher

## Installation

Clone the repository:

```bash
git clone https://github.com/yourusername/concurrent-web-crawler.git
cd concurrent-web-crawler
