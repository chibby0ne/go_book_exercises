package main

import (
	"github.com/chibby0ne/go_book_exercises/chapter2/popcount"
)

func PopCountShiftLoop(x uint64) int {
	count := 0
	for ; x != 0; x >>= 1 {
		if x&1 == 1 {
			count++
		}
	}
	return count
}

func main() {
	// Compare performace between the two implementations
	popcount.PopCount(123)
	PopCountShiftLoop(123)
}
