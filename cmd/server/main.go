package main

import (
	"fmt"
	"time"
)

func producer(ch chan<- int) {
	for i := 1; i <= 5; i++ {
		fmt.Println("Producing", i)
		ch <- i
		time.Sleep(200 * time.Millisecond)
	}
	close(ch) // signal no more values
}

func consumer(ch <-chan int) {
	for v := range ch {
		fmt.Println("Consuming", v)
		time.Sleep(300 * time.Millisecond)
	}
}

func main() {
	ch := make(chan int, 2) // buffered channel = shared queue
	go producer(ch)
	consumer(ch)
}
