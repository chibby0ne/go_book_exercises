package main

import (
	"flag"
	"fmt"
	"github.com/chibby0ne/go_book_exercises/chapter7/tempconv"
)

var temp = tempconv.CelsiusFlag("temp", 20.0, "the temperature")

func main() {
    flag.Parse()
    fmt.Println(*temp)
}
