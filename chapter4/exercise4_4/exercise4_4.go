// Write a version of rotate that operates in a single pass

package main

import (
	"fmt"
)

func rotate(slice []int, positions int, left bool) {
	if left {
		temp := make([]int, positions, positions)
		copy(temp, slice[:positions])
		copy(slice[:len(slice)-positions], slice[positions:])
		copy(slice[len(slice)-positions:], temp)

	} else {
		temp := make([]int, positions, positions)
		copy(temp, slice[len(slice)-positions:])
		copy(slice[len(slice)-positions-1:], slice[:len(slice)-positions])
		copy(slice[:positions], temp)
	}
}

func main() {
	slice := []int{0, 1, 2, 3, 4}
	fmt.Println(slice)
	rotate(slice, 2, true)
	fmt.Println(slice)
	slice = []int{0, 1, 2, 3, 4}
	fmt.Println(slice)
	rotate(slice, 2, false)
	fmt.Println(slice)
}
