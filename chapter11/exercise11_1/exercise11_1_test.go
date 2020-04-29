package main

import (
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"
)

type Output struct {
	invalid int
	counts  map[rune]int
	utflen  [utf8.UTFMax + 1]int
}

func equalMaps(a, b map[rune]int) (bool, error) {
	if len(a) != len(b) {
		return false, fmt.Errorf("lengths of maps are different. len(got): %v and len(want): %v", len(a), len(b))
	}
	for k, v := range a {
		vv, ok := b[k]
		if !ok {
			return false, fmt.Errorf("char %c is not found in wants map", k)
		}
		if vv != v {
			return false, fmt.Errorf("count of char %c is %v but want: %v", k, v, vv)
		}
	}
	return true, nil
}

func equalSlicesInt(a, b []int) (bool, error) {
	if len(a) != len(b) {
		return false, fmt.Errorf("lengths of arrays are different. len(got): %v and len(want): %v", len(a), len(b))
	}
	for i, v := range a {
		if b[i] != v {
			return false, fmt.Errorf("count of chars of length %v is: %v but want: %v", i, v, b[i])
		}
	}
	return true, nil
}

func TestCountChars(t *testing.T) {
	var tests = []struct {
		input  string
		output Output
	}{
		{

			input: "asdfasdfhoho la;18123",
			output: Output{
				invalid: 0,
				counts:  map[rune]int{'a': 3, 's': 2, 'd': 2, 'f': 2, 'h': 2, 'o': 2, ' ': 1, 'l': 1, ';': 1, '1': 2, '8': 1, '2': 1, '3': 1},
				utflen:  [utf8.UTFMax + 1]int{0, 21, 0, 0, 0},
			},
		},
		{

			input: "^ADao la;18123",
			output: Output{
				invalid: 0,
				counts:  map[rune]int{'^': 1, 'A': 1, 'D': 1, 'a': 2, 'o': 1, ' ': 1, 'l': 1, ';': 1, '1': 2, '8': 1, '2': 1, '3': 1},
				utflen:  [utf8.UTFMax + 1]int{0, 14, 0, 0, 0},
			},
		},
	}
	for _, test := range tests {
		invalid, counts, utflen := CountChars(strings.NewReader(test.input))
		if invalid != test.output.invalid {
			t.Errorf("CountChars(%q) = invalid: %v, want %v", test.input, invalid, test.output.invalid)
		}
		if ok, err := equalMaps(counts, test.output.counts); !ok {
			t.Errorf("CountChars(%q) = %s ", test.input, err)
		}
		if ok, err := equalSlicesInt(utflen[:], test.output.utflen[:]); !ok {
			t.Errorf("CountChars(%q) = %s", test.input, err)
		}
	}
}
