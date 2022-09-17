package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Implement the dining philosopher’s problem with the following constraints/modifications.
// 1. There should be 5 philosophers sharing chopsticks, with one chopstick between each adjacent pair of philosophers.
// 2. Each philosopher should eat only 3 times (not in an infinite loop as we did in lecture)
// 3. The philosophers pick up the chopsticks in any order, not lowest-numbered first (which we did in lecture).
// 4. In order to eat, a philosopher must get permission from a host which executes in its own goroutine.
// 5. The host allows no more than 2 philosophers to eat concurrently.
// 6. Each philosopher is numbered, 1 through 5.
// 7. When a philosopher starts eating (after it has obtained necessary locks) it prints “starting to eat <number>” on a line by itself, where <number> is the number of the philosopher.
// 8. When a philosopher finishes eating (before it has released its locks) it prints “finishing eating <number>” on a line by itself, where <number> is the number of the philosopher.
// Submission: Upload your source code for the program.

type ChopStick struct{ sync.Mutex }

type Philosopher struct {
	number                        int
	leftChopstick, rightChopstick *ChopStick
}

func eat(p *Philosopher, wg *sync.WaitGroup) {
	p.leftChopstick.Lock()
	p.rightChopstick.Lock()
	fmt.Printf("start to eat %d\n", p.number)

	time.Sleep(time.Duration(rand.Int()%300) * time.Millisecond)

	fmt.Printf("finishing eating %d\n", p.number)
	p.leftChopstick.Unlock()
	p.rightChopstick.Unlock()
	wg.Done()
}

func main() {

	allChopSticks := make([]*ChopStick, 5)
	for i, _ := range allChopSticks {
		allChopSticks[i] = new(ChopStick)
	}

	allPhilosopher := make([]*Philosopher, 5)
	for i, _ := range allPhilosopher {
		allPhilosopher[i] = &Philosopher{
			number:         i + 1,
			leftChopstick:  allChopSticks[i],
			rightChopstick: allChopSticks[(i+1)%5],
		}
	}

	maxEatersAtOnce := 5
	atOnce := maxEatersAtOnce
	wg := &sync.WaitGroup{}
	randomOrder := generateRandomPerumationXtimesThenShuffle(5, 3)
	fmt.Println(randomOrder)
	for _, philoNumber := range randomOrder {
		if atOnce == 0 {
			wg.Wait()
			atOnce = maxEatersAtOnce
			fmt.Printf("---- atOnce: %d ----\n", atOnce)
		}
		wg.Add(1)
		go eat(allPhilosopher[philoNumber], wg)
		atOnce -= 1
	}
	wg.Wait()
	fmt.Println("end of program")
}

func generateRandomPerumationXtimesThenShuffle(n int, X int) []int {
	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(n)
	for i := 0; i < X; i++ {
		perm = append(perm, rand.Perm(n)...)
	}
	rand.Shuffle(len(perm), func(i, j int) {
		perm[i], perm[j] = perm[j], perm[i]
	})

	return perm
}
