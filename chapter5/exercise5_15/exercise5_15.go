// Write variadic funcions max and minx, analogous to sum. What shoulds these functions do when
// called with no arguments? Write variants that require at least one argument.
package main

import (
	"fmt"
	"log"
	"math"
)

func max(vals ...int) (int, error) {
	// It might be better to panic in these cases?
	if len(vals) == 0 {
		return 0, fmt.Errorf("max needs at least one value")
	}
	max := math.MinInt64
	for _, val := range vals {
		if val > max {
			max = val
		}
	}
	return max, nil
}

func min(vals ...int) (int, error) {
	// It might be better to panic in these cases?
	if len(vals) == 0 {
		return 0, fmt.Errorf("min needs at least one value")
	}
	min := math.MaxInt64
	for _, val := range vals {
		if val < min {
			min = val
		}
	}
	return min, nil
}

func main() {
	v := []int{1, 2, -3, 10, 123, 8}
	maxV, err := max(v...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("max of %v, is: %v\n", v, maxV)
	minV, err := min(v...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("min of %v, is: %v\n", v, minV)
}
