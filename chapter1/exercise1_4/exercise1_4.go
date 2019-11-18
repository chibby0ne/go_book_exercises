package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	counts := make(map[string]int)
	filenames := make(map[string][]string)
	for _, filename := range os.Args[1:] {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "dup3: %v\n", err)
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			if line != "" {
				counts[line]++
				if filenames[line] == nil {
					filenames[line] = []string{filename}
				} else {
					filenames[line] = append(filenames[line], filename)
				}
			}
		}
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
			fmt.Printf("Found in these files: ")
			for _, filename := range filenames[line] {
				fmt.Printf("%s\t", filename)
			}
			fmt.Println()
		}
	}
}
