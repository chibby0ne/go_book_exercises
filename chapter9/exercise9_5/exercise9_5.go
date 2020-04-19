// Write a program with two goroutines that send messages back and forth over
// two unbuffered channels in ping-pong fashion. How many communications per
// second can the program sustain?
package main

import (
	"fmt"
	"time"
)

func main() {
	inputA, inputB := make(chan int), make(chan int)
	tick := time.NewTicker(10 * time.Second)
	var count uint64
	go func() {
		for {
			select {
			case v := <-inputA:
				inputB <- v
			}
		}
	}()
	go func() {
		for {
			select {
			case v := <-inputB:
				inputA <- v
				count++
			}
		}
	}()
	inputA <- 5
	<-tick.C
	fmt.Printf("\rNumber of communications per second in average was: %.3f", float64(count*2)/10)

}
