package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
)

type Task struct {
	ID   string `json:"id"`
	BODY string `json:"body"`
}

// {\"id\": \"3\",\"body\":\"this is body\"}
func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	ctx := context.Background()

	for {
		x := rand.Intn(3) + 1
		time.Sleep(time.Second * time.Duration(x))
		taskJSON, err := rdb.RPop(ctx, "task_queue").Result()
		if err != nil {
			log.Printf("(%vs)Error dequeuing task: %v", x, err)
			continue // Retry dequeuing on error
		}

		fmt.Println(taskJSON)
		var task Task
		if err := json.Unmarshal([]byte(taskJSON), &task); err != nil {
			log.Println("Error unmarshalling task:", err)
			continue // Retry unmarshalling on error
		}
		fmt.Println(task.BODY)

		// Process the task
		err = func() error {
			// Do something with task
			return nil
		}()
		if err != nil {
			log.Println("Error processing task:", err)
			log.Println("retry TODO")
		} else {
			// Acknowledge successful task completion by removing it from the tracking hash
			// rdb.HDel(ctx, "task_hash", task.ID)
		}
	}
}
