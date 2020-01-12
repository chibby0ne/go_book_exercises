// Write a function expand(s string, f func(string) string) string that replaces
// each substring "$foo" withing s by the text returned by f("foo").
package main

import (
	"fmt"
	"strings"
)

const VAR = "$foo"

func expand(s string, f func(string) string) string {
	var buf strings.Builder
	parts := strings.Split(s, VAR)
	if len(parts) == 1 {
		return s
	}
	for _, part := range parts {
		if part == "" {
			buf.WriteString(f(VAR[1:]))
		} else {
			buf.WriteString(part)
		}
	}
	return buf.String()
}

func itRocks(s string) string {
	return fmt.Sprintf("%s rocks", s)
}

func main() {
	s := "$foo = $foo"
	fmt.Println("before expansion:", s)
	new_s := expand(s, itRocks)
	fmt.Println("after expansion:", new_s)
}
