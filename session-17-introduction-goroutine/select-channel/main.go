package main

import (
	"time"
)

func GiveMeResponse(channel chan string) {
	time.Sleep(2 * time.Second)
	channel <- "Hello World"
}

func main() {
	channel1 := make(chan string)
	channel2 := make(chan string)
	defer close(channel1)
	defer close(channel2)

	go GiveMeResponse(channel1)
	go GiveMeResponse(channel2)

	counter := 0
	for {
		select {
		case data := <-channel1:
			counter++
			println("Data from channel 1:", data)
		case data := <-channel2:
			counter++
			println("Data from channel 2:", data)
		}

		if counter == 2 {
			break
		}
	}
}
