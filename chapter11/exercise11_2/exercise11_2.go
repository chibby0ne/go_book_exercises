// Write a set of tests for IntSet that checks that its behavior after each
// operation is equivalent to a set based on built-in maps. Save your
// implementation for benchmarking in Exercise 11.7

package main

import (
	"bytes"
	"fmt"
	"sort"
)

type Set interface {
	Has(x int) bool
	Add(x int)
	UnionWith(s Set)
	GetSet() []int
}

// A MapSet is a set of small non-negative integers.
// Its zero value represents the empty set.
type MapSet struct {
	words map[uint]bool
}

// Has reports whether the set contains the non-negative value x.
func (s *MapSet) Has(x int) bool {
	_, ok := s.words[uint(x)]
	return ok
}

// Add adds the non-negative value x to the set.
func (s *MapSet) Add(x int) {
	if s.words == nil {
		s.words = make(map[uint]bool)
	}
	s.words[uint(x)] = true
}

// UnionWith sets s to the union of s and t.
func (s *MapSet) UnionWith(t Set) {
	for _, word := range t.GetSet() {
		s.words[uint(word)] = true
	}
}

// GetsSet returns all elements in the set
func (s *MapSet) GetSet() []int {
	var vals []int
	for k := range s.words {
		vals = append(vals, int(k))
	}
	return vals
}

// String returns the set as a string in the form "{1 2 3}"
func (s MapSet) String() string {
	return commonString(&s)
}

// An IntSet is a set of small non-negative integers.
// Its zero value represents the empty set.
type IntSet struct {
	words []uint
}

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	word, bit := x/64, x%64
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set.
func (s *IntSet) Add(x int) {
	word, bit := x/64, x%64
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// GetsSet returns all elements in the set
func (s *IntSet) GetSet() []int {
	var vals []int
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<j) != 0 {
				vals = append(vals, 64*i+j)
			}
		}
	}
	return vals
}

// UnionWith sets s to the union of s and t.
func (s *IntSet) UnionWith(t Set) {
	for _, val := range t.GetSet() {
		s.Add(val)
	}
}

// String returns the set as a string in the form "{1 2 3}"
func (s IntSet) String() string {
	return commonString(&s)
}

func commonString(s Set) string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	vals := s.GetSet()
	if _, ok := s.(*MapSet); ok {
		sort.Ints(vals)
	}
	for _, val := range vals {
		if buf.Len() > len("{") {
			buf.WriteByte(' ')
		}
		fmt.Fprintf(&buf, "%d", val)
	}
	buf.WriteByte('}')
	return buf.String()
}

func main() {
	var x, y IntSet
	x.Add(1)
	x.Add(144)
	x.Add(9)
	fmt.Println(x.String())

	y.Add(9)
	y.Add(42)
	fmt.Println(y.String())

	x.UnionWith(&y)
	fmt.Println(x.String())

	fmt.Println(x.Has(9), x.Has(123))
}
