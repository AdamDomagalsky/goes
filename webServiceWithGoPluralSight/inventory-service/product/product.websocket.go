package product

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"golang.org/x/net/websocket"
)

type message struct {
	Data string `json:"data"`
	Type string `json:"type"`
}

func productSocket(ws *websocket.Conn) {
	ticker := time.NewTicker(10 * time.Second)
	defer func() {
		fmt.Println("closing the connection")
		ws.Close()
		ticker.Stop()
	}()
	done := make(chan bool)
	fmt.Println("new websocket connection established")
	go func(c *websocket.Conn) {
		for {
			var msg message
			if err := websocket.JSON.Receive(ws, &msg); err != nil {
				log.Println(err)
				break
			}
			fmt.Printf("received message %s\n", msg.Data)
		}
		close(done)
	}(ws)

	top10productsCache, err := GetTop10Products()
	if err != nil {
		log.Println(err)
		return
	}
	for {
		select {
		case <-done:
			fmt.Println("connection was closed, lets break out of here")
			return
		case t := <-ticker.C:
			fmt.Println("Tick at", t)
			newProducts, err := GetTop10Products()
			if err != nil {
				log.Println(err)
				return
			}
			if reflect.DeepEqual(top10productsCache, newProducts) {
				fmt.Println("nothing new with top10Products")
				continue
			} else {
				fmt.Println("got something new...saving to cache and sending via ws")
				top10productsCache = newProducts
			}

			if err := websocket.JSON.Send(ws, top10productsCache); err != nil {
				log.Println(err)
				return
			}
		}
	}
}
