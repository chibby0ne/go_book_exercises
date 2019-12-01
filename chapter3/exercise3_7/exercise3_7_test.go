package main

import (
	"testing"
)

var roundtests = []struct {
	in     float64
	out    float64
	digits int
}{
	{1.234528123, 1.2345, 4},
	{0.00012, 0.0000, 4},
	{8888.1234567891011, 8888.1235, 4},
	{0.999999, 1, 2},
	{1.999999, 2, 3},
	{0.999999, 1, 2},
	{0.999999, 1, 1},
}

func TestRound(t *testing.T) {
	for _, test := range roundtests {
		out := round(test.in, test.digits)
		if test.out != out {
			t.Errorf("do not match. expected %v, actual %v\n", test.out, out)
		}
	}
}
