// Write a function to populate a mapping from element names -- p, div, span,
// and so on -- to the number of elements with that name in an HTML document
// tree.

package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"os"
	"sort"
)

func populateMapping(n *html.Node, m map[string]int) map[string]int {
	if n == nil {
		return m
	}
	if n.DataAtom.String() != "" {
		m[n.DataAtom.String()]++
	}
	m = populateMapping(n.FirstChild, m)
	m = populateMapping(n.NextSibling, m)
	return m
}

func main() {
	node, err := html.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	m := make(map[string]int)
	m = populateMapping(node, m)
	// Sort the tag names so that we iterate through the tags determiniscally every time (easier to test too)
	var tags []string
	for tag := range m {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	fmt.Printf("%-10v %v\n", "tag", "ocurrences")
	for _, tag := range tags {
		fmt.Printf("%-10v %v\n", tag, m[tag])
	}
}
