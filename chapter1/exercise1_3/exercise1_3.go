package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("hello world")
}

func using_range_concat(slice []string) {
	s, sep := "", ""
	for _, arg := range slice {
		s += sep + arg
		sep = " "
	}
	fmt.Println(s)
}

func using_strings_join(slice []string) {
	fmt.Println(strings.Join(slice, " "))
}
