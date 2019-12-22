// Rewrite reverse to use an array pointer instead of a slice

package main

import (
	"fmt"
)

func reverse(s *[5]int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func main() {
	s := [5]int{0, 1, 2, 3, 4}
	fmt.Println(s)
	reverse(&s)
	fmt.Println(s)
}
