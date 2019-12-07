// Write a non-recursive version of comma, using bytes.Buffer instead of string concatenation.
package main

import (
	"bytes"
	"fmt"
)

func main() {
	fmt.Println(comma("1234567"))
	fmt.Println(comma("123456"))
	fmt.Println(comma("12345"))
	fmt.Println(comma("1234"))
	fmt.Println(comma("123"))
	fmt.Println(comma("12"))
}

func comma(s string) string {
	var buf bytes.Buffer
	n := len(s)
	if n <= 3 {
		return s
	}
	// Write the part until the first comma
	res := n % 3
	if res != 0 {
		buf.WriteString(s[:res])
		buf.WriteRune(',')
		s = s[res:]
	}
	// Writes the part until the last comma
	var i int
	for i = res + 3; i < n; i += 3 {
		buf.WriteString(s[:i])
		buf.WriteRune(',')
	}
	// Write the part after the last comma
	buf.WriteString(s[i-3:])
	return buf.String()
}
