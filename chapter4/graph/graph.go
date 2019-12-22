package main

import (
	"fmt"
)

var graph = make(map[string]map[string]bool)

func addEdge(from, to string) {
	edges := graph[from]
	if edges == nil {
		edges = make(map[string]bool)
		graph[from] = edges
	}
	edges[to] = true
}

func hasEdge(from, to string) bool {
	return graph[from][to]
}

func main() {
	node0 := map[string]bool{"1": true}
	graph["0"] = node0
	addEdge("1", "0")
	fmt.Println(hasEdge("0", "1"))
	fmt.Println(hasEdge("1", "0"))
}
