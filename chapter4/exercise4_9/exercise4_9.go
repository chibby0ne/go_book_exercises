// Write a program wordfreq to report the frequency of each word in a input text file.
// Call input.Split(bufio.ScanWords) before the first call to Scan to break the input into words instead of lines.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("USAGE: wordfreq FILE")
		os.Exit(1)
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	counts := make(map[string]int)
	input := bufio.NewScanner(file)
	input.Split(bufio.ScanWords)
	for input.Scan() {
		counts[input.Text()]++
	}
	if input.Err() != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%-40vcounts\n", "word")
	for word, count := range counts {
		fmt.Printf("%-40v%v\n", word, count)
	}
}
