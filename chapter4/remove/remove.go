package main

import (
	"fmt"
)

func remove(slice []int, i int) []int {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

func removeunordered(slice []int, i int) []int {
	slice[i] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}

func main() {
	s := []int{5, 6, 7, 8, 9}
	fmt.Println(s)
	fmt.Println(remove(s, 2))
	s = []int{5, 6, 7, 8, 9}
	fmt.Println(s)
	fmt.Println(removeunordered(s, 2))
}
