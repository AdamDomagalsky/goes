package main

import (
	"fmt"
	"reflect"
	"sync"
)

func main() {
	var counter int // just to have a view
	myPool := &sync.Pool{
		New: func() interface{} {
			fmt.Printf("Creating new instance number %d\n", counter)
			counter++
			return struct{}{}
		},
	}

	myPool.Get()              // will create new instance1
	instance1 := myPool.Get() // will create new instance2
	fmt.Println(reflect.TypeOf(instance1))
	myPool.Put(instance1)
	myPool.Get() // won't create new instance as we return one via .Put

	fmt.Println("------calcPoolExample----")
	var numCalcsCreated int
	calcPool := &sync.Pool{New: func() interface{} {
		numCalcsCreated++
		mem := make([]byte, 1024)
		return &mem
	}}

	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	const numWorkers = 1024 * 1024
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := numWorkers; i > 0; i-- {
		go func() {
			defer wg.Done()
			mem := calcPool.Get().(*[]byte)
			defer calcPool.Put(mem)
			// Assume something interesting, but quick is being done with
			// this memory.
		}()

	}
	wg.Wait()
	fmt.Printf("%d calculators were created.\n", numCalcsCreated)

}
