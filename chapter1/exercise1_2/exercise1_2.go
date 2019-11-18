package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	s, arg := "", ""
	index := 0
	fmt.Println("index\tvalue")
	for index, arg = range os.Args {
		s += strconv.Itoa(index) + "\t" + arg
		if index != len(os.Args)-1 {
			s += "\n"
		}
	}
	fmt.Println(s)
}
