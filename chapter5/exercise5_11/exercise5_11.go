// The instructor of the linear algebra course decids that calculus is now a prerequisite.
// Extend the toposort function to report cycles
package main

import (
	"fmt"
	"log"
	"reflect"
	"sort"
)

type Stack struct {
	stack []string
}

type Stacker interface {
	Pop() string
	Push(string) string
}

type Queue struct {
	queue []string
}

type Queuer interface {
	Enqueue(string)
	Dequeue() string
}

func (q *Queue) Enqueue(element ...string) {
	q.queue = append(q.queue, element...)
}

func (q *Queue) Dequeue() string {
	front := q.queue[0]
	q.queue = q.queue[1:]
	return front
}

func (q *Queue) isEmpty() bool {
	return len(q.queue) == 0
}

func (s *Stack) Pop() string {
	top := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return top
}

func (s *Stack) Push(element string) {
	s.stack = append(s.stack, element)
}

func (s *Stack) isEmpty() bool {
	return len(s.stack) == 0
}

func NewQueue() *Queue {
	return &Queue{}
}

func NewStack() Stack {
	return Stack{}
}

var prereqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},
	// "linear algebra": {"calculus"},
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
	var visitAll func(items []string)

	visitAll = func(items []string) {
		queue := NewQueue()
		for _, item := range items {
			alreadySeen := make(map[string]bool)
			queue.Enqueue(item)
			var v string
			for !queue.isEmpty() {
				// fmt.Printf("%+v\n", queue)
				v = queue.Dequeue()
				if alreadySeen[v] {
					log.Fatalf("There's a loop with these courses: %v\n", reflect.ValueOf(alreadySeen).MapKeys())
				}
				if !seen[v] {
					seen[v] = true
					alreadySeen[v] = true
					queue.Enqueue(m[v]...)
					fmt.Printf("after appending %+v\n", queue)

					if queue.isEmpty() {
						sort.Sort(sort.Reverse(sort.StringSlice(order)))

					}
					order = append(order, v)
				}
			}
		}
	}
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	visitAll(keys)
	return order
}
