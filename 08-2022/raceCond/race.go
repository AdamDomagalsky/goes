package main

import (
	"fmt"
)

func main() {

	// go run -race .
	// race conditions is a situation when more than 1 process, thread (or go routine)
	// attempt to access the same shared resource at the same time
	// for example both go routines want to modify x

	var x, y int
	go func() {
		x = x * 4
		y += x + 2
	}()
	go func() {
		x += 1
	}()

	fmt.Println(y)
}
