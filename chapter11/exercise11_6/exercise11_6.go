// Write benchmarks to compare the PopCount implementation in section 2.6.23
// with your solutions to exercise 2.4 and exercise 2.5. At what point does the
// table-based approach break even?
package main

import (
	"fmt"
)

var pc [256]byte

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

func PopCount(x uint64) int {
	return int(pc[byte(x>>(0*8))] +
		pc[byte(x>>(1*8))] +
		pc[byte(x>>(2*8))] +
		pc[byte(x>>(3*8))] +
		pc[byte(x>>(4*8))] +
		pc[byte(x>>(5*8))] +
		pc[byte(x>>(6*8))] +
		pc[byte(x>>(7*8))])
}

// Exercise 2.4
func PopCountShiftLoop(x uint64) int {
	count := 0
	for ; x != 0; x >>= 1 {
		if x&1 == 1 {
			count++
		}
	}
	return count
}

// The expression x & (x- 1) clears the rightmost non-zero bit of x. Write a
// version of PopCount that counts bits by using this fact, and assess its
// performance.
// Exercise 2.5
func PopCountClearingLoop(x uint64) int {
	count := 0
	for x != 0 {
		x = x & (x - 1)
		count += 1
	}
	return count
}

func main() {
	values := []uint64{
		81,
		123,
		437,
		412384,
		1000123,
		78956781234,
		123890491234,
		12389049123432141234,
		15986149086411123128,
	}

	for _, val := range values {
		fmt.Printf("%d = %d\n", val, PopCount(val))
	}
}
