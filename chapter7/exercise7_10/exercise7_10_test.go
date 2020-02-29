package exercise7_10

import (
	"testing"
)

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		word     word
		expected bool
	}{
		{word("11/11/11"), true},
		{word("11:11"), true},
		{word("madam"), true},
		{word("racecar"), true},
		{word("amanaplanacanalpanama"), true},
		{word("A man, a plan, a canal, Panama!"), false},
		{word("simple"), false},
		{word("word"), false},
		{word("a good old sentence"), false},
	}
	for _, tt := range tests {
		tt := tt // NOTE: https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(string(tt.word), func(t *testing.T) {
			t.Parallel() // marks each test case as capable of running in parallel with each other
			if actual := IsPalindrome(tt.word); actual != tt.expected {
				t.Errorf("expected: %v, actual: %v", tt.expected, actual)
			}
		})
	}
}
