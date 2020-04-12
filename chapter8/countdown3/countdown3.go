package main

import (
	"fmt"
	"os"
	"time"
)

func launch() {
	fmt.Println("Liftoff!!!")
}

func main() {
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		abort <- struct{}{}
	}()
	tick := time.Tick(1 * time.Second)
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		select {
		case <-tick:
			// do nothing
		case <-abort:
			fmt.Println("Launch aborted!")
			return
		}
	}
	launch()

}
