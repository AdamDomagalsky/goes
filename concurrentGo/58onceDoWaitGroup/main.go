package main

import (
	"fmt"
	"sync"
)

func main() {
	var count int

	increment := func() {
		count++
	}
	var once sync.Once
	var increments sync.WaitGroup
	increments.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer increments.Done()
			once.Do(increment)
		}()
	}
	increments.Wait()
	fmt.Printf("Count is %d\n", count)

	//	Example 2 as I'm lazy to do separate file
	var count2 int
	increment2 := func() {
		fmt.Println("increment2 called")
		count2++
	}
	decrement2 := func() {
		fmt.Println("decrement2 called")
		count2--
	}
	var once2 sync.Once
	once2.Do(decrement2)
	once2.Do(increment2) // won't be called
	fmt.Printf("Count2: %d\n", count2)

	//// circular reference -- DEADLOCK
	//var onceA, onceB sync.Once
	//var initB func()
	//initA := func() { onceB.Do(initB) }
	//initB = func() { onceA.Do(initA) } // \1
	//onceA.Do(initA)                    // \2

	//// \1 This call can’t proceed until the call at \2 returns.
	////	This program will deadlock because the call to Do at \1 won’t proceed until the call to
	////Do at \2 exits—a classic example of a deadlock. For some, this may be slightly counterintuitive
	////since it appears as though we’re using sync.Once as intended to guard
	////against multiple initialization, but the only thing sync.Once guarantees is that your
	////functions are only called once. Sometimes this is done by deadlocking
}
