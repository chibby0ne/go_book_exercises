// Exercise 3.11: Enhance comma so that it deals correctly with floating-point numbers and an
// optional sign.
package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {
	fmt.Println(comma("1234567"))
	fmt.Println(comma("123456"))
	fmt.Println(comma("12345"))
	fmt.Println(comma("1234"))
	fmt.Println(comma("123"))
	fmt.Println(comma("12"))

	fmt.Println(comma("1234567.890"))
	fmt.Println(comma("123456.7890"))
	fmt.Println(comma("-1234567.890"))
	fmt.Println(comma("-123456.7890"))
	fmt.Println(comma("123.381"))
	fmt.Println(comma("-12"))

}

func comma(s string) (string, string) {
	var buf bytes.Buffer
	// Split string into integer part and decimal part
	parts := strings.Split(s, ".")
	integer := parts[0]
	n := len(integer)
	// Write the negative sign if any
	if integer[0] == '-' {
		n -= 1
		buf.WriteRune('-')
		integer = integer[1:]
	}
	if n <= 3 {
		return s, s
	}
	// Write the part until the first comma
	res := n % 3
	if res != 0 {
		buf.WriteString(integer[:res])
		buf.WriteRune(',')
	}
	// Writes the part until the last comma
	var i, prev int
	prev = res
	for i = res + 3; i < n; i += 3 {
		buf.WriteString(integer[prev:i])
		buf.WriteRune(',')
		prev = i
	}
	// Write the part after the last comma
	buf.WriteString(integer[i-3:])
	// Write the decimal decimal part including decimal point if there was any
	if len(parts) > 1 {
		buf.WriteRune('.')
		buf.WriteString(parts[1])
	}
	return s, buf.String()
}
