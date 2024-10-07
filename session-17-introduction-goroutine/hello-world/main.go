package main

import (
	"fmt"
	"time"
)

func RunHelloWorld() {
	fmt.Println("Hello World")
}

func main() {
	go RunHelloWorld()
	fmt.Println("Main function")
	time.Sleep(1 * time.Second)

}
