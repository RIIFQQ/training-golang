package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// shareData represents the data structure where the scraped data will be stored
type SharedData struct {
	mu   sync.Mutex
	data map[string]float64
}

// scrapeWebsite simulates  scarping financial data from a website
func ScrapeWebsite(url string, sharedData *SharedData, wg *sync.WaitGroup, r *rand.Rand) {
	defer wg.Done()

	// Simulate data scraping with a random delay
	time.Sleep(time.Duration(r.Intn(100)) * time.Millisecond)

	//simulatedthe data scraped (random stock price)
	scrapeData := r.Float64() * 1000

	//use Mutex to safely update the shared data structure
	sharedData.mu.Lock()
	sharedData.data[url] = scrapeData
	sharedData.mu.Unlock()

	fmt.Printf("Scraped data from %s: %f\n", url, scrapeData)
}

func main() {
	//seed the random number generator for simulating data scraping
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	//shared data structure for storing scraped financial data
	sharedData := &SharedData{
		data: make(map[string]float64),
	}

	//list of websites to scrape financial data from
	websites := []string{
		"https://www.finance.yahoo.com/",
		"https://www.investing.com/",
		"https://alphavantage.com/",
		"https://www.google.com/finance/",
		"https://www.nasdaq.com/",
		"https://www.bloomberg.com/",
		"https://www.moringstart.com/",
		"https://coinmarketcap.com/",
		"https://data.worldbank.org/",
		"https://www.quandl.com/",
	}

	//create a Waitgroup to manage the goroutines
	var wg sync.WaitGroup

	//start scrapping each website concurrently
	for _, url := range websites {
		wg.Add(1)
		go ScrapeWebsite(url, sharedData, &wg, r)
	}

	//wait for all goroutines to complete
	wg.Wait()

	//Dispalay the collected financial data
	fmt.Println("Collected Financial Data:")
	sharedData.mu.Lock()
	for site, value := range sharedData.data {
		fmt.Printf("%s: %f\n", site, value)
	}
	sharedData.mu.Unlock()
}
