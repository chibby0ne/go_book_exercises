// Write an in-place function that squashes each run of adjacent Unicode spaces
// (see unicode.IsSpace) in a UTF-8-encoded []byte slice into a single ASCII
// space
package main

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

func squashSpace(slice []byte) []byte {
	prev := false
	builder := strings.Builder{}
	sliceAsString := string(slice)

	for i := 0; i < len(sliceAsString); {
		r, size := utf8.DecodeRuneInString(sliceAsString[i:])
		if unicode.IsSpace(r) {
			if prev {
				i += size
				continue
			}
			builder.WriteRune(r)
			prev = true
		} else {
			builder.WriteRune(r)
			prev = false
		}
		i += size
	}
	return []byte(builder.String())
}

func main() {
	slice := []byte("hello  \t\n元気. メリークリスマス")
	fmt.Println(slice)
	slice = squashSpace(slice)
	fmt.Println(string(slice))
}
