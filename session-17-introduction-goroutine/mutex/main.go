package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	x := 0
	var mutex sync.Mutex

	for i := 1; i <= 1000; i++ {
		go func() { //go 1, go 2, go 3, go 4
			for j := 1; j <= 100; j++ {
				mutex.Lock()
				x++ //variable x
				mutex.Unlock()
			}
		}()
	}

	time.Sleep(5 * time.Second)
	fmt.Println("Counter:", x)
}
