// Write a String method for the *tree type in gopl.io/ch4/treesort (4.4) that
// reveals the sequence of values in the tree

package main

import (
	"bytes"
	"fmt"
)

type tree struct {
	value       int
	left, right *tree
}

// Sort sorts values in place
func Sort(values []int) {
	var root *tree
	for _, v := range values {
		root = add(root, v)
	}
	appendValues(values[:0], root)
}

// appendValues appends the elements of t to values in order and returns the
// resulting slice.
func appendValues(values []int, t *tree) []int {
	if t != nil {
		values = appendValues(values, t.left)
		values = append(values, t.value)
		values = appendValues(values, t.right)
	}
	return values
}

func add(t *tree, value int) *tree {
	if t == nil {
		// Equivaluent to return &tree{value: value}
		t = new(tree)
		t.value = value
		return t
	}
	if value < t.value {
		t.left = add(t.left, value)
	} else {
		t.right = add(t.right, value)
	}
	return t
}

func stringSubtree(t *tree, buf *bytes.Buffer) {
	if t != nil {
		stringSubtree(t.left, buf)
		buf.WriteString(fmt.Sprintf("%d ", t.value))
		stringSubtree(t.right, buf)
	}
}

func (t *tree) String() string {
	var buf bytes.Buffer
	fmt.Fprint(&buf, "{")
	stringSubtree(t, &buf)
	fmt.Fprint(&buf, "}")
	return buf.String()
}

func main() {
	var t *tree
	for i := 0; i < 10; i++ {
		t = add(t, i)
	}
	fmt.Println(t)
}
