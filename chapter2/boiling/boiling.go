package main

import (
	"fmt"
)

const boilingF = 212.0

func main() {
	var f = boilingF
	var c = (f - 32) * 5 / 9
	fmt.Printf("boling point = %gF or %gC\n", f, c)
}
