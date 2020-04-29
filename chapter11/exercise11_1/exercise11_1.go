// Write tests for the charcount program in Section 4.3
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"
)

func CountChars(r io.Reader) (int, map[rune]int, [utf8.UTFMax + 1]int) {
	counts := make(map[rune]int)    // counts of unicode chars
	var utflen [utf8.UTFMax + 1]int // count of lengths of utf-8 encodings
	invalid := 0                    // count of invalid utf-8 characters

	in := bufio.NewReader(r)
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
	}
	return invalid, counts, utflen
}

func main() {
	invalid, counts, utflen := CountChars(os.Stdin)
	fmt.Printf("rune\tcount\n")
	for c, n := range counts {
		fmt.Printf("%q\t%d\n", c, n)
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
