// Use panic and recover to write a function that contains
// no return statement yet return a non-zero value.
package main

import (
	"fmt"
	"log"
)

func division(x, y float64) (res float64, err error) {
	type nonzero struct{}
	defer func() {
		switch p := recover(); p {
		case nonzero{}:
			res = x / y
			err = nil
		default:
			panic(p)
		}
	}()
	if y != 0 {
		panic(nonzero{})
	} else {
		panic(fmt.Errorf("Division by zero: dividend %v, divisor: %v", x, y))
	}
}

type pairs struct {
	dividend, divisor float64
}

func main() {
	for _, v := range []pairs{{3.0, 2}, {3.0, 0}} {
		quotient, err := division(v.dividend, v.divisor)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("quotient:", quotient)
	}
}
