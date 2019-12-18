package main

import (
	"fmt"
)

const (
	KB = string(iota) + "000"
	MB = KB + "000"
	GB = MB + "000"
	TB = GB + "000"
	PB = TB + "000"
	EB = PB + "000"
	ZB = EB + "000"
	YB = ZB + "000"
)

func main() {
	fmt.Printf("%v\n", MB)
}
