package main

import (
	"github.com/chibby0ne/go_book_exercises/chapter2/popcount"
)

func PopCountLoop(x uint64) int {
	count := 0
	for i := 0; i < 8; i++ {
		count += int(popcount.PC[byte(x>>(i*8))])
	}
	return count
}

func main() {
	// Compare performace between the two implementations
	popcount.PopCount(123)
	PopCountLoop(123)
}
