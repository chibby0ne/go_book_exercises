// Modify reverse to reverse the characters of a []byte slice that reperesents a utf-8 encoded string, in place. Can you do it without allocation new memory?

package main

import (
	"fmt"
	"unicode/utf8"
)

func reverseSlice(slice []byte) {
	size := len(slice)
	for i := 0; i < size/2; i++ {
		slice[i], slice[size-i-1] = slice[size-i-1], slice[i]
	}
}

func reverseUtf8(slice []byte) {
	length := len(slice)
	for i := 0; i < length; {
		_, size := utf8.DecodeRune(slice[i:])
		reverseSlice(slice[i : i+size])
		i += size

	}
	reverseSlice(slice)
}

func main() {
	s := []byte("hello, 世界")
	fmt.Printf("%v\t%s\n", s, s)
	reverseUtf8(s)
	fmt.Printf("%v\t%s\n", s, s)
}
