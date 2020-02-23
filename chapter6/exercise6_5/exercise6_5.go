// The type of each word used by IntSet is uint64, but 64-bit airthmetic may be
// ineffictinetn on a 32-bit platform. Modify the program to use the uint type,
// which is the most efficent unsigned interger type for the platform. Instead
// of dividying by 64, define a constant holding the effective size of unint in
// bits, 32, or 64. You can use the perhaps too-clever expression 32 <<
// (^uint(0) >> 63) for this purpose.
package main

import (
	"bytes"
	"fmt"
	"math/bits"
)

// An IntSet is a set of small non-negative integers.
// Its zero value represents the empty set.
type IntSet struct {
	words []uint
	size  uint
}

func (s *IntSet) Elems() []uint {
	var buf []uint
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < bits.UintSize; j++ {
			if word&(1<<uint(j)) != 0 {
				buf = append(buf, uint(bits.UintSize*i+j))
			}
		}
	}
	return buf
}

func (s *IntSet) IntersectWith(t *IntSet) {
	var newSize uint
	for i, tword := range t.words {
		if i < len(s.words) {
			newSize += countBitsSet(s.words[i] & uint(tword))
			s.words[i] &= tword
		}
	}
	s.size = newSize
}

func (s *IntSet) DifferenceWith(t *IntSet) {
	var newSize uint
	for i, tword := range t.words {
		if i < len(s.words) {
			newSize += countBitsSet(s.words[i] ^ uint(tword))
			s.words[i] ^= tword
		} else {
			newSize += countBitsSet(tword)
			s.words = append(s.words, tword)
		}
	}
	s.size = newSize
}

func (s *IntSet) SymmetricDifference(t *IntSet) {
	var newSize uint
	for i, tword := range t.words {
		if i < len(s.words) {
			newSize += countBitsSet(s.words[i]) - countBitsSet(s.words[i]&uint(tword))
			s.words[i] = s.words[i] & ^(s.words[i] & tword)
		} else {
			newSize += countBitsSet(tword)
			s.words = append(s.words, tword)
		}
	}
	s.size = newSize
}

// AddAll adds a list of values to the set
func (s *IntSet) AddAll(values ...int) {
	for _, val := range values {
		s.Add(val)
	}
}

// Len returns the number of elements
func (s *IntSet) Len() uint {
	return s.size
}

// Remove remove x from the set
func (s *IntSet) Remove(x int) {
	word, bit := x/bits.UintSize, uint(x%bits.UintSize)
	if word < len(s.words) {
		s.words[word] &= ^(1 << bit)
	}
	s.size -= 1
}

// Clear remove all elements from the set
func (s *IntSet) Clear() {
	s.words = []uint{}
	s.size = 0
}

// Copy returns a copy of the set
func (s *IntSet) Copy() *IntSet {
	copyS := &IntSet{}
	for _, word := range s.words {
		copyS.words = append(copyS.words, word)
	}
	copyS.size = s.size
	return copyS
}

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	word, bit := x/bits.UintSize, uint(x%bits.UintSize)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set.
func (s *IntSet) Add(x int) {
	word, bit := x/bits.UintSize, uint(x%bits.UintSize)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
	s.size += 1
}

func countBitsSet(word uint) uint {
	var count uint
	var temp uint = word & (1<<bits.UintSize - 1)
	for i := 0; i < bits.UintSize && temp != 0; i++ {
		count += temp & 1
		temp = temp >> 1
	}
	return count
}

// UnionWith sets s to the union of s and t.
func (s *IntSet) UnionWith(t *IntSet) {
	var newSize uint
	for i, tword := range t.words {
		if i < len(s.words) {
			newSize += countBitsSet(s.words[i] | uint(tword))
			s.words[i] |= tword
		} else {
			newSize += countBitsSet(tword)
			s.words = append(s.words, tword)
		}
	}
	s.size = newSize
}

// String returns the set as a string in the form "{1 2 3}"
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < bits.UintSize; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", bits.UintSize*i+j)
			}
		}
	}
	buf.WriteString(fmt.Sprintf("}, size: %d", s.size))
	return buf.String()
}

func main() {
	var x, y IntSet
	x.Add(1)
	y.Add(144)
	x.Add(9)
	fmt.Println("x:", x.String())

	y.Add(9)
	y.Add(42)
	fmt.Println("y:", y.String())

	xUnionY := x.Copy()
	xUnionY.UnionWith(&y)
	fmt.Println("after x.Copy().UnionWith(y), x:", xUnionY.String())

	xDiffY := x.Copy()
	xDiffY.DifferenceWith(&y)
	fmt.Println("after x.Copy().DifferenceWith(y), x:", xDiffY.String())

	xIntersectY := x.Copy()
	xIntersectY.IntersectWith(&y)
	fmt.Println("after x.Copy().IntersectWith(y), x:", xIntersectY.String())

	xSymmDiffY := x.Copy()
	xSymmDiffY.SymmetricDifference(&y)
	fmt.Println("after x.Copy().SymmetricDifference(y), x:", xSymmDiffY.String())

	fmt.Println("Iterating over xUnionY using Elems with range")
	for i, val := range xUnionY.Elems() {
		fmt.Printf("i: %d, xUnionY[%d]: %v\n", i, i, val)
	}
}
