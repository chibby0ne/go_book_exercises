// The expression x & (x- 1) clears the rightmost non-zero bit of x. Write a version of PopCount that counts bits by using this fact, and assess its performance.

package main

import (
	"github.com/chibby0ne/go_book_exercises/chapter2/popcount"
)

func PopCountClearingLoop(x uint64) int {
	count := 0
	for x != 0 {
		x = x & (x - 1)
		count += 1
	}
	return count
}

func main() {
	popcount.PopCount(123)
	PopCountClearingLoop(123)
}
