// Write an in-place function to eliminate adjacent duplicates in a []string slice

package main

import (
	"fmt"
)

func removeDuplicates(strings []string) []string {
	prev := strings[0]
	i := 1
	for _, s := range strings[1:] {
		if s != prev {
			strings[i] = s
			prev = s
			i++
		}
	}
	return strings[:i]
}

func main() {
	slice := []string{"hello", "hello", "baby", "baby", "baby", "babies", "something", "hello", "hello"}
	fmt.Println(slice)
	slice = removeDuplicates(slice)
	fmt.Println(slice)
}
