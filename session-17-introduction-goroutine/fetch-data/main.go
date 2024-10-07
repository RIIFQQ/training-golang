package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func FetchData(url string, wg *sync.WaitGroup, ch chan<- string) {
	defer wg.Done()

	startTime := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprintf("Error fetching %s: %v", url, err)
		return

	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- fmt.Sprintf("Error reading response from %s: %v", url, err)
		return
	}

	duration := time.Since(startTime)
	ch <- fmt.Sprintf("Fetched %d bytes from %s in %v", len(body), url, duration)

}

func main() {
	urls := []string{
		"https://jsonplaceholder.typicode.com/todo/1",
		"https://httpbin.org/delay/2",
		"https://dog.ceo/api/breeds/image/random",
		"https://api.open-notify.org/iss-now.json",
		"https://randomuser.me/api/",
	}

	results := make(chan string, len(urls))

	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go FetchData(url, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}

}
