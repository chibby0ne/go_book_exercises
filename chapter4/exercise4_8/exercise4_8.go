// Modify charcount to count letters, digits, and so on in ther Unicode categories, using functions like unicode.IsLetter

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"
)

func main() {
	counts := make(map[rune]int)       // counts of unicode chars
	categories := make(map[string]int) // counts of unicode chars categories
	var utflen [utf8.UTFMax + 1]int    // count of lengths of utf-8 encodings
	invalid := 0                       // count of invalid utf-8 characters

	in := bufio.NewReader(os.Stdin)
	for {
		r, n, err := in.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "charcount: %v\n", err)
			os.Exit(1)
		}
		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}
		counts[r]++
		utflen[n]++
		if unicode.IsLetter(r) {
			categories["letter"]++
		} else if unicode.IsDigit(r) {
			categories["digit"]++
		} else if unicode.IsControl(r) {
			categories["control"]++
		} else if unicode.IsPunct(r) {
			categories["punct"]++
		} else if unicode.IsSymbol(r) {
			categories["symbol"]++
		} else if unicode.IsMark(r) {
			categories["mark"]++
		} else if unicode.IsSpace(r) {
			categories["space"]++
		}
	}
	fmt.Printf("rune\tcount\n")
	for c, n := range counts {
		fmt.Printf("%q\t%d\n", c, n)
	}
	fmt.Printf("\ncategory\tcount\n")
	for category, n := range categories {
		fmt.Printf("%v\t\t%d\n", category, n)
	}
	fmt.Print("\nlen\tcount\n")
	for i, n := range utflen {
		if i > 0 {
			fmt.Printf("%d\t%d\n", i, n)
		}
	}
	if invalid > 0 {
		fmt.Printf("\n%d invalid utf-8 characters\n", invalid)
	}
}
