package main

import (
	"fmt"
)

func appendInt(x []int, y ...int) []int {
	var z []int
	zlen := len(x) + len(y)
	if zlen <= cap(x) {
		// There is room to grow extend the slice.
		z = x[:zlen]
	} else {
		//  There is insufficient space Allocate a new array.
		// Grow by doubling, for amortized linear complexity
		zcap := zlen
		if zcap < 2*len(x) {
			zcap = 2 * len(x)
		}
		z = make([]int, zlen, zcap)
		copy(z, x)
	}
	copy(z[len(x):], y)
	return z
}

func main() {
	var z []int
	fmt.Printf("z is %v, len: %v cap: %v\n", z, len(z), cap(z))
	for i := 0; i < 10; i++ {
		z = appendInt(z, i)
		fmt.Printf("z is %v, len: %v cap: %v\n", z, len(z), cap(z))
	}
	z = appendInt(z, z...)
	fmt.Printf("z is %v, len: %v cap: %v\n", z, len(z), cap(z))
}
