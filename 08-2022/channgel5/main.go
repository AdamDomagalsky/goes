package main

import (
	"fmt"
	"sync"
)

func main() {

	wg := &sync.WaitGroup{}
	channel := make(chan int)

	wg.Add(2)
	go func(ch chan int, wg *sync.WaitGroup) {
		fmt.Println(<-ch)
		wg.Done()
	}(channel, wg)

	go func(ch chan int, wg *sync.WaitGroup) {
		ch <- 42
		close(ch)

		wg.Done()
	}(channel, wg)

	wg.Wait()
	fmt.Println("EOF script")

}
