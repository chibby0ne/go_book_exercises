// Extend TestSplit to use a table of inputs and expected outputs

package whatever

import (
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {
	var tests = []struct {
		s    string
		sep  string
		want int
	}{
		{"a:b:c", ":", 3},
		{"a b c d", " ", 4},
		{"100002000040000", " ", 1},
		{"100002000040000", "2", 2},
	}
	for _, test := range tests {
		words := strings.Split(test.s, test.sep)
		if got := len(words); got != test.want {
			t.Errorf("Split(%q, %q) returned %d words, want: %d", test.s, test.sep, got, test.want)
		}
	}
}
