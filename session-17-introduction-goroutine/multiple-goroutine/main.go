package main

import (
	"fmt"
	"time"
)

func DisplayNumber(number int) {
	fmt.Println("Display", number)
}

func main() {
	start := time.Now()
	for i := 0; i < 100000; i++ {
		go DisplayNumber(i)
	}

	duration := time.Since(start)
	fmt.Printf("execution time: %v\n", duration)
	time.Sleep(3 * time.Second)
}
