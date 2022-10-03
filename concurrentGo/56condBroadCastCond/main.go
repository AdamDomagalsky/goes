package main

import "sync"

type Button struct {
	Clicked *sync.Cond
}

func main() {
	button := Button{Clicked: sync.NewCond(&sync.Mutex{})}

	subscribe := func(c *sync.Cond, fn func()) {
		var innerWG sync.WaitGroup
		innerWG.Add(1)
		go func() {
			innerWG.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			fn()
		}()
		innerWG.Wait()
	}

	var clickRegistered sync.WaitGroup
	clickRegistered.Add(9)
	subscribe(button.Clicked, func() {
		println("subscriber 1 got a signal!")
		clickRegistered.Done()
	})

	subscribe(button.Clicked, func() {
		println("subscriber 2 got a signal!")
		clickRegistered.Done()
	})
	subscribe(button.Clicked, func() {
		println("subscriber 3 got a signal!")
		clickRegistered.Done()
	})
	button.Clicked.Broadcast() // Were it not for the clickRegistered WaitGroup,
	//button.Clicked.Signal() will deadlock as only 1 subscriber will pick a Signal
	clickRegistered.Wait()
}
