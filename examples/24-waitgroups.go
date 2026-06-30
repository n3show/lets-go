// To wait for multiple goroutines to finish, we can
// use a wait group.

package main

import (
	"fmt"
	"sync"
	"time"
)

// This is the function we'll run in every goroutine.
func worker(id int) {
	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Second)
	fmt.Printf("Worker %d done\n", id)
}

func main() {
	// A WaitGroup waits for a collection of goroutines
	// to finish. The main goroutine calls Add to set
	// the number of goroutines to wait for.
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)

		// Wrap the worker call in a closure that makes
		// sure to tell the WaitGroup that this worker is
		// done.
		go func() {
			defer wg.Done()
			worker(i)
		}()
	}

	// Block until the WaitGroup counter goes back to 0;
	// all the workers notified they're done.
	wg.Wait()
}
