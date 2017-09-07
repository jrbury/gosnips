package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// create a context that can be canceled. This allows the parent thread to send a signal
	// to all the child threads that have this context object
	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)                    // create a channel that will be used with our notify hook below
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM) // notify on sig int and sig term

	var wg *sync.WaitGroup // create a waitgroup so we can keep track of our background threads

	// create a background thread, passing it the context object so it can be canceled
	// as well as the wait group so it can tell us when it has exited
	wg.Add(1) // Add one to the wait group for every thread we make
	go thread(ctx, "1", wg)

	wg.Add(1) // Add one to the wait group for every thread we make
	go thread(ctx, "2", wg)

	wg.Add(1) // Add one to the wait group for every thread we make
	go thread(ctx, "3", wg)

	<-signals // block here waiting for a ctrl+c
	cancel()  // cancel the background threads

	// now that we've notified all the background threads with the cancel
	// we wait for all of them to exit
	wg.Wait()

	log.Println("All threads have finished exiting.")
}

func thread(ctx context.Context, name string, wg *sync.WaitGroup) {
	defer wg.Done() // once this functin returns, let the wait group know we're done

	x := 0
	for {
		select {
		case <-ctx.Done(): // the parent context was canceled which means we want to shut down
			// here you could do any cleanup necessary before exiting
			log.Printf("Thread %s: received cancel, exiting thread.", name)
			return
		default:
			log.Printf("Thread %s: %d", name, x)
			x++
			time.Sleep(time.Second * time.Duration(rand.Intn(5))) // sleep for a random amount between 1-5 seconds
		}
	}
}
