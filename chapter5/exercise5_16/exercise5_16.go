// Write a variadic version of strings.Join
package main

import (
	"fmt"
	"strings"
)

func Join(sep string, a ...string) string {
	var builder strings.Builder
	var length = len(a)
	for _, v := range a[:length-1] {
		builder.WriteString(v)
		builder.WriteString(sep)
	}
	builder.WriteString(a[length-1])
	return builder.String()
}

func main() {
	sep := ","
	variadicJoin := Join(sep, "apples", "oranges", "lemons", "cranberries", "bananas")
	standardJoin := strings.Join([]string{"apples", "oranges", "lemons", "cranberries", "bananas"}, sep)
	fmt.Printf("variadicJoin: %v\n", variadicJoin)
	fmt.Printf("standardJoin: %v\n", standardJoin)
	if variadicJoin == standardJoin {
		fmt.Println("They are identical in behavior")
	}
}
