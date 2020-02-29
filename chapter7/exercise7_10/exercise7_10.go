// The sort.Interface type can be adapted to other uses. Write a function
// IsPalindrome (s sort.Interface) bool that reports whether the sequnece s is
// a palindrom., in other words, reversing the seuqnece wou.ld not change it.
// Assume that the elements at indices i and j are equal if
// !s.Less(i, j) && !s.Less(j, i)

package exercise7_10

import (
	"sort"
)

type word string

func (w word) Len() int           { return len(w) }
func (w word) Less(i, j int) bool { return []rune(w)[i] < []rune(w)[j] }
func (w word) Swap(i, j int)      { []rune(w)[i], []rune(w)[j] = []rune(w)[j], []rune(w)[i] }

func IsPalindrome(s sort.Interface) bool {
	for i, j := 0, s.Len()-1; i < s.Len()/2; i, i = i+1, j-1 {
		if s.Less(i, j) || s.Less(j, i) {
			return false
		}
	}
	return true
}
