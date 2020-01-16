// Rewrite topoSort to use maps isntead of slices and elimitnate the inisial sort.
// Verify that the results though nondeterministic are valdi topological orderings

package main

import (
	"fmt"
)

var prereqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},
	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},
	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

func main() {
	for i, course := range topoSort(prereqs) {
		fmt.Printf("%d\t%s\n", i+1, course)
	}
}

func topoSort(m map[string][]string) []string {
	var order []string
	seen := make(map[string]bool)
	var visitAll func(string)
	visitAll = func(item string) {
		for _, v := range m[item] {
			if !seen[v] {
				seen[v] = true
				visitAll(v)
				order = append(order, v)
			}
		}
		if !seen[item] {
			order = append(order, item)
			seen[item] = true
		}
	}
	for k := range m {
		visitAll(k)
	}
	return order
}
